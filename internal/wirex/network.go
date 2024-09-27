package wirex

import (
	"context"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/ginx-server/pkg/email"
)

func NewEmail(ctx context.Context, emailConf conf.Email) (*email.Sender, error) {
	return email.NewSender(email.Options{
		Host:        emailConf.Host,
		Port:        emailConf.Port,
		Username:    emailConf.Username,
		Password:    emailConf.Password,
		TemplateDir: emailConf.Template,
	})
}
