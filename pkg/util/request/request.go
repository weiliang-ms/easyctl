package request

import (
	"io"
	"time"
)

type HTTPRequestItem struct {
	Url      string
	Method   string
	Body     io.Reader
	Headers  map[string]string
	User     string
	Timeout  time.Duration
	Password string
	Mock     bool
}
