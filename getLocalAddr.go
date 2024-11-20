package getlocaladdr

import (
	"fmt"
	"net"
	"runtime"
	"sync"
)

var GetLocalIpFail = fmt.Errorf("failed get local ip address")

// GetLocalIPWithTargetAddr 根据目标地址获取本地连接地址
func GetLocalIPWithTargetAddr(ip string, port int) (string, error) {
	//todo 参考exp引擎即可
	//todo 完成中后将exp引擎等地方的函数替换下来
	panic("todo")
}

var DefaultGetLocalIP *GetLocalIPToConnectTarget

type GetLocalIPToConnectTarget struct {
	lock sync.Mutex
}

func (g *GetLocalIPToConnectTarget) GetLocalIPWithTargetIP(ip string) string {
	//g.lock.Lock()
	re := getLocalIPWithTargetIP(ip)
	//g.lock.Unlock()
	return re
}

func GetLocalIPWithTargetIP(ip string) string {
	return DefaultGetLocalIP.GetLocalIPWithTargetIP(ip)
}

// GetLocalIPWithTargetIP 根据目标ip获取本地连接ip
func getLocalIPWithTargetIP(ip string) string {
	//获取当前网卡数量，如果只有一个就直接返回其地址
	card, _ := getLocalIP4NetCard()
	if card != nil && len(card) == 1 {
		return card[0].IP.String()
	}
	//根据路由表获取
	var ipByRoute string
	switch runtime.GOOS {
	case "linux":
		ipByRoute, _ = GetLocalIPWithCmdOnLinux(ip)
	case "darwin":
		ipByRoute, _ = GetLocalIPWithCmdOnDarwin(ip)
	case "windows":
		ipByRoute, _ = GetLocalIPWithCmdOnWindows(ip)
	default:
		ipByRoute, _ = GetLocalIPWithCmd(ip)
	}
	if ipByRoute != "" && net.ParseIP(ipByRoute) != nil {
		return ipByRoute
	}
	//查看网卡包含关系
	var ipByContains string
	ipByContains, _ = NetCardContains(net.ParseIP(ip))
	return ipByContains
}

// 获取本地除localhost所有ipv4的网卡
func getLocalIP4NetCard() ([]*net.IPNet, error) {
	var re []*net.IPNet
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range addr {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				re = append(re, ipNet)
			}
		}
	}
	return re, nil
}

// GetInterfaceIP 根据网卡名字来获取其ip地址
func GetInterfaceIP(ifaceName string) (string, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return "", err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("IP address not found for interface")
}

// NetCardContains 网卡包含筛选，如果目标地址仅在一张网卡的范围内（应该不会出现多张，但不排除特殊情况,这里只取第一张），就采用该网卡地址
func NetCardContains(ip net.IP) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.Contains(ip) {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("interface not found")
}

func init() {
	DefaultGetLocalIP = &GetLocalIPToConnectTarget{lock: sync.Mutex{}}
}
