package strings

import (
	"fmt"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	ips := IPS{"192.168.174.10", "192.168.174.4", "1.1.1.1", "192.168.174.1", "192.168.1.10"}
	sort.Sort(ips)
	fmt.Printf("%#v", ips)
}
