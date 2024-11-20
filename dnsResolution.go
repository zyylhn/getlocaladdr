package getlocaladdr

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
)

// DNSResolution 将给定的域名解析成ip(默认返回第一个)
func DNSResolution(domain string) (net.IP, error) {
	//先检查是否是ip地址
	ip := net.ParseIP(domain)
	if ip != nil {
		return ip, nil
	}
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}
	if len(ips) > 0 {
		return ips[0], nil
	} else {
		return nil, errors.New("the number of resolved ip addresses is 0")
	}
}

func DNSResolutionWithDNSServer(domain string, server []string, ctx context.Context) (net.IP, error) {
	var servers []*net.Resolver
	for _, s := range server {
		if net.ParseIP(s) == nil {
			continue
		}
		resolver := &net.Resolver{
			PreferGo: true, // 使用Go的DNS解析器
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				// 指定自定义的DNS服务器地址
				dialer := net.Dialer{}
				return dialer.DialContext(ctx, "udp", net.JoinHostPort(s, strconv.Itoa(53))) // 以8.8.8.8为例
			},
		}
		servers = append(servers, resolver)
	}
	for _, s := range servers {
		ips, err := s.LookupIP(ctx, "ip", domain)
		if err != nil {
			continue
		}
		if len(ips) > 0 {
			return ips[0], nil
		}
	}
	return nil, errors.New("the number of resolved ip addresses is 0")
}

// DNSResolutionFromUrl 从url中提取域名进行解析
func DNSResolutionFromUrl(url string) (net.IP, error) {
	domain, err := MatchingDomain(url)
	if err != nil {
		return nil, err
	}
	return DNSResolution(domain)
}

// MatchingDomain 从url中匹配域名或者ip地址
func MatchingDomain(url string) (string, error) {
	domainRegex := regexp.MustCompile(`^(?:https?://)?([^:/\s]+)`)
	domain := domainRegex.FindStringSubmatch(url)
	if len(domain) < 1 {
		return "", fmt.Errorf("description Failed to obtain the domain name or ip address:%v", url)
	}
	return domain[1], nil
}
