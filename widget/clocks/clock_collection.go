package clocks

import (
	"sort"
	"time"
)

const (
	SortAlphabetical  SortType = "alphabetical"
	SortChronological SortType = "chronological"
)

type SortType string

type ClockCollection struct {
	Clocks []Clock
	Sort   SortType
}

func (clocks *ClockCollection) Sorted() []Clock {
	if SortChronological == clocks.Sort {
		clocks.SortedChronologically()
	} else {
		clocks.SortedAlphabetically()
	}

	return clocks.Clocks
}

func (clocks *ClockCollection) SortedAlphabetically() {
	sort.Slice(clocks.Clocks, func(i, j int) bool {
		clock := clocks.Clocks[i]
		other := clocks.Clocks[j]

		return clock.Label < other.Label
	})
}

func (clocks *ClockCollection) SortedChronologically() {
	now := time.Now()
	sort.Slice(clocks.Clocks, func(i, j int) bool {
		clock := clocks.Clocks[i]
		other := clocks.Clocks[j]

		return clock.ToLocal(now).String() < other.ToLocal(now).String()
	})
}
