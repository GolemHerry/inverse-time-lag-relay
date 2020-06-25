package main

import (
	"fmt"
	"github.com/nathan-osman/go-rpigpio"
	"time"
)

func main() {
	pin, _ := rpi.OpenPin(2, rpi.IN)
	var lastVal, val rpi.Value
	var err error
	for {
		if val, err = pin.Read(); val != lastVal && err == nil {
			fmt.Println(val)
		}
		lastVal = val
		time.Sleep(time.Millisecond * 20)
	}
}
