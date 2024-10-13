package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/handler"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx-server/pkg/token"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"github.com/golang-jwt/jwt/v5"
)

type AuthAPI struct {
	TokenResolver  *token.Resolver
	AuthHandler    handler.AuthHandler
	CaptchaHandler handler.CaptchaHandler
}

// Login
// @Summary      Login
// @Description  login with password, and returns jwt token pair
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        LoginOptions  body  types.LoginOptions  true "LoginOptions"
// @Success      200  {object}  types.Response{data=types.TokenResult}
// @Router       /auth/login [POST]
func (a *AuthAPI) Login(ctx *gin.Context) {
	var loginOpt types.LoginOptions
	if err := ginx.ShouldValidateJSON(ctx, &loginOpt); err != nil {
		return
	}

	// login by username and password
	tokenPair, err := a.AuthHandler.LoginWithPassword(ctx, loginOpt)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
		return
	}

	resp.Ok(ctx).Msg("login ok").Data(types.TokenResult{
		AccessToken:  tokenPair.Access.Raw,
		RefreshToken: tokenPair.Refresh.Raw,
	}).JSON()
}

// Register
// @Summary      Register
// @Description  register a new user with verification code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        RegisterOptions  body  types.RegisterOptions  true "register params"
// @Success      200  {object}  types.Response
// @Router       /auth/register [POST]
func (a *AuthAPI) Register(ctx *gin.Context) {
	var registerOpt types.RegisterOptions
	if err := ginx.ShouldValidateJSON(ctx, &registerOpt); err != nil {
		return
	}

	_, err := a.AuthHandler.RegisterNewUser(ctx, registerOpt)
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
// @Param        ResetOptions   body  types.ResetOptions  true "ResetOptions"
// @Success      200  {object}  types.Response
// @Router       /auth/reset [POST]
func (a *AuthAPI) ResetPassword(ctx *gin.Context) {
	var restOpt types.ResetOptions
	if err := ginx.ShouldValidateJSON(ctx, &restOpt); err != nil {
		return
	}

	if err := a.AuthHandler.ResetPassword(ctx, restOpt); err != nil {
		resp.Fail(ctx).Error(err).JSON()
		return
	}
	resp.Ok(ctx).Msg("reset password ok").JSON()
}

// Refresh
// @Summary      Refresh
// @Description  ask for refresh access TokenHandler lifetime with refresh TokenHandler
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        RefreshTokenOptions  body  types.RefreshTokenOptions  true "refresh params"
// @Success      200  {object}  types.Response{data=types.TokenResult}
// @Router       /auth/refresh [POST]
func (a *AuthAPI) Refresh(ctx *gin.Context) {
	var refreshOpt types.RefreshTokenOptions
	if err := ginx.ShouldValidateJSON(ctx, &refreshOpt); err != nil {
		return
	}

	// ask for refresh TokenHandler
	tokenPair, err := a.TokenResolver.Refresh(ctx, refreshOpt.AccessToken, refreshOpt.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			resp.Fail(ctx).Error(types.ErrCredentialExpired).JSON()
		} else {
			ctx.Error(err)
			resp.Fail(ctx).Error(types.ErrCredentialInvalid).JSON()
		}
		return
	}
	resp.Ok(ctx).Msg("refresh ok").Data(types.TokenResult{
		AccessToken:  tokenPair.Access.Raw,
		RefreshToken: tokenPair.Refresh.Raw,
	}).JSON()
}

// Captcha
// @Summary      Captcha
// @Description  send captcha code mail to specified email address
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        CaptchaOption   body   types.CaptchaOption  true  "CaptchaOption"
// @Success      200  {object}  types.Response
// @Router       /auth/captcha [POST]
func (a *AuthAPI) Captcha(ctx *gin.Context) {
	var verifyOpt types.CaptchaOption
	if err := ginx.ShouldValidateJSON(ctx, &verifyOpt); err != nil {
		return
	}

	// check usage
	if err := types.CheckValidUsage(verifyOpt.Usage); err != nil {
		resp.Fail(ctx).Error(types.ErrVerifyCodeUsageUnsupported).JSON()
		return
	}

	err := a.CaptchaHandler.SendCaptchaEmail(ctx, verifyOpt.To, verifyOpt.Usage)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
		return
	}
	resp.Ok(ctx).Msg("mail has been sent").JSON()
}
