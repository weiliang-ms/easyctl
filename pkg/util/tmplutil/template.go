package tmplutil

import (
	"github.com/pkg/errors"
	"strings"
	"text/template"
)

// TmplRenderData 渲染模板的数据集
type TmplRenderData map[string]interface{}

// Render text template with given `variables` Render-context
func Render(tmpl *template.Template, variables map[string]interface{}) (string, error) {

	var buf strings.Builder

	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", errors.Wrap(err, "Failed to render template")
	}
	return buf.String(), nil
}

func RenderPanicErr(tmpl *template.Template, variables map[string]interface{}) string {
	var buf strings.Builder
	if err := tmpl.Execute(&buf, variables); err != nil {
		panic(err)
	}
	return buf.String()
}
