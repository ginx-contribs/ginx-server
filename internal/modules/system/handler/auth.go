package handler

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/ginx-contribs/ginx-server/ent"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/cache"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/repo"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx-server/pkg/email"
	"github.com/ginx-contribs/ginx-server/pkg/token"
	"github.com/ginx-contribs/ginx-server/pkg/utils/captcha"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/str2bytes"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/wneessen/go-mail"
	"golang.org/x/net/context"
)

// AuthHandler is responsible for user authentication
type AuthHandler struct {
	Token          *token.Resolver
	UserRepo       repo.UserRepo
	CaptchaHandler CaptchaHandler
}

// EncryptPassword encrypts password with sha512
func (a AuthHandler) EncryptPassword(s string) string {
	sum512 := sha1.Sum(str2bytes.Str2Bytes(s))
	return base64.StdEncoding.EncodeToString(sum512[:])
}

// LoginWithPassword user login by password
func (a AuthHandler) LoginWithPassword(ctx context.Context, option types.LoginOptions) (token.Pair, error) {
	// find user from repository
	queryUser, err := a.UserRepo.FindByNameOrMail(ctx, option.Username)
	if ent.IsNotFound(err) {
		return token.Pair{}, err
	} else if err != nil { // db error
		return token.Pair{}, statuserr.InternalError(err)
	}

	// check password
	hashPaswd := a.EncryptPassword(option.Password)
	if queryUser.Password != hashPaswd {
		return token.Pair{}, types.ErrPasswordMismatch
	}

	// issue token
	tokenPair, err := a.Token.Issue(ctx, types.TokenPayload{
		Username: queryUser.Username,
		UserId:   queryUser.UID,
	}, option.Remember)

	if err != nil {
		return token.Pair{}, statuserr.InternalError(err)
	}

	return tokenPair, nil
}

// RegisterNewUser registers new user and returns it
func (a AuthHandler) RegisterNewUser(ctx context.Context, option types.RegisterOptions) (*ent.User, error) {

	// check verify code if is valid
	err := a.CaptchaHandler.Check(ctx, option.Email, option.Code, types.UsageRegister)
	if err != nil {
		return nil, err
	}

	// check username if is duplicate
	userByName, err := a.UserRepo.FindByName(ctx, option.Username)
	if !ent.IsNotFound(err) && err != nil {
		return nil, statuserr.InternalError(err)
	} else if userByName != nil {
		return nil, types.ErrUserAlreadyExists
	}

	// check email if is duplicate
	userByEmail, err := a.UserRepo.FindByEmail(ctx, option.Email)
	if !ent.IsNotFound(err) && err != nil {
		return nil, statuserr.InternalError(err)
	} else if userByEmail != nil {
		return nil, types.ErrEmailAlreadyUsed
	}

	// create new user
	user, err := a.UserRepo.CreateNewUser(ctx, option.Username, option.Email, a.EncryptPassword(option.Password))
	if err != nil {
		return nil, statuserr.InternalError(err)
	}

	// remove verify code
	err = a.CaptchaHandler.Remove(ctx, option.Code, types.UsageRegister)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ResetPassword resets specified user password and returns uid
func (a AuthHandler) ResetPassword(ctx context.Context, option types.ResetOptions) error {

	// check verify code if is valid
	err := a.CaptchaHandler.Check(ctx, option.Email, option.Code, types.UsageReset)
	if err != nil {
		return err
	}

	// check email if is already registered
	queryUser, err := a.UserRepo.FindByEmail(ctx, option.Email)
	if ent.IsNotFound(err) {
		return types.ErrUserNotFund
	}

	// update password
	_, err = a.UserRepo.UpdateOnePassword(ctx, queryUser.ID, a.EncryptPassword(option.Password))
	if err != nil {
		return statuserr.InternalError(err)
	}

	// remove verify code
	err = a.CaptchaHandler.Remove(ctx, option.Code, types.UsageReset)
	if err != nil {
		return err
	}
	return nil
}

type CaptchaHandler struct {
	CaptchaCache cache.CaptchaCache
	EmailHandler EmailHandler

	MetaInfo conf.MetaInfo
}

// SendCaptchaEmail send a verify code email to the specified address
func (v CaptchaHandler) SendCaptchaEmail(ctx context.Context, to string, usage types.Usage) error {
	ttl := v.EmailHandler.Config.Code.TTL
	retryttl := v.EmailHandler.Config.Code.RetryTTL
	var code string

	for try := 0; ; try++ {
		// max retry 10 times
		if try > 10 {
			return types.ErrVerifyCodeRetryLater
		}

		// generated a verification code
		code = captcha.GenCaptcha(8)
		tryOk, err := v.CaptchaCache.Set(ctx, usage, code, to, ttl.Duration(), retryttl.Duration())

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

	msg := email.Message{
		ContentType: mail.TypeTextHTML,
		To:          []string{to},
		Subject:     fmt.Sprintf("you are applying for verification code tp %s.", usage.String()),
		Message: map[string]any{
			"to":       to,
			"action":   usage.String(),
			"duration": ttl.String(),
			"code":     code,
			"author":   v.MetaInfo.Author,
		},
		Template: email.TemplateCaptcha,
	}

	// send email
	err := v.EmailHandler.Publish(ctx, msg)
	if err != nil {
		if err := v.CaptchaCache.Del(ctx, usage, code); err != nil {
			return statuserr.InternalError(err)
		}
		return err
	}
	return nil
}

// Check checks captcha if is valid
func (v CaptchaHandler) Check(ctx context.Context, to, code string, usage types.Usage) error {
	getTo, err := v.CaptchaCache.Get(ctx, usage, code)
	if errors.Is(err, redis.Nil) || getTo != to {
		return types.ErrVerifyCodeInvalid
	} else if err != nil {
		return statuserr.InternalError(err)
	}
	return nil
}

// Remove removes captcha from cache
func (v CaptchaHandler) Remove(ctx context.Context, code string, usage types.Usage) error {
	err := v.CaptchaCache.Del(ctx, usage, code)
	if err != nil {
		return statuserr.InternalError(err)
	}
	return nil
}
