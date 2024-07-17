package getLocalAddr

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

func ExtractPortFromURL(u string) (int, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return 0, err
	}

	port := parsedURL.Port()
	if port == "" {
		// 如果 URL 中未指定端口号，则根据协议默认使用对应的端口号
		switch parsedURL.Scheme {
		case "http":
			port = "80"
		case "https":
			port = "443"
		default:
			return 0, fmt.Errorf("unknown protocol: %s", parsedURL.Scheme)
		}
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}

	return portNumber, nil
}

// IsLocalIP 判断给定地址是否是当前主机的ip地址
func IsLocalIP(ip string) bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			if v.IP.Equal(net.ParseIP(ip)) {
				return true
			}
		case *net.IPAddr:
			if v.IP.Equal(net.ParseIP(ip)) {
				return true
			}
		}
	}

	return false
}
