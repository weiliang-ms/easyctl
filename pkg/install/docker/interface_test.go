package docker

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Detect(t *testing.T) {
	var h Handler
	require.Nil(t, h.Detect("", mockServer, mockLocal, mockLogger, mockTimeout))
}

func Test_Prune(t *testing.T) {
	var h Handler
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	h.Prune(mockServer, mockLocal, mockLogger, mockTimeout)
}

func Test_Install(t *testing.T) {
	var h Handler
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	h.Install("", mockServer, mockLocal, mockLogger, mockTimeout)
}

func Test_Boot(t *testing.T) {
	var h Handler
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	h.Boot(mockServer, mockLocal, mockLogger, mockTimeout)
}

func Test_SetConfig(t *testing.T) {
	var h Handler
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	h.SetConfig("", mockServer, mockLocal, mockLogger, mockTimeout)
}

func Test_SetSystemd(t *testing.T) {
	var h Handler
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	h.SetSystemd("", mockServer, mockLocal, mockLogger, mockTimeout)
}

func Test_SetUpRuntime(t *testing.T) {
	var h Handler
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	h.SetUpRuntime("", mockServer, mockLocal, mockLogger, mockTimeout)
}

func Test_HandPackage(t *testing.T) {
	var h Handler
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	h.HandPackage(mockServer, "", mockLocal, mockLogger, mockTimeout)
}

func Test_HandPackageLocal(t *testing.T) {
	var h Handler
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	h.HandPackage(mockServer, "", true, mockLogger, mockTimeout)
}

func Test_Exec(t *testing.T) {
	var h Handler
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()
	h.Exec("", mockServer, true, mockLogger, mockTimeout)
}
