package getLocalAddr

import (
	"fmt"
	"testing"
	"time"
)

func TestGetFreePort(t *testing.T) {
	for i := 0; i < 1; i++ {
		go func(a int) {
			for {
				fmt.Print(a, "  ")
				//fmt.Println(getAvailablePortByRandom(nil))
				//_, err := getAvailablePortByRandom(nil)
				//if err != nil {
				//	fmt.Print(a, "  ")
				//	fmt.Println(err)
				//}
			}
		}(i)
	}
	time.Sleep(time.Second * 10)
}

func TestGetFreePortWithRange(t *testing.T) {
	fmt.Println(GetFreePortWithRange("127.0.0.3", 7001, 7003))
}
