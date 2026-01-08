package xhttpclient

import (
	"io"
)

func (c *Client) readResponse(b io.Reader, result any) error {
	switch t := result.(type) {
	case io.Writer:
		_, err := io.Copy(t, b)
		return err
	case *[]byte:
		body, err := io.ReadAll(b)
		if err != nil {
			return err
		}
		c.logBody(body)
		*t = body
	default:
		body, err := io.ReadAll(b)
		if err != nil {
			return err
		}
		c.logBody(body)
		if c.unmarshal != nil {
			return c.unmarshal(body, result)
		}
	}

	return nil
}
