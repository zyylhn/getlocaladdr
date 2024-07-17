package main

import (
	"fmt"
	"github.com/zyylhn/getlocaladdr"
)

func main() {
	a := make(map[int]struct{})
	a[10000] = struct{}{}
	a[10001] = struct{}{}
	fmt.Println(getLocalAddr.GetFreePortMap(a))
}
