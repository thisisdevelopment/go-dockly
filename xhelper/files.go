package xhelper

import (
	"io/ioutil"
)

func (h *Helper) BytesFromFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)

	h.suite.Require().NoError(err)

	return bytes
}
