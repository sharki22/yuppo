package debouncer

import (
	"sync"
	"time"
)

type Debouncer struct {
	mu       sync.Mutex
	timer    *time.Timer
	duration time.Duration
}

func New(d time.Duration) *Debouncer {
	return &Debouncer{duration: d}
}

func (d *Debouncer) Trigger(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		if !d.timer.Stop() {
			select {
			case <-d.timer.C:
			default:
			}
		}
		d.timer.Reset(d.duration)
		return
	}

	d.timer = time.AfterFunc(d.duration, func() {
		d.mu.Lock()
		d.timer = nil
		d.mu.Unlock()
		fn()
	})
}

func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
