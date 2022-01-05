package format

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Person struct {
	Name string
	Age  int
	Sex  string
}

func TestObject(t *testing.T) {
	p := Person{
		Name: "Tom",
		Age:  16,
		Sex:  "male",
	}
	b := ObjectToJson(p)

	out := `{
	"Name": "Tom",
	"Age": 16,
	"Sex": "male"
}`
	require.Equal(t, out, b.String())
}
