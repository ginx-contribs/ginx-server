package handler

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/ginx-contribs/ginx-server/ent"
	"github.com/ginx-contribs/ginx-server/internal/common/data/cache"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/repo"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx-server/pkg/email"
	"github.com/ginx-contribs/ginx-server/pkg/utils/captcha"
	"github.com/ginx-contribs/ginx-server/pkg/utils/ts"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/jwtx"
	"github.com/ginx-contribs/str2bytes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/wneessen/go-mail"
	"golang.org/x/net/context"
	"time"
)

// AuthHandler is responsible for user authentication
type AuthHandler struct {
	Token          TokenHandler
	UserRepo       repo.UserRepo
	CaptchaHandler CaptchaHandler
}

// EncryptPassword encrypts password with sha512
func (a AuthHandler) EncryptPassword(s string) string {
	sum512 := sha1.Sum(str2bytes.Str2Bytes(s))
	return base64.StdEncoding.EncodeToString(sum512[:])
}

// LoginWithPassword user login by password
func (a AuthHandler) LoginWithPassword(ctx context.Context, option types.LoginOptions) (*types.TokenPair, error) {
	// find user from repository
	queryUser, err := a.UserRepo.FindByNameOrMail(ctx, option.Username)
	if ent.IsNotFound(err) {
		return nil, err
	} else if err != nil { // db error
		return nil, statuserr.InternalError(err)
	}

	// check password
	hashPaswd := a.EncryptPassword(option.Password)
	if queryUser.Password != hashPaswd {
		return nil, types.ErrPasswordMismatch
	}

	// issue token
	tokenPair, err := a.Token.Issue(ctx, types.TokenPayload{
		Username: queryUser.Username,
		UserId:   queryUser.UID,
		Remember: option.Remember,
	}, option.Remember)

	if err != nil {
		return nil, statuserr.InternalError(err)
	}

	return &tokenPair, nil
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

func NewTokenHandler(jwtConf conf.Jwt, client *redis.Client) TokenHandler {
	return TokenHandler{
		Method:       jwt.SigningMethodHS256,
		AccessCache:  cache.NewRedisTokenCache("access", client),
		RefreshCache: cache.NewRedisTokenCache("refresh", client),
		JwtConf:      jwtConf,
	}
}

// TokenHandler is responsible for maintaining authentication tokens
type TokenHandler struct {
	Method       jwt.SigningMethod
	AccessCache  cache.TokenCache
	RefreshCache cache.TokenCache
	JwtConf      conf.Jwt
}

func (t TokenHandler) Issue(ctx context.Context, payload types.TokenPayload, refresh bool) (types.TokenPair, error) {
	now := time.Now()
	var tokenPair types.TokenPair

	// issue access token
	accessToken, err := t.newToken(now, t.JwtConf.Access.Key, payload)
	if err != nil {
		return tokenPair, err
	}

	// consider network latency
	latency := time.Second * 10

	ttl := t.JwtConf.Access.Expire.Duration() + t.JwtConf.Access.Delay.Duration() + latency
	// store into the cache
	if err := t.AccessCache.Set(ctx, accessToken.Claims.ID, accessToken.Claims.ID, ttl); err != nil {
		return types.TokenPair{}, err
	}

	tokenPair.AccessToken = accessToken
	// no need to refresh the token
	if !refresh {
		return tokenPair, nil
	}

	// issue refresh token
	refreshToken, err := t.newToken(now, t.JwtConf.Refresh.Key, payload)
	if err != nil {
		return tokenPair, err
	}

	// associated with access token
	if err := t.RefreshCache.Set(ctx, refreshToken.Claims.ID, accessToken.Claims.ID, t.JwtConf.Refresh.Expire.Duration()); err != nil {
		return tokenPair, nil
	}
	tokenPair.RefreshToken = refreshToken

	return tokenPair, nil
}

// Refresh refreshes the access token lifetime with the given refresh token
func (t TokenHandler) Refresh(ctx context.Context, accessToken string, refreshToken string) (types.TokenPair, error) {
	now := time.Now()
	var pair types.TokenPair
	// return directly if refresh token is expired
	refresh, err := t.VerifyRefresh(ctx, refreshToken)
	if err != nil {
		return pair, err
	}
	pair.RefreshToken = refresh

	// parse access token
	access, err := t.VerifyAccess(ctx, accessToken)
	if errors.Is(err, jwt.ErrTokenExpired) {
		// return if over the delay time
		if access.Claims.ExpiresAt.Add(t.JwtConf.Access.Delay.Duration()).Sub(now) < 0 {
			return pair, jwt.ErrTokenExpired
		}
	} else if err != nil {
		return pair, err
	}

	// check access token if is associated with refresh token
	id, err := t.RefreshCache.Get(ctx, refresh.Claims.ID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return pair, err
	}
	if access.Claims.ID != id {
		return pair, jwt.ErrTokenUnverifiable
	}

	// use a new token to replace the old one
	newAccess, err := t.newToken(now, t.JwtConf.Access.Key, access.Claims.TokenPayload)
	if err != nil {
		return pair, err
	}
	pair.AccessToken = newAccess

	// get rest ttl
	ttl, err := t.AccessCache.TTL(ctx, access.Claims.ID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return pair, statuserr.InternalError(err)
	}
	// extend lifetime of access token
	ttl += t.JwtConf.Access.Expire.Duration()
	if err := t.AccessCache.Set(ctx, newAccess.Claims.ID, newAccess.Claims.ID, ttl); err != nil {
		return pair, statuserr.InternalError(err)
	}

	// update association
	if err := t.RefreshCache.Set(ctx, refresh.Claims.ID, newAccess.Claims.ID, -1); err != nil {
		return pair, statuserr.InternalError(err)
	}

	return pair, nil
}

// VerifyAccess verifies the access token if is valid and parses the payload in the token.
func (t TokenHandler) VerifyAccess(ctx context.Context, token string) (types.Token, error) {
	now := ts.Now()
	parsedToken, err := t.parse(token, t.JwtConf.Access.Key)
	if errors.Is(err, jwt.ErrTokenExpired) {
		// check if token needs to be refreshed
		if parsedToken.Claims.Remember && parsedToken.Claims.ExpiresAt.Add(t.JwtConf.Access.Delay.Duration()).Sub(now) > 0 {
			return parsedToken, types.ErrTokenNeedsRefresh
		}
		return parsedToken, err
	} else if err != nil {
		return parsedToken, err
	}

	// check if exists in cache
	if _, err := t.AccessCache.Get(ctx, parsedToken.Claims.ID); errors.Is(err, redis.Nil) {
		return parsedToken, jwt.ErrTokenExpired
	} else if err != nil {
		return parsedToken, statuserr.InternalError(err)
	}
	return parsedToken, nil
}

// VerifyRefresh verifies the refresh token if is valid.
func (t TokenHandler) VerifyRefresh(ctx context.Context, token string) (types.Token, error) {
	parsedToken, err := t.parse(token, t.JwtConf.Refresh.Key)
	if err != nil {
		return parsedToken, err
	}
	// check if exists in cache
	if _, err := t.RefreshCache.Get(ctx, parsedToken.Claims.ID); errors.Is(err, redis.Nil) {
		return parsedToken, jwt.ErrTokenExpired
	} else if err != nil {
		return parsedToken, statuserr.InternalError(err)
	}
	return parsedToken, nil
}

func (t TokenHandler) newToken(now time.Time, key string, payload types.TokenPayload) (types.Token, error) {
	// create the token claims
	claims := types.TokenClaims{
		TokenPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.JwtConf.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.JwtConf.Access.Expire.Duration())),
			ID:        uuid.NewString(),
		},
	}

	// issue the token
	token, err := jwtx.IssueWithClaims(str2bytes.Str2Bytes(key), t.Method, claims)
	if err != nil {
		return types.Token{}, err
	}

	return types.Token{
		Token:       token.Token,
		Claims:      claims,
		TokenString: token.SignedString,
	}, err
}

func (t TokenHandler) parse(token, secret string) (types.Token, error) {
	parseJwt, err := jwtx.VerifyWithClaims(token, str2bytes.Str2Bytes(secret), t.Method, &types.TokenClaims{})
	if err == nil || errors.Is(err, jwt.ErrTokenExpired) {
		return types.Token{
			Token:       parseJwt.Token,
			Claims:      *parseJwt.Claims.(*types.TokenClaims),
			TokenString: parseJwt.SignedString,
		}, nil
	} else {
		return types.Token{}, err
	}
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
