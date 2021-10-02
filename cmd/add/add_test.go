package add

import (
	"io"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {

	entity := Entity{DefaultConfig: nil}
	entity.setDefault()

	if entity.DefaultConfig == nil {
		t.Error("未配置默认值")
	}

	f, _ := os.Open("asset/config.yaml")
	b, _ := io.ReadAll(f)

	if string(entity.DefaultConfig) != string(b) {
		t.Error("配置不相等")
	}
}
