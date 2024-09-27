package wirex

import (
	"context"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/wneessen/go-mail"
)

// NewEmailClient initialize email client
func NewEmailClient(ctx context.Context, emailConf conf.Email) (*mail.Client, error) {
	client, err := mail.NewClient(emailConf.Host,
		mail.WithPort(emailConf.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(emailConf.Username),
		mail.WithPassword(emailConf.Password),
	)
	if err != nil {
		return nil, err
	}
	// test if smtp server is available
	err = client.DialWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
