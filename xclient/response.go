package xclient

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// readResponse() will try to unmarshal the response body into the
// desired result interface or return an error
func (cli *Client) readResponse(b io.ReadCloser, result interface{}) error {
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return errors.Wrap(err, "reading response failed")
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return errors.Wrapf(err, "unmarshal response into expected interface failed for: %s", string(body))
	}

	return nil
}
