package wtf

import (
	"context"
	"time"
)

type Refresher interface {
	Refresh()
	RefreshInterval() int
}

type Refreshable interface {
	Refresher
	Enabler
}

type RefreshScheduler struct {
	items  []Refreshable
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRefresher(items []Refreshable) *RefreshScheduler {
	return &RefreshScheduler{
		items:  items,
		cancel: func() {},
	}
}

func (s *RefreshScheduler) Refresh() {
	for _, item := range s.items {
		item.Refresh()
	}
}

func (s *RefreshScheduler) ScheduleAutoRefresh() {
	s.cancel() // cancel previous schedulers
	s.ctx, s.cancel = context.WithCancel(context.Background())

	for _, item := range s.items {
		go schedule(item, s.ctx)
	}
}

func schedule(widget Refreshable, ctx context.Context) {
	// Kick off the first refresh and then leave the rest to the timer
	widget.Refresh()

	interval := time.Duration(widget.RefreshInterval()) * time.Second

	if interval <= 0 {
		return
	}

	tick := time.NewTicker(interval)

	for {
		select {
		case <-tick.C:
			if widget.Enabled() {
				widget.Refresh()
			} else {
				tick.Stop()
				return
			}
		case <-ctx.Done():
			tick.Stop()
			return
		}
	}
}
