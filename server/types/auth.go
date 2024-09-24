package types

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/golang-jwt/jwt/v5"
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

	ErrRateLimitExceeded = statuserr.Errorf("server busy").SetCode(1_429_001).SetStatus(status.TooManyRequests)
)

type AuthLoginOption struct {
	// username or email
	Username string `json:"username" binding:"required"`
	// user password
	Password string `json:"password" binding:"required"`
	// remember user or not
	Remember bool `json:"remember"`
}

type AuthRegisterOption struct {
	// username must be alphanumeric
	Username string `json:"username" binding:"required,alphanum"`
	// user password
	Password string `json:"password" binding:"required"`
	// user email address
	Email string `json:"email" binding:"email"`
	// verification code from verify email
	Code string `json:"code" binding:"required,alphanum"`
}

type AuthResetPasswordOption struct {
	// user email address
	Email string `json:"email" binding:"email"`
	// new password
	Password string `json:"password" binding:"required"`
	// verification code from verify email
	Code string `json:"code" binding:"required"`
}

type AuthRefreshTokenOption struct {
	// access token
	AccessToken string `json:"accessToken" binding:"required"`
	// refresh token
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type AuthVerifyCodeOption struct {
	// email receiver
	To string `json:"to" binding:"email"`
	// verify code usage: 1-register 2-reset password
	Usage Usage `json:"usage" binding:"required,gte=1,lte=2"`
}

type TokenPayload struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
	Remember bool   `json:"remember"`
}

// TokenClaims is payload info in jwt
type TokenClaims struct {
	TokenPayload
	jwt.RegisteredClaims
}

// Token represents a jwt token
type Token struct {
	Token       *jwt.Token
	Claims      TokenClaims
	TokenString string
}

// TokenPair represents a jwt token pair composed of access token and refresh token
type TokenPair struct {
	AccessToken  Token
	RefreshToken Token
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

const tokenKey = "auth.token.context.info.key"

// SetTokenInfo stores token information into context
func SetTokenInfo(ctx *gin.Context, token *Token) {
	ctx.Set(tokenKey, token)
}

// GetTokenInfo returns token information from context
func GetTokenInfo(ctx *gin.Context) (*Token, error) {
	value, exists := ctx.Get(tokenKey)
	if !exists {
		return nil, errors.New("there is no token in context")
	}

	if token, ok := value.(*Token); !ok {
		return nil, fmt.Errorf("expected %T, got %T", &Token{}, value)
	} else if token == nil {
		return nil, errors.New("nil token in context")
	} else {
		return token, nil
	}
}

// MustGetTokenInfo returns token information from context, panic if err != nil
func MustGetTokenInfo(ctx *gin.Context) *Token {
	info, err := GetTokenInfo(ctx)
	if err != nil {
		panic(err)
	}
	return info
}
