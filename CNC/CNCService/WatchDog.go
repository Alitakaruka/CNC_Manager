package CNCService

import (
	"context"
	"fmt"
	"time"
)

type WatchDog struct {
	ttl  int64
	stop context.CancelFunc
	// WG       sync.WaitGroup
	timer    *time.Timer
	isStoped chan struct{}
}

func NewWatchDog(Seconds int64, killFunc func()) *WatchDog {

	wd := &WatchDog{
		ttl:      Seconds,
		isStoped: make(chan struct{}),
	}

	ctx, fn := context.WithCancel(context.Background())
	wd.stop = fn

	wd.timer = time.NewTimer(time.Second * time.Duration(Seconds))
	go func() {
		for {
			select {
			case <-wd.timer.C:
				close(wd.isStoped)
				return
			case <-ctx.Done():
				fmt.Println("WD stop!")
				// wd.WG.Done()
				close(wd.isStoped)
				if !wd.timer.Stop() {
					<-wd.timer.C
					wd.timer.Stop()
				}
				return
			}
		}

	}()

	return wd
}

func (wd *WatchDog) Wait() chan struct{} {
	return wd.isStoped
}

func (wd *WatchDog) Alive() {
	// log.Println(string(debug.Stack()))
	// log.Println("alive")
	wd.timer.Reset(time.Second * time.Duration(wd.ttl))
}

func (wd *WatchDog) Close() {
	wd.Alive()
	// wd.WG.Done()
	if wd.stop != nil {
		wd.stop()
	}
}
