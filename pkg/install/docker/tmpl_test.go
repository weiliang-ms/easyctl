package docker

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"testing"
)

func TestDockerConfigTmpl(t *testing.T) {
	shell, err := tmplutil.Render(DockerConfigTmpl, tmplutil.TmplRenderData{
		"DataPath":           "/data/lib/docker",
		"Mirrors":            "\"harbor.aa.io\",\"harbor.bb.io\"",
		"InsecureRegistries": "\"xxx.xxx.io\",\"xxx.yyy.io\"",
	})
	assert.Equal(t, nil, err)
	fmt.Println(shell)
}
