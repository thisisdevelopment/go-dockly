package xslice_test

import (
	"reflect"
	"testing"

	"github.com/thisisdevelopment/go-dockly/xslice"
	//"github.com/stretchr/testify/assert"
)

// make a Test function with a goroutine (=integration test), call couple of SyncSlice functions concurrently. + create big slices for testing.
// do not sort : find alternative to map.

/*
func AllTest(t *testing.T) {


	expdata := make(map[string]struct{}, 0)


	uniqdata := []struct{ in, out []interface{} }{
		{[]interface{}{}, []interface{}{}},
		{[]interface{}{"", "", ""}, []interface{}{""}},
		{[]interface{}{"a", "a"}, []interface{}{"a"}},
		{[]interface{}{"a", "b", "a"}, []interface{}{"a", "b"}},
		{[]interface{}{"a", "b", "a", "b"}, []interface{}{"a", "b"}},
		{[]interface{}{"a", "b", "b", "a", "b"}, []interface{}{"a", "b"}},
		{[]interface{}{"a", "a", "b", "b", "a", "b"}, []interface{}{"a", "b"}},
		{[]interface{}{"a", "b", "c", "a", "b", "c"}, []interface{}{"a", "b", "c"}},
	}

	cutdata := []struct {
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
		go func() {
			s.Append(exp.in)
			s.Uniq()


		}


}
*/

func TestUniq(t *testing.T) {

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
		s.Append(exp.in)
		s.Uniq()
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestCut(t *testing.T) {

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
		s.Append(exp.in)
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

func TestStrip(t *testing.T) {

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
		s.Append(exp.in)
		s.Strip(exp.val)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestDel(t *testing.T) {

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
		s.Append(exp.in)
		s.Del(exp.index)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestPop(t *testing.T) {

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
		s.Append(exp.in)
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

func TestShift(t *testing.T) {

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
		s.Append(exp.in)
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

func TestUnShift(t *testing.T) {

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
		s.Append(exp.in)
		s.UnShift(exp.unshift)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestFilter(t *testing.T) {

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
		s.Append(exp.in)
		s.Filter(exp.filter)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

/*
func TestContains(t *testing.T) {
  var ok bool
  data := []interface{}{"abc","def"}
  ok = xslice.Contains(data,"abc")
  assert.Equal(t, ok, true, "did not match")
}
*/

func TestReverse(t *testing.T) {

	data := []struct{ in, out []interface{} }{
		{[]interface{}{}, []interface{}{}},
		{[]interface{}{"c"}, []interface{}{"c"}},
		{[]interface{}{"a", "b"}, []interface{}{"b", "a"}},
		{[]interface{}{"a", "b", "c"}, []interface{}{"c", "b", "a"}},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in)
		s.Reverse()
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}

func TestContainsAny(t *testing.T) {

	data := []struct{ src, tgt, out []interface{} }{
		{[]interface{}{}, []interface{}{}, []interface{}{}},
		{[]interface{}{"c"}, []interface{}{"c"}, []interface{}{"c"}},
		{[]interface{}{"a", "b"}, []interface{}{"b", "c", "d"}, []interface{}{"b"}},
		{[]interface{}{"a", "b", "c"}, []interface{}{"b", "c", "d"}, []interface{}{"b", "c"}},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.src)
		out, ok := s.ContainsAny(exp.tgt)
		if ok && !reflect.DeepEqual(out, exp.out) {
			t.Fatalf("%q didn't match %q\n", out, exp.out)
		}
		s.Clear()
	}
}

func TestInsert(t *testing.T) {

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
		s.Append(exp.in)
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
		{[]interface{}{}, []interface{}{"a", ""}, []interface{}{"a"}},
		{[]interface{}{"a"}, []interface{}{"b"}, []interface{}{"a", "b"}},
		{[]interface{}{"b"}, []interface{}{"a", "", "", "c"}, []interface{}{"a", "b", "c"}},
	}

	s := xslice.NewSyncSlice()
	for _, exp := range data {
		s.Append(exp.in)
		s.AppendNotNil(exp.app)
		for item := range s.Iter() {
			if !reflect.DeepEqual(item.Val, exp.out[item.Idx]) {
				t.Fatalf("%q didn't match %q\n", item.Val, exp.out[item.Idx])
			}
		}
		s.Clear()
	}
}
