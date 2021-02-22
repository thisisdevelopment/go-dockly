package xslice

import (
	"math/rand"
)

// Reference https://github.com/golang/go/wiki/SliceTricks

// LastElm returns the last element of the slice
func LastElm(s []string) (string, bool) {
	if len(s) == 0 {
		return "", false
	}
	return s[len(s)-1], true
}

// Uniq deduplicates the slice
func Uniq(s []string) []string {
	var m = make(map[string]bool)

	t := s[:0]
	for _, v := range s {
		if _, exists := m[v]; !exists {
			t = append(t, v)
			m[v] = true
		}
	}

	return t
}

// Cut removes the elements between i and j from the slice
func Cut(s []string, i, j int) ([]string, bool) {
	if i < 0 || i > len(s) || j < 0 || j > len(s) || j < i {
		return s, false
	}

	copy(s[i:], s[j:])
	return s[:len(s)-(j-i)], true
}

// Strip removes all occurrences of val from slice
func Strip(s []string, val string) []string {
	res := s[:0]
	for _, x := range s {
		if x != val {
			res = append(res, x)
		}
	}

	return res
}

// Del removes the element at index from slice
func Del(s []string, i int) ([]string, bool) {
	if i < 0 || i >= len(s) {
		return s, false
	}

	s[i] = s[len(s)-1]
	return s[:len(s)-1], true
}

// Pop removes the last element from slice and returns both
func Pop(s []string) (string, []string, bool) {
	last, ok := LastElm(s)
	if !ok {
		return last, s, ok
	}

	return last, s[:len(s)-1], ok
}

// Shift removes the first element from slice and returns both
func Shift(s []string) (string, []string, bool) {
	if len(s) == 0 {
		return "", s, false
	}

	return s[0], s[1:], true
}

// UnShift adds an element at first index to slice
func UnShift(s []string, v string) []string {
	return append([]string{v}, s...)
}

// Filter ing without allocating
func Filter(s []string, val string) []string {
	b := s[:0]
	for _, x := range s {
		if val == x {
			b = append(b, x)
		}
	}

	return b
}

// Contains returns true if value is contained in slice
func Contains(s []string, val string) bool {
	for _, x := range s {
		if x == val {
			return true
		}
	}

	return false
}

// Reverse the order of a slice
func Reverse(s []string) []string {
	var opp int
	for i := len(s)/2 - 1; i >= 0; i-- {
		opp = len(s) - 1 - i
		s[i], s[opp] = s[opp], s[i]
	}

	return s
}

// Shuffle randomizes the order of a slice
func Shuffle(s []string) []string {
	for i := len(s) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}

	return s
}

// ContainsAny returns any matches between the src and tgt string slice
func ContainsAny(src, tgt []string) (hits []string, ok bool) {
	var m = make(map[string]bool)
	for _, v := range src {
		m[v] = true
	}
	for _, v := range tgt {
		if m[v] {
			hits = append(hits, v)
		}
	}

	return hits, len(hits) > 0
}

// Insert s one string slice in the other at the given index
func Insert(s, ins []string, i int) []string {
	if i > len(s) {
		i = len(s)
	}
	var start = make([]string, i)
	copy(start, s[:i])
	start = append(start, ins...)
	return append(start, s[i:]...)
}

// AppendNotEmpty appends string slice items which are not empty
func AppendNotEmpty(s []string, appends []string) []string {
	for _, v := range appends {
		if v != "" {
			s = append(s, v)
		}
	}
	return s
}
