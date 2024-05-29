package xclient

import (
	"encoding/json"
	"io"
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// readResponse() will try to unmarshal the response body into the
// desired result interface or return an error
func (cli *Client) readResponse(b io.ReadCloser, result interface{}) error {
	body, err := io.ReadAll(b)
	if err != nil {
		return errors.Wrap(err, "reading response failed")
	}

	switch result.(type) {
	case *[]byte:
		// assign the raw byte slice of body to the results interface as is
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(body))
	default:
		if cli.config.UseJsoniter {
			if err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, result); err != nil {
				return errors.Wrapf(err, "unmarshal response failed: %s", string(body))
			}
		} else {
			if err = json.Unmarshal(body, result); err != nil {
				return errors.Wrapf(err, "unmarshal response failed: %s", string(body))
			}
		}

	}

	return nil
}
