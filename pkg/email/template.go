package email

import (
	"bytes"
)

const (
	TemplateCaptcha = "captcha.tmpl"
)

// ParseTemplate parse specified named template with given data
func (s *Sender) ParseTemplate(name string, data any) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := s.template.ExecuteTemplate(buffer, name, data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
