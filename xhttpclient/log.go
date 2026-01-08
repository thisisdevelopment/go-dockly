package xhttpclient

func (c *Client) log(format string, v ...interface{}) {
	if c.verbose && c.logf != nil {
		c.logf(format, v...)
	}
}

func (c *Client) logBody(body []byte) {
	n := len(body)
	postfix := ""
	if n > 48 {
		n = 48
		postfix = "..."
	}
	c.log("response body: %s%s", string(body)[:n], postfix)
}
