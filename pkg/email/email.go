package email

import (
	"context"
	"github.com/ginx-contribs/ginx-server/pkg/email/internal/templates"
	"github.com/ginx-contribs/str2bytes"
	"github.com/wneessen/go-mail"
	"html/template"
)

type Options struct {
	// smtp server host
	Host    string
	SSLPort bool
	// smtp server port
	Port     int
	Username string
	Password string
	// email template resolve dir
	TemplateDir string
}

// NewSender initialize email sender
func NewSender(options Options) (*Sender, error) {
	client, err := mail.NewClient(options.Host,
		mail.WithPort(options.Port),
		mail.WithSSLPort(true),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(options.Username),
		mail.WithPassword(options.Password),
	)
	if err != nil {
		return nil, err
	}
	// test smtp server if is available
	err = client.DialWithContext(context.Background())
	if err != nil {
		return nil, err
	}

	// parse template
	var tmpl *template.Template
	if options.TemplateDir != "" {
		customTmpl, err := template.ParseFiles(options.TemplateDir, "*.tmpl")
		if err != nil {
			return nil, err
		}
		tmpl = customTmpl
	} else {
		embedTemplates, err := templates.ParseEmbedTemplates()
		if err != nil {
			return nil, err
		}
		tmpl = embedTemplates
	}

	return &Sender{client: client, Options: options, template: tmpl}, nil
}

// Sender is responsible for sending email
type Sender struct {
	client   *mail.Client
	Options  Options
	template *template.Template
}

// SendEmail sends an email with given message
func (s *Sender) SendEmail(ctx context.Context, message Message) error {
	if message.From == "" {
		message.From = s.Options.Username
	}
	email, err := s.BuildEmail(message)
	if err != nil {
		return err
	}
	err = s.client.DialAndSendWithContext(ctx, email)
	if err != nil {
		return err
	}
	return nil
}

// Message represents an email message
type Message struct {
	ContentType mail.ContentType `mapstructure:"contentType"`
	From        string           `mapstructure:"from"`
	To          []string         `mapstructure:"to"`
	CC          []string         `mapstructure:"cc"`
	Bcc         []string         `mapstructure:"bcc"`
	Subject     string           `mapstructure:"subject"`
	Message     any              `mapstructure:"message"`
	Template    string           `mapstructure:"template"`
}

// BuildEmail builds *mail.Msg from Message
func (s *Sender) BuildEmail(msg Message) (*mail.Msg, error) {
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

	// parsed email body
	var body string
	if msg.Template != "" {
		tmplBody, err := s.ParseTemplate(msg.Template, msg.Message)
		if err != nil {
			return nil, err
		}
		body = str2bytes.Bytes2Str(tmplBody)
	} else {
		body = msg.Message.(string)
	}
	mailMsg.SetBodyString(msg.ContentType, body)
	return mailMsg, nil
}
