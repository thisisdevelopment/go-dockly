package xhelper

import "os"

func (h *Helper) BytesFromFile(path string) []byte {
	bytes, err := os.ReadFile(path)

	h.suite.Require().NoError(err)

	return bytes
}
