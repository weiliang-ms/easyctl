package harbor

import (
	"io"
	"net/http"
)

func get(url string, body io.Reader, user string, password string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(user, password)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func post(url string, body io.Reader, user string, password string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(user, password)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func put(url string, body io.Reader, user string, password string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(user, password)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}
