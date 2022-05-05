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
	"github.com/weiliang-ms/easyctl/pkg/util/request"
)

//go:generate mockery --name=HandlerInterface
type Handler struct{}

func ApiListChartsUrl(schema, endpoint string) string {
	return fmt.Sprintf("%s://%s/api/chartrepo/charts/charts", schema, endpoint)
}

func ApiChartBytesUrl(schema, endpoint, name string) string {
	return fmt.Sprintf("%s://%s/api/chartrepo/charts/charts/%s", schema, endpoint, name)
}

type HandlerInterface interface {
	DoRequest(httpRequestItem request.HTTPRequestItem) ([]byte, error)
	GetChartByte(httpRequestItem request.HTTPRequestItem) ([]byte, error)
	GetChartList(httpRequestItem request.HTTPRequestItem) ([]byte, error)
}

func (h Handler) DoRequest(httpRequestItem request.HTTPRequestItem) ([]byte, error) {
	return request.DoRequest(httpRequestItem)
}

func (h Handler) GetChartList(httpRequestItem request.HTTPRequestItem) ([]byte, error) {
	return h.DoRequest(httpRequestItem)
}

func (h Handler) GetChartByte(httpRequestItem request.HTTPRequestItem) ([]byte, error) {
	return h.DoRequest(httpRequestItem)
}
