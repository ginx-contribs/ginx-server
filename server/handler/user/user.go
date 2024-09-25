package user

import (
	"context"
	"github.com/ginx-contribs/ginx-server/server/data/ent"
	"github.com/ginx-contribs/ginx-server/server/data/repo"
	"github.com/ginx-contribs/ginx-server/server/types"
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
)

func NewUserHandler(userRepo *repo.UserRepo) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

type UserHandler struct {
	userRepo *repo.UserRepo
}

func (u *UserHandler) FindByUID(ctx context.Context, uid string) (types.UserInfo, error) {
	record, err := u.userRepo.FindByUID(ctx, uid)
	if ent.IsNotFound(err) {
		return types.UserInfo{}, types.ErrUserNotFund
	} else if err != nil {
		return types.UserInfo{}, statuserr.InternalError(err)
	}
	return types.EntToUser(record), nil
}

func (u *UserHandler) ListUserByPage(ctx context.Context, page, size int, pattern string) (types.UserSearchResult, error) {
	users, err := u.userRepo.ListByPage(ctx, page, size, pattern)
	if err != nil {
		return types.UserSearchResult{}, statuserr.InternalError(err)
	}
	toUsers := types.EntsToUsers(users)
	return types.UserSearchResult{Total: int64(len(users)), List: toUsers}, nil
}

func (u *UserHandler) CreateUser(ctx context.Context, username string, email string, password string) (types.UserInfo, error) {
	newUser, err := u.userRepo.CreateNewUser(ctx, username, email, password)
	if err != nil {
		return types.UserInfo{}, statuserr.InternalError(err)
	}
	return types.EntToUser(newUser), nil
}

func (u *UserHandler) RemoveUser(ctx context.Context, uid string) error {
	deleted, err := u.userRepo.RemoveByUID(ctx, uid)
	if err != nil {
		return statuserr.InternalError(err)
	} else if deleted == 0 {
		return statuserr.Errorf("expected deleted 1, but got 0, remove user failed: uid %s", uid).SetStatus(status.InternalServerError)
	}
	return nil
}
