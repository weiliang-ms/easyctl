package chart

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/weiliang-ms/easyctl/pkg/util/request"
	"testing"
)

func Test_GetChartByte(t *testing.T) {
	var h Handler
	h.GetChartByte(request.HTTPRequestItem{
		Url:      "",
		Method:   "",
		Body:     nil,
		Headers:  nil,
		User:     "",
		Timeout:  0,
		Password: "",
		Mock:     false,
	})
}

func Test_GetChartList(t *testing.T) {
	var h Handler
	h.GetChartList(request.HTTPRequestItem{
		Url:      "",
		Method:   "",
		Body:     nil,
		Headers:  nil,
		User:     "",
		Timeout:  0,
		Password: "",
		Mock:     false,
	})
}

func Test_ApiListChartsUrl(t *testing.T) {
	ep := "1.1.1.1:80"
	re := ApiListChartsUrl("http", ep)
	require.Equal(t, fmt.Sprintf("http://%s/api/chartrepo/charts/charts", ep), re)
}

func Test_ApiChartBytesUrl(t *testing.T) {
	ep := "1.1.1.1:80"
	re := ApiChartBytesUrl("http", ep, "nginx")
	require.Equal(t, fmt.Sprintf("http://%s/api/chartrepo/charts/charts/nginx", ep), re)
}
