package types

import "github.com/wneessen/go-mail"

// EmailBody represents an email message body
type EmailBody struct {
	ContentType mail.ContentType `mapstructure:"contentType"`
	From        string           `mapstructure:"from"`
	To          []string         `mapstructure:"to"`
	CC          []string         `mapstructure:"cc"`
	Bcc         []string         `mapstructure:"bcc"`
	Subject     string           `mapstructure:"subject"`
	Body        string           `mapstructure:"body"`
}
