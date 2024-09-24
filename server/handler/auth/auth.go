package auth

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/ginx-contribs/ginx-server/server/data/ent"
	"github.com/ginx-contribs/ginx-server/server/data/repo"
	"github.com/ginx-contribs/ginx-server/server/types"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/str2bytes"
	"golang.org/x/net/context"
)

func NewAuthHandler(userRepo *repo.UserRepo, tokenHandler *TokenHandler, verify *VerifyCodeHandler) *AuthHandler {
	return &AuthHandler{token: tokenHandler, userRepo: userRepo, verifyCode: verify}
}

// AuthHandler is responsible for user authentication
type AuthHandler struct {
	token      *TokenHandler
	userRepo   *repo.UserRepo
	verifyCode *VerifyCodeHandler
}

// EncryptPassword encrypts password with sha512
func (a *AuthHandler) EncryptPassword(s string) string {
	sum512 := sha1.Sum(str2bytes.Str2Bytes(s))
	return base64.StdEncoding.EncodeToString(sum512[:])
}

// LoginWithPassword user login by password
func (a *AuthHandler) LoginWithPassword(ctx context.Context, option types.AuthLoginOption) (*types.TokenPair, error) {
	// find user from repository
	queryUser, err := a.userRepo.FindByNameOrMail(ctx, option.Username)
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
	tokenPair, err := a.token.Issue(ctx, types.TokenPayload{
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
func (a *AuthHandler) RegisterNewUser(ctx context.Context, option types.AuthRegisterOption) (*ent.User, error) {

	// check verify code if is valid
	err := a.verifyCode.CheckVerifyCode(ctx, option.Email, option.Code, types.UsageRegister)
	if err != nil {
		return nil, err
	}

	// check username if is duplicate
	userByName, err := a.userRepo.FindByName(ctx, option.Username)
	if !ent.IsNotFound(err) && err != nil {
		return nil, statuserr.InternalError(err)
	} else if userByName != nil {
		return nil, types.ErrUserAlreadyExists
	}

	// check email if is duplicate
	userByEmail, err := a.userRepo.FindByEmail(ctx, option.Email)
	if !ent.IsNotFound(err) && err != nil {
		return nil, statuserr.InternalError(err)
	} else if userByEmail != nil {
		return nil, types.ErrEmailAlreadyUsed
	}

	// create new user
	user, err := a.userRepo.CreateNewUser(ctx, option.Username, option.Email, a.EncryptPassword(option.Password))
	if err != nil {
		return nil, statuserr.InternalError(err)
	}

	// remove verify code
	err = a.verifyCode.RemoveVerifyCode(ctx, option.Code, types.UsageRegister)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ResetPassword resets specified user password and returns uid
func (a *AuthHandler) ResetPassword(ctx context.Context, option types.AuthResetPasswordOption) error {

	// check verify code if is valid
	err := a.verifyCode.CheckVerifyCode(ctx, option.Email, option.Code, types.UsageReset)
	if err != nil {
		return err
	}

	// check email if is already registered
	queryUser, err := a.userRepo.FindByEmail(ctx, option.Email)
	if ent.IsNotFound(err) {
		return types.ErrUserNotFund
	}

	// update password
	_, err = a.userRepo.UpdateOnePassword(ctx, queryUser.ID, a.EncryptPassword(option.Password))
	if err != nil {
		return statuserr.InternalError(err)
	}

	// remove verify code
	err = a.verifyCode.RemoveVerifyCode(ctx, option.Code, types.UsageReset)
	if err != nil {
		return err
	}
	return nil
}
