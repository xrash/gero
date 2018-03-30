package msg

import (
	"bytes"
	"text/template"
)

func ProcessTemplate(ts string, data interface{}) (string, error) {
	t, err := template.New("whatever").Parse(ts)
	if err != nil {
		return "", err
	}

	b := bytes.NewBuffer(nil)
	if err := t.Execute(b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
