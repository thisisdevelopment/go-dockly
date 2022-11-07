package xhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func ExtractHeader(key string, r *http.Request) (value string, err error) {
	value = r.Header.Get(key)
	if value == "" {
		return value, fmt.Errorf("no %s provided in header", key)
	}

	return value, nil
}

func ExtractIntParam(key string, r *http.Request) (value int, err error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return value, fmt.Errorf("no %s provided in url", key)
	}

	return strconv.Atoi(param)
}

func ExtractParam(key string, r *http.Request) (value string, err error) {
	value = r.URL.Query().Get(key)
	if value == "" {
		return value, fmt.Errorf("no %s provided in url", key)
	}

	return value, nil
}

func ExtractBody(body io.ReadCloser, expected interface{}) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, expected)
}
