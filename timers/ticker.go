package timers

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Ticker struct {
	ticker *time.Ticker
	period int
}

//
// specify period int in seconds amount you need
//
func NewTicker(ctx context.Context, period int) (*Ticker, error) {
	if period == 0 {
		return nil, fmt.Errorf("period can't be 0")
	}

	t := new(Ticker)

	interval := time.Duration(period) * time.Second
	//minutes := period / 60
	//hours := minutes / 60
	//seconds := float32(period-minutes*60-hours*60*60) / 100
	//log.Printf("Starting ticker. Interval: %d hours %.2f minutes", hours, float32(minutes)+seconds)

	t.ticker = time.NewTicker(time.Second * interval)

	return t, nil
}

func (t *Ticker) Run(ctx context.Context, callbackFn func()) {
	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				t.Stop(ctx)
				break loop
			case tick := <-t.ticker.C:
				minutes := t.period / 60
				hours := minutes / 60
				seconds := float32(t.period-minutes*60-hours*60*60) / 100
				log.Printf("Ticker tick! Time is: `%s`. Next tick after %d hours %.2f minutes, on %s", tick.Format(time.RFC1123), hours, float32(minutes)+seconds, tick.Add(time.Duration(t.period)*time.Second))
				if callbackFn != nil {
					callbackFn()
				}
			}
		}
	}()
}

func (t *Ticker) Stop(ctx context.Context) {
	t.ticker.Stop()
}
