package auth

import (
	"errors"
	"fmt"
	"github.com/ginx-contribs/ginx-server/server/data/cache"
	"github.com/ginx-contribs/ginx-server/server/handler/email"
	"github.com/ginx-contribs/ginx-server/server/types"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"math/rand/v2"
)

func NewVerifyCodeHandler(codeCache cache.VerifyCodeCache, sender *email.Handler) *VerifyCodeHandler {
	return &VerifyCodeHandler{
		codeCache: codeCache,
		sender:    sender,
	}
}

type VerifyCodeHandler struct {
	codeCache cache.VerifyCodeCache
	sender    *email.Handler
}

// SendVerifyCodeEmail send a verify code email to the specified address
func (v *VerifyCodeHandler) SendVerifyCodeEmail(ctx context.Context, to string, usage types.Usage) error {
	ttl := v.sender.Cfg.Code.TTL
	retryttl := v.sender.Cfg.Code.RetryTTL
	var code string

	for try := 0; ; try++ {
		// max retry 10 times
		if try > 10 {
			return types.ErrVerifyCodeRetryLater
		}

		// generated a verification code
		code = NewVerifyCode(8)
		tryOk, err := v.codeCache.Set(ctx, usage, code, to, ttl.Duration(), retryttl.Duration())

		if errors.Is(err, cache.ErrCodeRepeated) {
			continue
		} else if err != nil {
			return statuserr.InternalError(err)
		} else if !tryOk {
			return types.ErrVerifyCodeRetryLater
		} else {
			break
		}
	}

	pendingMail := email.TmplConfirmCode(usage.String(), to, code, ttl.Duration())

	// send email
	err := v.sender.PublishHermesEmail(ctx, fmt.Sprintf("you are applying for verification code for %s.", usage.String()), []string{to}, pendingMail)
	if err != nil {
		if err := v.codeCache.Del(ctx, usage, code); err != nil {
			return statuserr.InternalError(err)
		}
		return err
	}

	return nil
}

// CheckVerifyCode check verify code if is valid
func (v *VerifyCodeHandler) CheckVerifyCode(ctx context.Context, to, code string, usage types.Usage) error {
	getTo, err := v.codeCache.Get(ctx, usage, code)
	if errors.Is(err, redis.Nil) || getTo != to {
		return types.ErrVerifyCodeInvalid
	} else if err != nil {
		return statuserr.InternalError(err)
	}
	return nil
}

func (v *VerifyCodeHandler) RemoveVerifyCode(ctx context.Context, code string, usage types.Usage) error {
	err := v.codeCache.Del(ctx, usage, code)
	if err != nil {
		return statuserr.InternalError(err)
	}
	return nil
}

// NewVerifyCode returns a verification code with specified length.
// the bigger n is, the conflicts will be less, the recommended n is 8.
func NewVerifyCode(n int) string {
	code := make([]byte, n)
	for i, _ := range code {
		if rand.Int()%2 == 1 {
			code[i] = '0' + byte(rand.Int()%10)
		} else {
			code[i] = 'A' + byte(rand.Int()%26)
		}
	}
	return string(code)
}
