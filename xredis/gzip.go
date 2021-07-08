package xredis

import (
	"bytes"
	"compress/gzip"

	"github.com/pkg/errors"
)

func (c *Redis) gzip(val []byte) ([]byte, error) {
	var buf bytes.Buffer
	var w = gzip.NewWriter(&buf)

	pos, err := w.Write(val)
	if err != nil {
		return nil, errors.Wrapf(err, "gzip content for pos %d", pos)
	}

	err = w.Close()
	if err != nil {
		return nil, errors.Wrap(err, "close gzip content writer")
	}

	return buf.Bytes(), nil
}

func (c *Redis) gunzip(val []byte) ([]byte, error) {
	var b = bytes.NewBuffer(val)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, errors.Wrap(err, "open gzip content reader")
	}

	var res bytes.Buffer
	pos, err := res.ReadFrom(r)
	if err != nil {
		return nil, errors.Wrapf(err, "read gzip content pos %d", pos)
	}

	err = r.Close()
	if err != nil {
		return nil, errors.Wrap(err, "close gzip content reader")
	}

	return res.Bytes(), nil
}
