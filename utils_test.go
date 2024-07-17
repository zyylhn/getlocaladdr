package getLocalAddr

import (
	"fmt"
	"testing"
)

func TestIsLocalIP(t *testing.T) {
	fmt.Println(IsLocalIP("127.0.0.2"))
}
