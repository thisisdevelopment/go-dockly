package xhelper

import (
	"encoding/json"
)

// FromJSON decodes json into its expected struct
func (h *Helper) FromJSON(b []byte, expected interface{}) {
	err := json.Unmarshal(b, expected)
	h.suite.Require().NoError(err)
}

// ToJSON encodes a struct into serialized json
func (h *Helper) ToJSON(obj interface{}) []byte {
	bytes, err := json.Marshal(obj)

	h.suite.Require().NoError(err)

	return append(bytes, 0xa)
}
