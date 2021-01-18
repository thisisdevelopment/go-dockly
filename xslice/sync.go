package xslice

import (
	"math/rand"
	"sync"
)

// SyncSlice for concurrent access
type SyncSlice struct {
	sync.RWMutex
	items []interface{}
}

// SyncSliceItem is an element in the slice
type SyncSliceItem struct {
	Idx int
	Val interface{}
}

// NewSyncSlice constructs a concurrent slice
func NewSyncSlice() *SyncSlice {
	return &SyncSlice{
		items: make([]interface{}, 0),
	}
}

// Append adds an element to the slice
func (s *SyncSlice) Append(item interface{}) {
	s.Lock()
	defer s.Unlock()
	s.items = append(s.items, item)
}

// Del removes the element at index from slice
func (s *SyncSlice) Del(i int) {
	s.Lock()
	defer s.Unlock()
	s.items = append(s.items[:i], s.items[i+1:]...)
}

// Pop removes the last element from slice and returns both
func (s *SyncSlice) Pop() <-chan SyncSliceItem {
	var ch = make(chan SyncSliceItem)
	go func() {
		s.Lock()
		defer s.Unlock()
		var last = len(s.items) - 1
		ch <- SyncSliceItem{last, s.items[last]}
		s.items = s.items[:last]
		close(ch)
	}()

	return ch
}

// Shift removes the first element from slice and returns both
func (s *SyncSlice) Shift() <-chan SyncSliceItem {
	var ch = make(chan SyncSliceItem)
	go func() {
		s.Lock()
		defer s.Unlock()
		ch <- SyncSliceItem{0, s.items[0]}
		s.items = s.items[1:]
		close(ch)
	}()

	return ch
}

// UnShift adds an element at first index to slice
func (s *SyncSlice) UnShift(v interface{}) {
	s.Lock()
	defer s.Unlock()
	var is []interface{} = make([]interface{}, 1)
	is[0] = v
	s.items = append(is, s.items...)
	return
}

// Contains returns true if value is contained in slice
func (s *SyncSlice) Contains(val interface{}) bool {
	s.Lock()
	defer s.Unlock()
	for _, x := range s.items {
		if x == val {
			return true
		}
	}

	return false
}

// Filter ing without allocating
func (s *SyncSlice) Filter(val interface{}) []interface{} {
	s.Lock()
	defer s.Unlock()
	b := s.items[:0]
	for _, x := range s.items {
		if val == x {
			b = append(b, x)
		}
	}

	return b
}

// Len returns the length of all elements in the slice
func (s *SyncSlice) Len() int {
	s.Lock()
	defer s.Unlock()
	return len(s.items)
}

// Reverse the order of a slice
func (s *SyncSlice) Reverse() {
	s.Lock()
	defer s.Unlock()
	for i := len(s.items)/2 - 1; i >= 0; i-- {
		var opp = len(s.items) - 1 - i
		s.items[i], s.items[opp] = s.items[opp], s.items[i]
	}
}

// Shuffle randomizes the order of a slice
func (s *SyncSlice) Shuffle() {
	s.Lock()
	defer s.Unlock()
	for i := len(s.items) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s.items[i], s.items[j] = s.items[j], s.items[i]
	}
}

// Get returns the element at index of slice or nil
func (s *SyncSlice) Get(i int) <-chan SyncSliceItem {
	var ch = make(chan SyncSliceItem)
	go func() {
		s.Lock()
		defer s.Unlock()
		ch <- SyncSliceItem{i, s.items[i]} // todo be careful with out of range errors - use Len
		close(ch)
	}()

	return ch
}

// LastElm returns the last element of the slice
func (s *SyncSlice) LastElm() <-chan SyncSliceItem {
	var ch = make(chan SyncSliceItem)
	go func() {
		s.Lock()
		defer s.Unlock()
		ch <- SyncSliceItem{len(s.items), s.items[len(s.items)-1]}
		close(ch)
	}()

	return ch
}

// Iter outputs each item over a channel to the caller
func (s *SyncSlice) Iter() <-chan SyncSliceItem {
	var ch = make(chan SyncSliceItem)

	go func() {
		s.Lock()
		defer s.Unlock()
		for i, val := range s.items {
			ch <- SyncSliceItem{i, val}
		}
		close(ch)
	}()

	return ch
}
