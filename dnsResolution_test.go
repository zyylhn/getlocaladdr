package getlocaladdr

import (
	"context"
	"fmt"
	"testing"
)

func TestDNSResolution(t *testing.T) {
	fmt.Println(DNSResolution("www.baidu.com"))
	fmt.Println(DNSResolutionWithDNSServer("www.baidu.com", []string{"8.8.8.8"}, context.Background()))
}
