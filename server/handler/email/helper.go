package email

import "github.com/wneessen/go-mail"

// Message represents an email message
type Message struct {
	ContentType mail.ContentType `mapstructure:"contentType"`
	From        string           `mapstructure:"from"`
	To          []string         `mapstructure:"to"`
	CC          []string         `mapstructure:"cc"`
	Bcc         []string         `mapstructure:"bcc"`
	Subject     string           `mapstructure:"subject"`
	Body        string           `mapstructure:"body"`
}

// BuildMail build *mail.Msg from Message
func BuildMail(msg Message) (*mail.Msg, error) {
	mailMsg := mail.NewMsg()
	steps := []error{
		mailMsg.From(msg.From),
		mailMsg.To(msg.To...),
		mailMsg.Cc(msg.CC...),
		mailMsg.Bcc(msg.Bcc...),
	}
	for _, err := range steps {
		if err != nil {
			return nil, err
		}
	}
	mailMsg.Subject(msg.Subject)
	mailMsg.SetBodyString(msg.ContentType, msg.Body)
	return mailMsg, nil
}
