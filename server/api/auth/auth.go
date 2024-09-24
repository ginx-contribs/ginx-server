package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	authandler "github.com/ginx-contribs/ginx-server/server/handler/auth"
	"github.com/ginx-contribs/ginx-server/server/types"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthAPI(token *authandler.TokenHandler, auth *authandler.AuthHandler, verifycode *authandler.VerifyCodeHandler) *AuthAPI {
	return &AuthAPI{token: token, auth: auth, verifycode: verifycode}
}

type AuthAPI struct {
	token      *authandler.TokenHandler
	auth       *authandler.AuthHandler
	verifycode *authandler.VerifyCodeHandler
}

// Login
// @Summary      Login
// @Description  login with password, and returns jwt token pair
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginOption  body  types.AuthLoginOption  true "login params"
// @Success      200  {object}  types.Response{data=types.TokenResult}
// @Router       /auth/login [POST]
func (a *AuthAPI) Login(ctx *gin.Context) {
	var loginOpt types.AuthLoginOption
	if err := ginx.ShouldValidateJSON(ctx, &loginOpt); err != nil {
		return
	}

	// login by username and password
	tokenPair, err := a.auth.LoginWithPassword(ctx, loginOpt)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Msg("login ok").Data(types.TokenResult{
			AccessToken:  tokenPair.AccessToken.TokenString,
			RefreshToken: tokenPair.RefreshToken.TokenString,
		}).JSON()
	}
}

// Register
// @Summary      Register
// @Description  register a new user with verification code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        AuthRegisterOption  body  types.AuthRegisterOption  true "register params"
// @Success      200  {object}  types.Response
// @Router       /auth/register [POST]
func (a *AuthAPI) Register(ctx *gin.Context) {
	var registerOpt types.AuthRegisterOption
	if err := ginx.ShouldValidateJSON(ctx, &registerOpt); err != nil {
		return
	}

	_, err := a.auth.RegisterNewUser(ctx, registerOpt)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Msg("register ok").JSON()
	}
}

// ResetPassword
// @Summary      ResetPassword
// @Description  reset user password with verification code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        AuthResetPasswordOption   body  types.AuthResetPasswordOption  true "reset params"
// @Success      200  {object}  types.Response
// @Router       /auth/reset [POST]
func (a *AuthAPI) ResetPassword(ctx *gin.Context) {
	var restOpt types.AuthResetPasswordOption
	if err := ginx.ShouldValidateJSON(ctx, &restOpt); err != nil {
		return
	}

	if err := a.auth.ResetPassword(ctx, restOpt); err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Msg("reset ok").JSON()
	}
}

// Refresh
// @Summary      Refresh
// @Description  ask for refresh access token lifetime with refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        AuthRefreshTokenOption  body  types.AuthRefreshTokenOption  true "refresh params"
// @Success      200  {object}  types.Response{data=types.TokenResult}
// @Router       /auth/refresh [POST]
func (a *AuthAPI) Refresh(ctx *gin.Context) {
	var refreshOpt types.AuthRefreshTokenOption
	if err := ginx.ShouldValidateJSON(ctx, &refreshOpt); err != nil {
		return
	}

	// ask for refresh token
	tokenPair, err := a.token.Refresh(ctx, refreshOpt.AccessToken, refreshOpt.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			resp.Fail(ctx).Error(types.ErrCredentialExpired).JSON()
		} else {
			ctx.Error(err)
			resp.Fail(ctx).Error(types.ErrCredentialInvalid).JSON()
		}
	} else {
		resp.Ok(ctx).Msg("refresh ok").Data(types.TokenResult{
			AccessToken:  tokenPair.AccessToken.TokenString,
			RefreshToken: tokenPair.RefreshToken.TokenString,
		}).JSON()
	}
}

// VerifyCode
// @Summary      VerifyCode
// @Description  send verification code mail to specified email address
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        AuthVerifyCodeOption   body   types.AuthVerifyCodeOption  true  "AuthVerifyCodeOption"
// @Success      200  {object}  types.Response
// @Router       /auth/code [POST]
func (a *AuthAPI) VerifyCode(ctx *gin.Context) {
	var verifyOpt types.AuthVerifyCodeOption
	if err := ginx.ShouldValidateJSON(ctx, &verifyOpt); err != nil {
		return
	}

	// check usage
	if err := types.CheckValidUsage(verifyOpt.Usage); err != nil {
		resp.Fail(ctx).Error(types.ErrVerifyCodeUsageUnsupported).JSON()
		return
	}

	err := a.verifycode.SendVerifyCodeEmail(ctx, verifyOpt.To, verifyOpt.Usage)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Msg("mail has been sent").JSON()
	}
}
