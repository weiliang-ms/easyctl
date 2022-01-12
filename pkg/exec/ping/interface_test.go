package ping

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_HandlerInterface(t *testing.T) {
	h := Handler{}
	require.NotNil(t, h.ICMP("1.1.1.1", 1, time.Millisecond, true))
	require.NotNil(t, h.Telnet("tcp", "1.1.1.1", time.Millisecond))

	// invalid ip address case
	require.NotNil(t, h.ICMP("444.444.44.4", 1, time.Millisecond, true))
	require.NotNil(t, h.ICMP("", 1, time.Millisecond, true))
}
