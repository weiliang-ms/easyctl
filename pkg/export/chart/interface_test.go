/*
	MIT License

Copyright (c) 2022 xzx.weiliang

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

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
