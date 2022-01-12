package request

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_DoRequest(t *testing.T) {
	// test DoRequest
	b, err := DoRequest(HTTPRequestItem{Method: "！！！"})
	require.Nil(t, b)
	require.Equal(t, "net/http: invalid method \"！！！\"", err.Error())

	// do request error case
	b2, err2 := DoRequest(HTTPRequestItem{Mock: false, Timeout: time.Millisecond, Url: "http://www.baidu.com"})
	require.NotNil(t, err2)
	require.Equal(t, "", string(b2))

	// test DoRequest
	defer func() {
		r := recover()
		require.Equal(t, "runtime error: invalid memory address or nil pointer dereference", fmt.Sprintf("%s", r))
	}()
	_, _ = DoRequest(HTTPRequestItem{Mock: true})

}
