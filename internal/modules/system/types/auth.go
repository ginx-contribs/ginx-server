package types

import (
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
)

var (
	ErrUserNotFund       = statuserr.Errorf("user not found").SetCode(1_400_000).SetStatus(status.BadRequest)
	ErrUserAlreadyExists = statuserr.Errorf("user already exists").SetCode(1_400_002).SetStatus(status.BadRequest)
	ErrPasswordMismatch  = statuserr.Errorf("password mismatch").SetCode(1_400_004).SetStatus(status.BadRequest)
	ErrEmailAlreadyUsed  = statuserr.Errorf("email already used by other").SetCode(1_400_016).SetStatus(status.BadRequest)

	ErrVerifyCodeRetryLater       = statuserr.Errorf("retry applying for verify code later").SetCode(1_400_032).SetStatus(status.BadRequest)
	ErrVerifyCodeInvalid          = statuserr.Errorf("invliad verify code").SetCode(1_400_033).SetStatus(status.BadRequest)
	ErrVerifyCodeUsageUnsupported = statuserr.Errorf("verify code usage unsupported").SetCode(1_400_036).SetStatus(status.BadRequest)

	ErrCredentialInvalid = statuserr.Errorf("invalid credential").SetCode(1_401_001).SetStatus(status.Unauthorized)
	ErrCredentialExpired = statuserr.Errorf("credential expired").SetCode(1_401_002).SetStatus(status.Unauthorized)
	ErrTokenNeedsRefresh = statuserr.Errorf("token need to refresh").SetCode(1_401_003).SetStatus(status.Unauthorized)
)

type LoginOptions struct {
	// username or email
	Username string `json:"username" binding:"required"`
	// user password
	Password string `json:"password" binding:"required"`
	// remember user or not
	Remember bool `json:"remember"`
}

type RegisterOptions struct {
	// username must be alphanumeric
	Username string `json:"username" binding:"required,alphanum"`
	// user password
	Password string `json:"password" binding:"required"`
	// user email address
	Email string `json:"email" binding:"email"`
	// verification code from verify email
	Code string `json:"code" binding:"required,alphanum"`
}

type ResetOptions struct {
	// user email address
	Email string `json:"email" binding:"email"`
	// new password
	Password string `json:"password" binding:"required"`
	// verification code from verify email
	Code string `json:"code" binding:"required"`
}

type RefreshTokenOptions struct {
	// access token
	AccessToken string `json:"accessToken" binding:"required"`
	// refresh token
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type CaptchaOption struct {
	// email receiver
	To string `json:"to" binding:"email"`
	// verify code usage: 1-register 2-reset password
	Usage Usage `json:"usage" binding:"required,gte=1,lte=2"`
}

type TokenPayload struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

type TokenResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

const (
	UsageUnknown  Usage = 0
	UsageRegister Usage = 1
	UsageReset    Usage = 2
)

type Usage int

func (u Usage) Name() string {
	switch u {
	case 1:
		return "register"
	case 2:
		return "reset"
	default:
		return "unknown"
	}
}

func (u Usage) String() string {
	switch u {
	case 1:
		return "register account"
	case 2:
		return "reset password"
	default:
		return "unknown usage"
	}
}

func CheckValidUsage(u Usage) error {
	if u.String() == UsageUnknown.String() {
		return ErrVerifyCodeUsageUnsupported
	}
	return nil
}
