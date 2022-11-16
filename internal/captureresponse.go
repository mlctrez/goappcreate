package internal

import (
	"bytes"
	"net/http"
)

var _ http.ResponseWriter = (*CapturingResponseWriter)(nil)

type CapturingResponseWriter struct {
	header     http.Header
	body       *bytes.Buffer
	statusCode int
}

func (c *CapturingResponseWriter) Header() http.Header {
	if c.header == nil {
		c.header = make(http.Header)
	}
	return c.header
}

func (c *CapturingResponseWriter) Write(bytes_ []byte) (int, error) {
	if c.body == nil {
		c.body = &bytes.Buffer{}
	}
	return c.body.Write(bytes_)
}

func (c *CapturingResponseWriter) WriteHeader(statusCode int) {
	c.statusCode = statusCode
}
