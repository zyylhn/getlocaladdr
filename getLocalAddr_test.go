package getlocaladdr

import (
	"fmt"
	"net"
	"testing"
)

func TestNetCard(t *testing.T) {
	fmt.Println(NetCardContains(net.ParseIP("172.16.95.24")))
}
