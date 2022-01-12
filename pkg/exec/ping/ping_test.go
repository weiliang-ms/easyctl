package ping

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/weiliang-ms/easyctl/pkg/exec/ping/mocks"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"os"
	"testing"
	"time"
)

var (
	mockIP1 = "192.168.1.1"
	mockIP2 = "192.168.1.2"
	mockIP3 = "192.168.1.3"

	mockAddress1 = fmt.Sprintf("%s:22", mockIP1)
	mockAddress2 = fmt.Sprintf("%s:22", mockIP2)
	mockAddress3 = fmt.Sprintf("%s:22", mockIP3)

	mockIPs = []string{
		mockIP1,
		mockIP2,
		mockIP3,
	}
	mockAddresses = []string{
		mockAddress1,
		mockAddress2,
		mockAddress3,
	}

	mockErr          = fmt.Errorf("mock error")
	mockLogger       = logrus.New()
	mockCheckCount   = 1
	mockCheckTimeout = time.Millisecond
)

func Test_OnlyPingSuccess_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}

	m := Manger{
		Handler:      mockInterface,
		CheckCount:   mockCheckCount,
		CheckTimeout: mockCheckTimeout,
	}
	mockInterface.On("ICMP",
		mockIP1, mockCheckCount, mockCheckTimeout, true).
		Return(nil)
	mockInterface.On("ICMP",
		mockIP2, mockCheckCount, mockCheckTimeout, true).
		Return(nil)
	mockInterface.On("ICMP",
		mockIP3, mockCheckCount, mockCheckTimeout, true).
		Return(nil)

	err := m.getSurviveList(mockIPs, mockLogger)
	require.Nil(t, err)
}

func Test_OnlyPingErr_Mock(t *testing.T) {

	defer os.Remove("server-list.txt")
	mockInterface := &mocks.HandlerInterface{}

	m := Manger{
		Handler:      mockInterface,
		CheckCount:   mockCheckCount,
		CheckTimeout: mockCheckTimeout,
	}
	mockInterface.On("ICMP",
		mockIP1, mockCheckCount, mockCheckTimeout, true).
		Return(mockErr)
	mockInterface.On("ICMP",
		mockIP2, mockCheckCount, mockCheckTimeout, true).
		Return(mockErr)
	mockInterface.On("ICMP",
		mockIP3, mockCheckCount, mockCheckTimeout, true).
		Return(mockErr)

	err := m.getSurviveList(mockIPs, mockLogger)
	require.Nil(t, err)
}

func Test_PingAndTelnetSuccess_Mock(t *testing.T) {

	defer os.Remove("server-list.txt")
	mockInterface := &mocks.HandlerInterface{}

	m := Manger{
		Handler:      mockInterface,
		CheckCount:   mockCheckCount,
		CheckTimeout: mockCheckTimeout,
	}
	mockInterface.On("ICMP",
		mockIP1, mockCheckCount, mockCheckTimeout, true).
		Return(nil)
	mockInterface.On("Telnet",
		"tcp", mockAddress1, mockCheckTimeout).
		Return(nil)

	mockInterface.On("ICMP",
		mockIP2, mockCheckCount, mockCheckTimeout, true).
		Return(nil)
	mockInterface.On("Telnet",
		"tcp", mockAddress2, mockCheckTimeout).
		Return(nil)

	mockInterface.On("ICMP",
		mockIP3, mockCheckCount, mockCheckTimeout, true).
		Return(nil)
	mockInterface.On("Telnet",
		"tcp", mockAddress3, mockCheckTimeout).
		Return(nil)

	err := m.getSurviveList(mockAddresses, mockLogger)
	require.Nil(t, err)
}

func Test_ParsePingItems(t *testing.T) {

	content := `
ping:
  - address: 192.168.1.
    start: 1
    end: 3
    #port: 22
  - address: 192.168.555.
    start: 1
    end: 40
    #port: 22
  - address: 192.168.3.
    start: 256
    end: 22
    #port: 22
  - address: 192.168.4.
    start: 1
    end: 256
    #port: 22
  - address: 192.168.5.
    start: -3
    end: 256
    #port: 22
  - address: 192.168.5.
    start: 30
    end: 20
    #port: 22
`

	re, err := ParsePingItems([]byte(content), mockLogger)
	require.Nil(t, err)
	require.Equal(t, 3, len(re))
}

func Test_ParsePingItems_ParseErrCase(t *testing.T) {
	defer os.Remove("server-list.txt")
	content := `
ping:
  - address: 192.168.1.
    - start: 1
    end: 3
    #port: 22
`

	re, err := ParsePingItems([]byte(content), mockLogger)
	require.NotNil(t, err)
	require.Equal(t, 0, len(re))
}

func TestPing(t *testing.T) {
	defer os.Remove("server-list.txt")
	content := `
ping:
 - address: 1.1.1.1
   start: 1
   end: 3
   #port: 22
`
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	Run(command.OperationItem{B: []byte(content), Logger: mockLogger})
}

func Test_PingParseConfigErr(t *testing.T) {
	content := `
ping:
- address: 1.1.1.1
  - start: 1
  end: 255
  #port: 22
`
	err := Run(command.OperationItem{B: []byte(content), Logger: logrus.New()})
	require.NotNil(t, err)
}

func Test_getHandlerInterface(t *testing.T) {

	var h HandlerInterface
	r := getHandlerInterface(h)
	require.Equal(t, new(Handler), r)

	h2 := &mocks.HandlerInterface{}
	r2 := getHandlerInterface(h2)
	require.Equal(t, h2, r2)
}
