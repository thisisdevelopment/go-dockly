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
func NewSyncSlice(initial []interface{}) *SyncSlice {

	var s = &SyncSlice{
		items: make([]interface{}, 0),
	}

	for _, v := range initial {
		s.items = append(s.items, v)
	}

	return s
}

// LastElm returns the last element of the slice
func (s *SyncSlice) LastElm() <-chan *SyncSliceItem {
	var ch = make(chan *SyncSliceItem)
	go func() {
		s.Lock()
		defer s.Unlock()
		if s.Len() > 0 {
			ch <- &SyncSliceItem{len(s.items), s.items[len(s.items)-1]}
		}
		close(ch)
	}()

	return ch
}

// Uniq deduplicates the slice
func (s *SyncSlice) Uniq() {
	s.Lock()
	defer s.Unlock()

	var m = make(map[interface{}]bool)
	for _, v := range s.items {
		m[v] = true
	}

	s.items = s.items[:0]
	for k, _ := range m {
		s.items = append(s.items, k)
	}
}

// Cut removes the elements between i and j from the slice
func (s *SyncSlice) Cut(i, j int) bool {
	s.Lock()
	defer s.Unlock()

	if i < 0 || i > len(s.items) || j < 0 || j > len(s.items) || j < i {
		return false
	}

	copy(s.items[i:], s.items[j:])
	s.items = s.items[:len(s.items)-(j-i)]
	return true
}

// Strip removes all occurrences of val from slice
func (s *SyncSlice) Strip(val interface{}) {
	s.Lock()
	defer s.Unlock()

	res := s.items[:0]
	for _, x := range s.items {
		if x != val {
			res = append(res, x)
		}
	}

	s.items = res
}

// Del removes the element at index from slice
func (s *SyncSlice) Del(i int) bool {
	s.Lock()
	defer s.Unlock()

	if i < 0 || i >= len(s.items) {
		return false
	}

	s.items[i] = s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return true
}

// Pop removes the last element from slice and returns both
func (s *SyncSlice) Pop() <-chan *SyncSliceItem {
	var ch = make(chan *SyncSliceItem)
	go func() {
		s.Lock()
		defer s.Unlock()
		if len(s.items) > 0 {
			var last = len(s.items) - 1
			ch <- &SyncSliceItem{last, s.items[last]}
			s.items = s.items[:last]
		}
		close(ch)
	}()

	return ch
}

// Shift removes the first element from slice and returns it
func (s *SyncSlice) Shift() <-chan *SyncSliceItem {
	var ch = make(chan *SyncSliceItem)
	go func() {
		s.Lock()
		defer s.Unlock()
		if len(s.items) > 0 {
			ch <- &SyncSliceItem{0, s.items[0]}
			s.items = s.items[1:]
		}
		close(ch)
	}()

	return ch
}

// UnShift adds an element at first index to slice
func (s *SyncSlice) UnShift(v interface{}) {
	s.Lock()
	defer s.Unlock()

	s.items = append([]interface{}{v}, s.items...)
	return
}

// Filter ing without allocating
func (s *SyncSlice) Filter(val interface{}) {
	s.Lock()
	defer s.Unlock()

	b := s.items[:0]
	for _, x := range s.items {
		if val == x {
			b = append(b, x)
		}
	}

	s.items = b
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

// Reverse the order of a slice
func (s *SyncSlice) Reverse() {
	s.Lock()
	defer s.Unlock()

	var opp int
	for i := len(s.items)/2 - 1; i >= 0; i-- {
		opp = len(s.items) - 1 - i
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

// ContainsAny returns any matches from tgt slice
func (s *SyncSlice) ContainsAny(tgt []interface{}) (hits []interface{}, ok bool) {
	s.Lock()
	defer s.Unlock()

	var m = make(map[interface{}]bool)
	for _, v := range s.items {
		m[v] = true
	}
	for _, v := range tgt {
		if ok := m[v]; ok {
			hits = append(hits, v)
		}
	}

	return hits, len(hits) > 0
}

// Insert s one string slice in the other at the given index
func (s *SyncSlice) Insert(ins []interface{}, i int) {
	s.Lock()
	defer s.Unlock()

	if i > len(s.items) {
		i = len(s.items)
	}
	var start = make([]interface{}, i)
	copy(start, s.items[:i])
	start = append(start, ins...)
	s.items = append(start, s.items[i:]...)
}

// Append elements to the slice
func (s *SyncSlice) Append(items ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		s.items = append(s.items, v)
	}
}

// Clear s all elements from the slice
func (s *SyncSlice) Clear() {
	s.Lock()
	defer s.Unlock()
	s.items = s.items[:0]
}

// Len returns the length of all elements in the slice
func (s *SyncSlice) Len() int {
	s.Lock()
	defer s.Unlock()
	return len(s.items)
}

// Get returns the element at index of slice or nil
func (s *SyncSlice) Get(i int) <-chan *SyncSliceItem {
	var ch = make(chan *SyncSliceItem)
	go func() {
		s.Lock()
		defer s.Unlock()
		if i < 0 || i >= len(s.items) {
			ch = nil
		} else {
			ch <- &SyncSliceItem{i, s.items[i]}
		}
		close(ch)
	}()

	return ch
}

// Iter outputs each item over a channel to the caller
func (s *SyncSlice) Iter() <-chan *SyncSliceItem {
	var ch = make(chan *SyncSliceItem)
	go func() {
		s.Lock()
		defer s.Unlock()
		for i, val := range s.items {
			ch <- &SyncSliceItem{i, val}
		}
		close(ch)
	}()

	return ch
}

// AppendNotNil appends string slice items which are not empty
func (s *SyncSlice) AppendNotNil(appends []interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, val := range appends {
		if val != nil {
			s.items = append(s.items, val)
		}
	}
}
