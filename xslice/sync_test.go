package xslice_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thisisdevelopment/go-dockly/xslice"
)

func TestSyncUniq(t *testing.T) {

	data := []struct{ in, out []interface{} }{
		{[]interface{}{}, []interface{}{}},
		{[]interface{}{"", "", ""}, []interface{}{""}},
		{[]interface{}{"a", "a"}, []interface{}{"a"}},
		{[]interface{}{"a", "b", "a"}, []interface{}{"a", "b"}},
		{[]interface{}{"a", "b", "a", "b"}, []interface{}{"a", "b"}},
		{[]interface{}{"a", "b", "b", "a", "b"}, []interface{}{"a", "b"}},
		{[]interface{}{"a", "a", "b", "b", "a", "b"}, []interface{}{"a", "b"}},
		{[]interface{}{"a", "b", "c", "a", "b", "c"}, []interface{}{"a", "b", "c"}},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		s.Uniq()
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestSyncCut(t *testing.T) {

	data := []struct {
		in, out    []interface{}
		start, end int
	}{
		{[]interface{}{}, []interface{}{}, 0, 0},
		{[]interface{}{"", "", ""}, []interface{}{""}, 1, 3},
		{[]interface{}{"a", "a"}, []interface{}{"a"}, 1, 2},
		{[]interface{}{"a", "b", "a"}, []interface{}{"a", "b"}, 2, 3},
		{[]interface{}{"a", "b", "a", "b"}, []interface{}{"a", "b"}, 2, 4},
		{[]interface{}{"a", "a", "b", "b", "a", "b"}, []interface{}{"a", "b"}, 0, 4},
		{[]interface{}{"a", "b", "c", "a", "b", "c"}, []interface{}{"a", "b", "c"}, 3, 6},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		go func() {
			s.Cut(exp.start, exp.end)
			for item := range s.Iter() {
				if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
					t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
				}
			}
		}()
		s.Clear()
	}
}

func TestSyncStrip(t *testing.T) {

	data := []struct {
		in, out []interface{}
		val     interface{}
	}{
		{[]interface{}{}, []interface{}{}, "blah"},
		{[]interface{}{"", "", ""}, []interface{}{}, ""},
		{[]interface{}{"a", "a"}, []interface{}{"a", "a"}, "b"},
		{[]interface{}{"a", "b", "a"}, []interface{}{"b"}, "a"},
		{[]interface{}{"c", "c", "c"}, []interface{}{}, "c"},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		s.Strip(exp.val)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestSyncDel(t *testing.T) {

	data := []struct {
		in, out []interface{}
		index   int
	}{
		{[]interface{}{}, []interface{}{}, 22},
		{[]interface{}{"", "", ""}, []interface{}{"", ""}, 1},
		{[]interface{}{"a", "a"}, []interface{}{"a"}, 0},
		{[]interface{}{"a", "b", "a"}, []interface{}{"a", "a"}, 1},
		{[]interface{}{"a", "b", "c"}, []interface{}{"a", "b"}, 2},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		s.Del(exp.index)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestSyncPop(t *testing.T) {

	data := []struct {
		in, out []interface{}
		pop     interface{}
	}{
		{[]interface{}{}, []interface{}{}, nil},
		{[]interface{}{"a"}, []interface{}{}, "a"},
		{[]interface{}{"a", "b"}, []interface{}{"a"}, "b"},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		pop := <-s.Pop()
		if pop != nil && pop.Val != exp.pop {
			t.Fatalf("%q didn't match %q\n", pop.Val, exp.pop)
		}
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestSyncShift(t *testing.T) {

	data := []struct {
		in, out []interface{}
		shift   interface{}
	}{
		{[]interface{}{}, []interface{}{}, nil},
		{[]interface{}{"a"}, []interface{}{}, "a"},
		{[]interface{}{"a", "b"}, []interface{}{"b"}, "a"},
		{[]interface{}{"a", "b", "c"}, []interface{}{"b", "c"}, "a"},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		shift := <-s.Shift()
		if shift != nil && shift.Val != exp.shift {
			t.Fatalf("%q didn't match %q\n", shift.Val, exp.shift)
		}
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestSyncUnShift(t *testing.T) {

	data := []struct {
		in, out []interface{}
		unshift interface{}
	}{
		{[]interface{}{}, []interface{}{"a"}, "a"},
		{[]interface{}{"b"}, []interface{}{"a", "b"}, "a"},
		{[]interface{}{"b", "c"}, []interface{}{"a", "b", "c"}, "a"},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		s.UnShift(exp.unshift)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestSyncFilter(t *testing.T) {

	data := []struct {
		in, out []interface{}
		filter  interface{}
	}{
		{[]interface{}{}, []interface{}{}, "a"},
		{[]interface{}{"c"}, []interface{}{}, "b"},
		{[]interface{}{"c"}, []interface{}{"c"}, "c"},
		{[]interface{}{"a", "b", "c"}, []interface{}{"b"}, "b"},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		s.Filter(exp.filter)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestSyncContains(t *testing.T) {
	var ok bool

	s := xslice.NewSyncSlice([]interface{}{"abc", "def"}...)
	ok = s.Contains("abc")
	assert.Equal(t, ok, true, "did not match")
}

func TestSyncReverse(t *testing.T) {

	data := []struct{ in, out []interface{} }{
		{[]interface{}{}, []interface{}{}},
		{[]interface{}{"c"}, []interface{}{"c"}},
		{[]interface{}{"a", "b"}, []interface{}{"b", "a"}},
		{[]interface{}{"a", "b", "c"}, []interface{}{"c", "b", "a"}},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		s.Reverse()
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestSyncContainsAny(t *testing.T) {

	data := []struct{ src, tgt, out []interface{} }{
		{[]interface{}{}, []interface{}{}, []interface{}{}},
		{[]interface{}{"c"}, []interface{}{"c"}, []interface{}{"c"}},
		{[]interface{}{"a", "b"}, []interface{}{"b", "c", "d"}, []interface{}{"b"}},
		{[]interface{}{"a", "b", "c"}, []interface{}{"b", "c", "d"}, []interface{}{"b", "c"}},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.src...)
		out, ok := s.ContainsAny(exp.tgt)
		if ok && !reflect.DeepEqual(out, exp.out) {
			t.Fatalf("%q didn't match %q\n", out, exp.out)
		}
		s.Clear()
	}
}

func TestSyncInsert(t *testing.T) {

	data := []struct {
		in, ins, out []interface{}
		index        int
	}{
		{[]interface{}{}, []interface{}{}, []interface{}{}, 0},
		{[]interface{}{}, []interface{}{}, []interface{}{}, 22},
		{[]interface{}{"a"}, []interface{}{"b"}, []interface{}{"a", "b"}, 1},
		{[]interface{}{"b"}, []interface{}{"a"}, []interface{}{"a", "b"}, 0},
		{[]interface{}{"a", "c"}, []interface{}{"b"}, []interface{}{"a", "b", "c"}, 1},
		{[]interface{}{"a", "d"}, []interface{}{"b", "c"}, []interface{}{"a", "b", "c", "d"}, 1},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		s.Insert(exp.ins, exp.index)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestAppendNotNil(t *testing.T) {

	data := []struct {
		in, app, out []interface{}
	}{
		{[]interface{}{}, []interface{}{}, []interface{}{}},
		{[]interface{}{}, []interface{}{"a", nil}, []interface{}{"a"}},
		{[]interface{}{"a"}, []interface{}{"b"}, []interface{}{"a", "b"}},
		{[]interface{}{"b"}, []interface{}{"a", nil, nil, "c"}, []interface{}{"b", "a", "c"}},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in...)
		s.AppendNotNil(exp.app)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}
