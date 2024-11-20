package getlocaladdr

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
)

// GetFreePortWithError 获取一个空闲的端口，并且会绕过传入的白名单端口，获取的时候可以保证是空闲的，使用时并不能保证是空闲的，所以在并发获取的时候要配合白名单使用
func GetFreePortWithError(whitelist []int) (int, error) {
	return Default.RandomGetFreePort(whitelist)
}

// GetFreePort 获取一个空闲端口，并且会绕过传入的白名单端口，获取的时候可以保证是空闲的，使用时并不能保证是空闲的，所以在并发获取的时候要配合白名单使用
func GetFreePort(whitelist []int) int {
	port, err := Default.RandomGetFreePort(whitelist)
	if err != nil {
		panic(err)
	}
	return port
}

func GetFreePortMap(whitelist map[int]struct{}) int {
	return Default.GetFreePortMap(whitelist)
}

type GetLocalFreePort struct {
	lock      *sync.Mutex
	basePort  int //获取端口的基准端口号
	usePort   int //分配出去的端口，会不断做累加
	maxLength int //单个对象的最大长度，也就是需要保证当前对象分配不重复的端口最多端口数，需要小于64000-basePort，当为负数的时候就可以随机获取不需要保证不重复
}

var Default GetLocalFreePort

// GetFreePortMap 使用累加的方式获取可用端口
func (g *GetLocalFreePort) GetFreePortMap(whitelist map[int]struct{}) int {
	var port int
	//判断可用端口是否被用完了并进行获取端口
G:
	port = g.getPort()
	if isPortWhitelistedMap(port, whitelist) {
		goto G
	}
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		goto G
	}
	_ = l.Close()
	return port
}

// 获取一个端口
func (g *GetLocalFreePort) getPort() int {
	var port int
	g.lock.Lock()
	if g.usePort >= g.basePort+g.maxLength {
		g.usePort = g.basePort
	}
	port = g.usePort
	g.usePort++
	g.lock.Unlock()
	return port
}

func (g *GetLocalFreePort) RandomGetFreePort(whitelist []int) (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port
	if isPortWhitelisted(port, whitelist) {
		return g.RandomGetFreePort(whitelist) // Try again if the port is whitelisted
	}

	return port, nil
}

func isPortWhitelisted(port int, whitelist []int) bool {
	for _, p := range whitelist {
		if port == p {
			return true
		}
	}
	return false
}

func isPortWhitelistedMap(port int, whitelist map[int]struct{}) bool {
	for p := range whitelist {
		if port == p {
			return true
		}
	}
	return false
}

//
//func isPortWhitelistedSyncMap(port int, whitelist *sync.Map) bool {
//	var flag bool
//	whitelist.Range(func(key, value any) bool {
//		if key.(int)==port{
//			flag=true
//			return false
//		}
//		return true
//	})
//	return flag
//}

//func getAvailablePortMap(whitelist map[int]struct{}) (int, error) {
//	l, err := net.Listen("tcp", ":0")
//	//todo 并发过高会出现listen tcp :0: socket: too many open files，改成使用随机数加验证的穿插方式
//	if err != nil {
//		return 0, err
//	}
//	_ = l.Close()
//
//	port := l.Addr().(*net.TCPAddr).Port
//	if isPortWhitelistedMap(port, whitelist) {
//		return getAvailablePortMap(whitelist) // Try again if the port is whitelisted
//	}
//
//	return port, nil
//}

//// 使用监听的方式寻找可用的端口号
//func getAvailablePortByRandom(whitelist map[int]struct{}) (int, error) {
//	random := randomPort(10000, 65535)
//	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", random))
//	if err != nil {
//		return 0, err
//	}
//	_ = l.Close()
//	if isPortWhitelistedMap(random, whitelist) {
//		return getAvailablePortByRandom(whitelist)
//	}
//	return random, nil
//}

// GetFreePortWithRange 在指定范围获取一个指定地址可用端口号
func GetFreePortWithRange(address string, startPort, endPort int) int {
	if !IsLocalIP(address) {
		address = "0.0.0.0"
	}
	var numThreads = 10
	randPort := randomPort(startPort, endPort) // 在范围内随机一个端口号
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	//fmt.Println("随机端口：", randPort)
	freePort := make(chan int, 1)
	defer func() {
		close(freePort)
	}()
	// 确定遍历方向和步长
	direction := 1
	if randPort > (startPort+endPort)/2 {
		direction = -1
	}

	ports := make(chan int, numThreads)
	var wg sync.WaitGroup

	// 启动并发线程
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range ports {
				select {
				case <-ctx.Done():
					return
				default:
				}
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
				if err != nil {
					if strings.Contains(err.Error(), "refused") {
						select {
						case freePort <- port:
						default:
						}
						cancel()
						return
					}
					continue
				}
				_ = conn.Close()
			}
		}()
	}

	// 向通道中发送待检查的端口号
	go func() {
		defer close(ports)
		port := randPort
		var flag bool
		//fmt.Println("步长", direction)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if port < startPort || port > endPort {
				if flag == false {
					//fmt.Println("步长转换")
					direction = direction - direction*2
					port = randPort + direction
					flag = true
					continue
				} else {
					//fmt.Println("步长已转换")
					break
				}
			}
			//fmt.Println("任务端口", port)
			ports <- port
			port += direction
		}

	}()

	wg.Wait() // 等待所有线程完成
	//re := <-freePort
	select {
	case re := <-freePort:
		return re
	default:
		return 0
	}
}

func randomPort(startPort, endPort int) int {
	return startPort + rand.Intn(endPort-startPort+1)
}

func init() {
	Default = GetLocalFreePort{
		lock:      &sync.Mutex{},
		basePort:  10000,
		usePort:   10001,
		maxLength: 50000,
	}
}
