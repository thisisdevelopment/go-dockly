package xslice_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thisisdevelopment/go-dockly/xslice"
)

var syncslice = xslice.NewSyncSlice([]interface{}{"one", "two"})

func TestShift(t *testing.T) {

	res := <-syncslice.Shift()
	assert.Equal(t, res.Val.(string), "one")
}
func TestInsert(t *testing.T) {

	syncslice.Append([]interface{}{"three", "four"})
	assert.GreaterOrEqual(t, syncslice.Len(), 3)
}
