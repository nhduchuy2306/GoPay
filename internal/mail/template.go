package mail

import (
	"bytes"
	"html/template"
)

func ParseTemplate(path string, data any) (string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}
	return body.String(), nil
}
