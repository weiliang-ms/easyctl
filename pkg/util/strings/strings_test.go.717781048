package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitIfContain(t *testing.T) {

	s, err := SplitIfContain("1:3", []string{":", "-", ".."})
	assert.Nil(t, err)
	assert.Equal(t, []string{"1", "3"}, s)

	s, err = SplitIfContain("1+2", []string{":", "-", ".."})
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(s))
	assert.Errorf(t, err, "1+2分割符不在[: - ..]内")
}

func TestSubSlash(t *testing.T) {
	assert.Equal(t, []string{"d:", "ddd", "1.txt"}, SubSlash("d:\\ddd\\1.txt"))
	assert.Equal(t, []string{"", "root", "ddd", "1.txt"}, SubSlash("/root/ddd/1.txt"))
	assert.Equal(t, []string{".", "1.txt"}, SubSlash("./1.txt"))
	assert.Equal(t, []string{".", "ddd", "redis.tar"}, SubSlash("./ddd\\redis.tar"))
	assert.Equal(t, []string{"redis.tar"}, SubSlash("redis.tar"))
}

func TestTrimPrefixAndSuffix(t *testing.T) {
	assert.Equal(t, "xxx", TrimPrefixAndSuffix("axxxa", "a"))
	assert.Equal(t, "xxxa", TrimPrefixAndSuffix("xxxa", "a"))
}

func TestSubFileName(t *testing.T) {
	assert.Equal(t, "nginx.tar.gz", SubFileName("/root/nginx.tar.gz"))
	assert.Equal(t, "nginx.tar.gz", SubFileName("C:\\root\\nginx.tar.gz"))
	assert.Equal(t, "redis.tar.gz", SubFileName("redis.tar.gz"))
	assert.Equal(t, "aaa.tar.gz", SubFileName("./aaa.tar.gz"))
	assert.Equal(t, "ddd.tar.gz", SubFileName(".\\ddd.tar.gz"))
}

func TestGetMemoryBytes(t *testing.T) {
	var v int64
	var err error
	v , err = GetMemoryBytes("1MB")
	assert.Equal(t, int64(1 * 1024 * 1024), v)
	assert.Equal(t, nil, err)

	v , err = GetMemoryBytes("1GB")
	assert.Equal(t, int64(1 * 1024 * 1024 * 1024), v)
	assert.Equal(t, nil, err)

	v , err = GetMemoryBytes("1Gb")
	assert.Equal(t, int64(0), v)
	assert.NotNil(t, err)
}
