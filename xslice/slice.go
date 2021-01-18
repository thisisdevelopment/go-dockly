package xslice

import "math/rand"

// Reference https://github.com/golang/go/wiki/SliceTricks

// LastElm returns the last element of the slice
func LastElm(s []string) string {
	return s[len(s)-1]
}

// Uniq deduplicates the slice
func Uniq(s []string) []string {
	var m = make(map[string]struct{})
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
		}
	}

	var uniq = make([]string, len(m))
	var i = 0
	for v := range m {
		uniq[i] = v
		i++
	}

	return uniq
}

// Cut removes the elements between i and j from the slice
func Cut(a []string, i, j int) []string {
	copy(a[i:], a[j:])
	for k, n := len(a)-j+i, len(a); k < n; k++ {
		a[k] = ""
	}
	return a[:len(a)-j+i]
}

// Strip removes all occurrences of val from slice
func Strip(s []string, val string) []string {
	for i, x := range s {
		if x == val {
			s = Del(s, i)
		}
	}

	return s
}

// Del removes the element at index from slice
func Del(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
}

// Pop removes the last element from slice and returns both
func Pop(s []string) (string, []string) {
	return LastElm(s), s[:len(s)-1]
}

// Shift removes the first element from slice and returns both
func Shift(s []string) (string, []string) {
	return s[0], s[1:]
}

// UnShift adds an element at first index to slice
func UnShift(s []string, v string) []string {
	return append([]string{v}, s...)
}

// Filter ing without allocating
func Filter(a []string, val string) []string {
	b := a[:0]
	for _, x := range a {
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
	for i := len(s)/2 - 1; i >= 0; i-- {
		var opp = len(s) - 1 - i
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
