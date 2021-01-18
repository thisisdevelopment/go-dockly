package xslice_test

import (
	"reflect"
	"testing"

	"github.com/thisisdevelopment/go-dockly/xslice"
)

func TestUniqStr(t *testing.T) {

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
