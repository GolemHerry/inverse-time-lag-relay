package main

import (
	"fmt"
	"inverse_time_lag_relay/curve"
	"inverse_time_lag_relay/relay"
	"log"
	"os"
	"os/signal"
)

func main() {
	fmt.Println(VERSION)
	fmt.Println("Inverse time lag relay developed by group 28")
	if err := ParseConfiguration("config.yaml"); err != nil {
		os.Exit(1)
	}
	curve.Init(config.Curve.Iop, config.Curve.C, config.Curve.K)
	//runtime.GOMAXPROCS(8)
	inverseTimeLagRelay := relay.NewRelay(config.Relay.SampleArgs)
	inverseTimeLagRelay.Run()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Printf("exited by interrupt \n")
}
