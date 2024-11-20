package getlocaladdr

import (
	"os/exec"
	"strings"
)

func GetLocalIPWithCmd(ip string) (string, error) {
	return "", GetLocalIpFail
}

func GetLocalIPWithCmdOnWindows(ip string) (string, error) {
	//todo 比较麻烦，可能需要自行从路由表中提取
	return "", GetLocalIpFail
}

func GetLocalIPWithCmdOnDarwin(ip string) (string, error) {
	cmd := exec.Command("route", "get", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	var netCard string
	lines := strings.Split(string(output), "\n")
	//获取到网卡地址
	for _, line := range lines {
		if strings.Contains(line, "interface:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				netCard = fields[1]
				break
			}
		}
	}
	//根据网卡名称获取网卡地址
	return GetInterfaceIP(netCard)
}

// GetLocalIPWithCmdOnLinux 给定要路由过去的ip，查询使用哪个地址连接
func GetLocalIPWithCmdOnLinux(ip string) (string, error) {
	cmd := exec.Command("ip", "route", "get", ip)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	fields := strings.Fields(string(output))
	for i := 0; i < len(fields)-1; i++ {
		if fields[i] == "src" {
			return fields[i+1], nil
		}
	}
	return "", GetLocalIpFail
}

// GetDefaultRouteIp 获取默认路由的网卡地址
func GetDefaultRouteIp() (string, error) {
	return "", GetLocalIpFail
}

// GetDefaultRouteIpOnWindows 获取默认路由的网卡地址
func GetDefaultRouteIpOnWindows() (string, error) {
	return "", GetLocalIpFail
}

// GetDefaultRouteIpOnLinux 获取默认路由的网卡地址
func GetDefaultRouteIpOnLinux() (string, error) {
	return "", GetLocalIpFail
}

// GetDefaultRouteIpOnDarwin 获取默认路由的网卡地址
func GetDefaultRouteIpOnDarwin() (string, error) {
	return "", GetLocalIpFail
}
