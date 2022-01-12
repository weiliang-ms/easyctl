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
