package runner

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

// 测试servers排序
func TestSortServers(t *testing.T) {
	servers := InternelServersSlice{
		ServerInternal{
			Host: "10.10.10.3",
		},
		ServerInternal{
			Host: "10.10.10.2",
		},
		ServerInternal{
			Host: "10.10.10.1",
		},
	}

	sort.Sort(servers)

	expect := InternelServersSlice{
		ServerInternal{
			Host: "10.10.10.1",
		},
		ServerInternal{
			Host: "10.10.10.2",
		},
		ServerInternal{
			Host: "10.10.10.3",
		},
	}

	assert.Equal(t, expect, servers)
}

// 测试ShellResult排序
func TestSortShellResult(t *testing.T) {
	result := ShellResultSlice{
		ShellResult{
			Host: "10.10.10.3",
		},
		ShellResult{
			Host: "10.10.10.2",
		},
		ShellResult{
			Host: "10.10.10.1",
		},
	}

	sort.Sort(result)

	expect := ShellResultSlice{
		ShellResult{
			Host: "10.10.10.1",
		},
		ShellResult{
			Host: "10.10.10.2",
		},
		ShellResult{
			Host: "10.10.10.3",
		},
	}

	assert.Equal(t, expect, result)
}
