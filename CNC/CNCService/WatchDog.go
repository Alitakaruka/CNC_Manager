package CNCService

import (
	"context"
	"sync"
	"time"
)

type WatchDog struct {
	ttl      int64
	stop     context.CancelFunc
	WG       sync.WaitGroup
	timer    *time.Timer
	isStoped chan struct{}
}

func NewWatchDog(Seconds int64, killFunc func()) *WatchDog {

	wd := &WatchDog{
		ttl: Seconds,
	}

	ctx, fn := context.WithCancel(context.Background())

	wd.WG = sync.WaitGroup{}
	wd.WG.Add(1)
	wd.stop = fn

	wd.timer = time.NewTimer(time.Second * time.Duration(Seconds))
	go func() {
		for {
			select {
			case <-wd.timer.C:
				if killFunc != nil {
					wd.WG.Done()
					killFunc()
					return
				} else {
					wd.WG.Done()
					return
				}
			case <-ctx.Done():
				wd.WG.Done()
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
	wd.WG.Wait()
	return wd.isStoped
}

func (wd *WatchDog) Alive() {
	wd.timer.Reset(time.Second * time.Duration(wd.ttl))
}

func (wd *WatchDog) Close() {
	wd.Alive()
	wd.stop()
}
