/*
	MIT License

Copyright (c) 2020 xzx.weiliang

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

	expect2 := InternelServersSlice{
		ServerInternal{
			Host: "10.10.10.1",
		},
		ServerInternal{
			Host: "10.10.10.1",
		},
	}

	assert.Equal(t, true, expect2.Less(0, 1))
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

	result2 := ShellResultSlice{
		ShellResult{
			Host: "10.10.10.1",
		},
		ShellResult{
			Host: "10.10.10.1",
		},
	}

	assert.Equal(t, true, result2.Less(0, 1))
}
