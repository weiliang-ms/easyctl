package exec

import "testing"

func TestRun(t *testing.T) {
	err := Run("executor.yaml", 0)
	if err != nil {
		t.Error(err)
	}
}
