package relay

import (
	"context"
	"fmt"
	"inverse_time_lag_relay/ADC"
	"inverse_time_lag_relay/curve"
	"inverse_time_lag_relay/fft"
	"log"
	"math"
	"sync"
	"time"
)

const (
	relayClosed = iota
	relayOpened
	relayFaulted
)

type Relay struct {
	lock           sync.RWMutex
	relayState     int
	tickChan       <-chan time.Time
	SampleArgs     []float64 `config:"sampleargs"`
	lastTimeToJump float64
}

func NewRelay(sample []float64) *Relay {
	return &Relay{
		lock:           sync.RWMutex{},
		relayState:     relayClosed,
		tickChan:       time.NewTicker(time.Millisecond * 20).C,
		SampleArgs:     sample,
		lastTimeToJump: math.MaxFloat64,
	}
}

func (r *Relay) detect() float64 {
	res := fft.Calculate(ADC.GenerateSample(r.SampleArgs))
	fmt.Printf("%#v\n", res)
	return curve.StdCurve.GetTime(res.Amplitude())
}

func (r *Relay) action(ctx context.Context, timeToJump float64) {
	time.Sleep(time.Millisecond * time.Duration(timeToJump))
	select {
	case <-ctx.Done():
		log.Println("attempt to jump but recovered")
		return
	default:
		//log.Println("relayState turns to opened")
		r.setRelayState(relayOpened)
	}
}

func (r *Relay) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Recovered in %v\n", r)
		}
	}()
	var lastCancel context.CancelFunc
	for {
		select {
		case <-r.tickChan:
			timeToJump := r.detect()
			ctx, cancel := context.WithCancel(context.Background())
			if timeToJump < 0 || r.lastTimeToJump < timeToJump {
				r.setLastTimeToJump(timeToJump)
				if lastCancel != nil {
					lastCancel()
				}
				return
			} else {
				r.setLastTimeToJump(timeToJump)
				if lastCancel != nil {
					lastCancel()
				}
				go r.action(ctx, timeToJump)
			}
			lastCancel = cancel
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func (r *Relay) setRelayState(state int) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.relayState = state
}

func (r *Relay) GetRelayState() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.relayState
}

func (r *Relay) setLastTimeToJump(timeToJump float64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.lastTimeToJump = timeToJump
}

func (r *Relay) getLastTimeToJump() float64 {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.lastTimeToJump
}
