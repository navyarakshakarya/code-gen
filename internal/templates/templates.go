package templates

import (
	"bytes"
	"text/template"
)

// Execute executes a template with the given data
func Execute(tmpl string, data interface{}) (string, error) {
	t, err := template.New("template").Funcs(GetTemplateFuncs()).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
