package wirex

import (
	"context"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/ginx-server/pkg/email"
	"github.com/ginx-contribs/ginx-server/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func NewEmailSender(ctx context.Context, emailConf conf.Email) (*email.Sender, error) {
	return email.NewSender(email.Options{
		Host:        emailConf.Host,
		Port:        emailConf.Port,
		Username:    emailConf.Username,
		Password:    emailConf.Password,
		TemplateDir: emailConf.Template,
	})
}

func NewTokenResolver(ctx context.Context, jwtconf conf.Jwt, client *redis.Client) (*token.Resolver, error) {
	return token.NewResolver(token.Options{
		Cache:          token.NewRedisTokenCache(client),
		Issuer:         jwtconf.Issuer,
		AccessSecret:   jwtconf.Access.Key,
		AccessMethod:   jwt.SigningMethodHS512,
		AccessExpired:  jwtconf.Access.Expire.Duration(),
		AccessDelay:    jwtconf.Access.Delay.Duration(),
		RefreshSecret:  jwtconf.Refresh.Key,
		RefreshMethod:  jwt.SigningMethodHS512,
		RefreshExpired: jwtconf.Refresh.Expire.Duration(),
	}), nil
}
