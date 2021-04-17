package util

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

type Data map[string]interface{}

// Render text template with given `variables` Render-context
func Render(tmpl *template.Template, variables map[string]interface{}) (string, error) {

	var buf strings.Builder

	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", errors.Wrap(err, "Failed to render template")
	}
	return buf.String(), nil
}

func OpenPortCmd(port int) string {
	return func(p int) string {
		reload := "firewall-cmd --reload"
		open := fmt.Sprintf("\nfirewall-cmd --zone=public --add-port=%d/tcp --permanent", port)
		return fmt.Sprintf("%s && %s\n", open, reload)
	}(port)
}

func Get(url string, username string, password string) []byte {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
