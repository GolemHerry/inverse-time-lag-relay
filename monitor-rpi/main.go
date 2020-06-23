package main

import (
	"fmt"
	"github.com/nathan-osman/go-rpigpio"
	"time"
)

func main() {
	pin, _ := rpi.OpenPin(1, rpi.IN)
	for {
		if val, _ := pin.Read(); val == rpi.HIGH {
			fmt.Println(val)
		}
		time.Sleep(time.Millisecond * 20)
	}
}
