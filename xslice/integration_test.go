package xslice_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/thisisdevelopment/go-dockly/xslice"
)

var syncslice = xslice.NewSyncSlice([]interface{}{"one", "two"}...)

func TestShift(t *testing.T) {
	var res = <-syncslice.Shift()
	spew.Dump(res)
	assert.Equal(t, res.Val, "one")
}

func TestInsert(t *testing.T) {
	syncslice.Append("three", "four")
	assert.GreaterOrEqual(t, syncslice.Len(), 3)
}
