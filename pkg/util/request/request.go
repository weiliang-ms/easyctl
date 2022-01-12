package request

import (
	"io"
	"net/http"
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

func DoRequest(httpRequestItem HTTPRequestItem) ([]byte, error) {

	req, err := http.NewRequest(httpRequestItem.Method, httpRequestItem.Url, httpRequestItem.Body)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: httpRequestItem.Timeout}

	req.SetBasicAuth(httpRequestItem.User, httpRequestItem.Password)

	resp, err := client.Do(req)

	if err != nil && !httpRequestItem.Mock {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}
