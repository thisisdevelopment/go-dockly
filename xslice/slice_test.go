package xslice_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thisisdevelopment/go-dockly/xslice"
)

func TestUniq(t *testing.T) {

	data := []struct{ in, out []string }{
		{[]string{}, []string{}},
		{[]string{"", "", ""}, []string{""}},
		{[]string{"a", "a"}, []string{"a"}},
		{[]string{"a", "b", "a"}, []string{"a", "b"}},
		{[]string{"a", "b", "a", "b"}, []string{"a", "b"}},
		{[]string{"a", "b", "b", "a", "b"}, []string{"a", "b"}},
		{[]string{"a", "a", "b", "b", "a", "b"}, []string{"a", "b"}},
		{[]string{"a", "b", "c", "a", "b", "c"}, []string{"a", "b", "c"}},
	}

	for _, exp := range data {
		res := xslice.Uniq(exp.in)
		if !reflect.DeepEqual(res, exp.out) {
			t.Fatalf("%q didn't match %q\n", res, exp.out)
		}
	}
}

func TestCut(t *testing.T) {

	data := []struct {
		in, out    []string
		start, end int
	}{
		{[]string{}, []string{}, 0, 0},
		{[]string{"", "", ""}, []string{""}, 1, 3},
		{[]string{"a", "a"}, []string{"a"}, 1, 2},
		{[]string{"a", "b", "a"}, []string{"a", "b"}, 2, 3},
		{[]string{"a", "b", "a", "b"}, []string{"a", "b"}, 2, 4},
		{[]string{"a", "a", "b", "b", "a", "b"}, []string{"a", "b"}, 0, 4},
		{[]string{"a", "b", "c", "a", "b", "c"}, []string{"a", "b", "c"}, 3, 6},
	}

	for _, exp := range data {
		res, _ := xslice.Cut(exp.in, exp.start, exp.end)
		if !reflect.DeepEqual(res, exp.out) {
			t.Fatalf("%q didn't match %q\n", res, exp.out)
		}
	}
}

func TestStrip(t *testing.T) {

	data := []struct {
		in, out []string
		val     string
	}{
		{[]string{}, []string{}, "blah"},
		{[]string{"", "", ""}, []string{}, ""},
		{[]string{"a", "a"}, []string{"a", "a"}, "b"},
		{[]string{"a", "b", "a"}, []string{"b"}, "a"},
		{[]string{"c", "c", "c"}, []string{}, "c"},
	}

	for _, exp := range data {
		res := xslice.Strip(exp.in, exp.val)
		if !reflect.DeepEqual(res, exp.out) {
			t.Fatalf("%q didn't match %q\n", res, exp.out)
		}
	}
}

func TestDel(t *testing.T) {

	data := []struct {
		in, out []string
		index   int
	}{
		{[]string{}, []string{}, 22},
		{[]string{"", "", ""}, []string{"", ""}, 1},
		{[]string{"a", "a"}, []string{"a"}, 0},
		{[]string{"a", "b", "a"}, []string{"a", "a"}, 1},
		{[]string{"a", "b", "c"}, []string{"a", "b"}, 2},
	}

	for _, exp := range data {
		res, _ := xslice.Del(exp.in, exp.index)
		if !reflect.DeepEqual(res, exp.out) {
			t.Fatalf("%q didn't match %q\n", res, exp.out)
		}
	}
}

func TestPop(t *testing.T) {

	data := []struct {
		in, out []string
		pop     string
	}{
		{[]string{}, []string{}, ""},
		{[]string{"a"}, []string{}, "a"},
		{[]string{"a", "b"}, []string{"a"}, "b"},
	}

	for _, exp := range data {
		pop, out, _ := xslice.Pop(exp.in)
		if pop != exp.pop {
			t.Fatalf("%q didn't match %q\n", pop, exp.pop)
		}
		if !reflect.DeepEqual(out, exp.out) {
			t.Fatalf("%q didn't match %q\n", out, exp.out)
		}
	}
}

func TestShift(t *testing.T) {

	data := []struct {
		in, out []string
		shift   string
	}{
		{[]string{}, []string{}, ""},
		{[]string{"a"}, []string{}, "a"},
		{[]string{"a", "b"}, []string{"b"}, "a"},
		{[]string{"a", "b", "c"}, []string{"b", "c"}, "a"},
	}

	for _, exp := range data {
		shift, out, _ := xslice.Shift(exp.in)
		if shift != exp.shift {
			t.Fatalf("%q didn't match %q\n", shift, exp.shift)
		}
		if !reflect.DeepEqual(out, exp.out) {
			t.Fatalf("%q didn't match %q\n", out, exp.out)
		}
	}
}

func TestUnShift(t *testing.T) {

	data := []struct {
		in, out []string
		unshift string
	}{
		{[]string{}, []string{"a"}, "a"},
		{[]string{"b"}, []string{"a", "b"}, "a"},
		{[]string{"b", "c"}, []string{"a", "b", "c"}, "a"},
	}

	for _, exp := range data {
		out := xslice.UnShift(exp.in, exp.unshift)
		if !reflect.DeepEqual(out, exp.out) {
			t.Fatalf("%q didn't match %q\n", out, exp.out)
		}
	}
}

func TestFilter(t *testing.T) {

	data := []struct {
		in, out []string
		filter  string
	}{
		{[]string{}, []string{}, "a"},
		{[]string{"c"}, []string{}, "b"},
		{[]string{"c"}, []string{"c"}, "c"},
		{[]string{"a", "b", "c"}, []string{"b"}, "b"},
	}

	for _, exp := range data {
		out := xslice.Filter(exp.in, exp.filter)
		if !reflect.DeepEqual(out, exp.out) {
			t.Fatalf("%q didn't match %q\n", out, exp.out)
		}
	}
}

func TestContains(t *testing.T) {
	var ok bool
	data := []string{"abc", "def"}
	ok = xslice.Contains(data, "abc")
	assert.Equal(t, ok, true, "did not match")
}

func TestReverse(t *testing.T) {

	data := []struct{ in, out []string }{
		{[]string{}, []string{}},
		{[]string{"c"}, []string{"c"}},
		{[]string{"a", "b"}, []string{"b", "a"}},
		{[]string{"a", "b", "c"}, []string{"c", "b", "a"}},
	}

	for _, exp := range data {
		out := xslice.Reverse(exp.in)
		if !reflect.DeepEqual(out, exp.out) {
			t.Fatalf("%q didn't match %q\n", out, exp.out)
		}
	}
}

func TestContainsAny(t *testing.T) {

	data := []struct{ src, tgt, out []string }{
		{[]string{}, []string{}, []string{}},
		{[]string{"c"}, []string{"c"}, []string{"c"}},
		{[]string{"a", "b"}, []string{"b", "c", "d"}, []string{"b"}},
		{[]string{"a", "b", "c"}, []string{"b", "c", "d"}, []string{"b", "c"}},
	}

	for _, exp := range data {
		out, _ := xslice.ContainsAny(exp.src, exp.tgt)
		if !reflect.DeepEqual(out, exp.out) {
			t.Fatalf("%q didn't match %q\n", out, exp.out)
		}
	}
}

func TestInsert(t *testing.T) {

	data := []struct {
		in, ins, out []string
		index        int
	}{
		{[]string{}, []string{}, []string{}, 0},
		{[]string{}, []string{}, []string{}, 22},
		{[]string{"a"}, []string{"b"}, []string{"a", "b"}, 1},
		{[]string{"b"}, []string{"a"}, []string{"a", "b"}, 0},
		{[]string{"a", "c"}, []string{"b"}, []string{"a", "b", "c"}, 1},
		{[]string{"a", "d"}, []string{"b", "c"}, []string{"a", "b", "c", "d"}, 1},
	}

	for _, exp := range data {
		out := xslice.Insert(exp.in, exp.ins, exp.index)
		if !reflect.DeepEqual(out, exp.out) {
			t.Fatalf("%q didn't match %q\n", out, exp.out)
		}
	}
}

func TestAppendNotEmpty(t *testing.T) {

	data := []struct {
		in, app, out []string
	}{
		{[]string{}, []string{}, []string{}},
		{[]string{}, []string{"a", ""}, []string{"a"}},
		{[]string{"a"}, []string{"b"}, []string{"a", "b"}},
		{[]string{"b"}, []string{"a", "", "", "c"}, []string{"b", "a", "c"}},
	}

	for _, exp := range data {
		out := xslice.AppendNotEmpty(exp.in, exp.app)
		if !reflect.DeepEqual(out, exp.out) {
			t.Fatalf("%q didn't match %q\n", out, exp.out)
		}
	}
}
