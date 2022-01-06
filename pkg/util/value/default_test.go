package value

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type Person struct {
	Name  string
	Age   int
	Age32 int32
	Age64 int64
	Bool  bool
}

func TestSetDefaultValue(t *testing.T) {
	p := Person{}
	err := SetStructDefaultValue(&p, "Name", "Tom")
	if err != nil {
		fmt.Println(err)
	}

	err = SetStructDefaultValue(&p, "Age", 22)
	if err != nil {
		fmt.Println(err)
	}

	err = SetStructDefaultValue(&p, "Age32", int32(32))
	if err != nil {
		fmt.Println(err)
	}

	err = SetStructDefaultValue(&p, "Age64", int64(64))
	if err != nil {
		fmt.Println(err)
	}

	err = SetStructDefaultValue(&p, "Bool", true)
	if err != nil {
		fmt.Println(err)
	}

	require.Equal(t, "Tom", p.Name)
	require.Equal(t, true, p.Bool)
	require.Equal(t, 22, p.Age)
	require.Equal(t, int32(32), p.Age32)
	require.Equal(t, int64(64), p.Age64)
}

func TestSetDefaultValue_ErrCase(t *testing.T) {
	var a string
	err := SetStructDefaultValue(a, "Name", "Tom")
	require.Equal(t, "必需为指针类型", err.Error())

	//m := make(map[string]string)
	err = SetStructDefaultValue(&a, "Name", "Tom")
	require.Equal(t, "必须为结构体类型", err.Error())

	d := struct {
		//A string
	}{}

	err = SetStructDefaultValue(&d, "DDD", nil)
	//require.Equal(t, "必须为结构体类型", err.Error())
}
