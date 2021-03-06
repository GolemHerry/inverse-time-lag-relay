package relay

import (
	"context"
	"github.com/nathan-osman/go-rpigpio"
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

const PHASES = 3

var pin, _ = rpi.OpenPin(2, rpi.OUT)

type Relay struct {
	lock           sync.RWMutex
	relayState     int
	ticker         *time.Ticker
	SampleArgs     SampleArgs `config:"sampleargs"`
	lastTimeToJump [3]float64
}

type SampleArgs struct {
	Current []Scalar `config:"current"`
	Voltage []Scalar `config:"voltage"`
}

type Scalar struct {
	Name   string    `config:"name"`
	Values []float64 `config:"values"`
}

func NewRelay(sample SampleArgs) *Relay {
	return &Relay{
		lock:           sync.RWMutex{},
		relayState:     relayClosed,
		ticker:         time.NewTicker(time.Millisecond * 20),
		SampleArgs:     sample,
		lastTimeToJump: [3]float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64},
	}
}

func (r *Relay) probe(i int) float64 {
	res := fft.Calculate(ADC.GenerateSample(r.SampleArgs.Current[i].Values))
	return curve.StdCurve.GetTime(res.Amplitude())
}

func (r *Relay) action(ctx context.Context, timeToJump float64, name string) {
	ticker := time.Tick(time.Millisecond * time.Duration(timeToJump))
	now := time.Now()
LOOP:
	for {
		select {
		case <-ticker:
			break LOOP
		case <-ctx.Done():
			log.Printf("the %s attempt to jump but recovered", name)
			return
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
	if r.GetRelayState() == relayOpened {
		return
	}
	log.Printf("A  failure on  %s happend at %v , relayState turns to opened", name, now)
	r.setRelayState(relayOpened)
}

func (r *Relay) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Recovered in %v\n", r)
		}
	}()
	go func() {
		for {
			if r.GetRelayState() == relayOpened {
				if err := pin.Write(rpi.LOW); err != nil {
					log.Println(err)
					return
				}
			} else if r.GetRelayState() == relayClosed {
				if err := pin.Write(rpi.HIGH); err != nil {
					log.Println(err)
					return
				}
			}
			time.Sleep(time.Millisecond * 20)
		}
	}()
	var lastCancel [3]context.CancelFunc
	for {
		if r.GetRelayState() == relayOpened {
			log.Println("relay opened")
			return
		}
		select {
		case <-r.ticker.C:
			for i := 0; i < PHASES; i++ {
				timeToJump := r.probe(i)
				ctx, cancel := context.WithCancel(context.Background())
				if timeToJump < 0 || r.getLastTimeToJump(i) < timeToJump {
					r.setLastTimeToJump(i, timeToJump)
					if lastCancel[i] != nil {
						lastCancel[i]()
					}
				} else if r.getLastTimeToJump(i) > timeToJump {
					r.setLastTimeToJump(i, timeToJump)
					if lastCancel[i] != nil {
						lastCancel[i]()
					}
					go r.action(ctx, timeToJump, r.SampleArgs.Current[i].Name)
				} else {
					time.Sleep(time.Millisecond)
					continue
				}
				lastCancel[i] = cancel
			}
		default:
			time.Sleep(time.Millisecond * 20)
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

func (r *Relay) setLastTimeToJump(i int, timeToJump float64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.lastTimeToJump[i] = timeToJump
}

func (r *Relay) getLastTimeToJump(i int) float64 {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.lastTimeToJump[i]
}
