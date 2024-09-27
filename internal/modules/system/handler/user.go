package handler

import (
	"context"
	"github.com/ginx-contribs/ginx-server/ent"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/repo"
	types2 "github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
)

type UserHandler struct {
	UserRepo *repo.UserRepo
}

func (u UserHandler) FindByUID(ctx context.Context, uid string) (types2.UserInfo, error) {
	record, err := u.UserRepo.FindByUID(ctx, uid)
	if ent.IsNotFound(err) {
		return types2.UserInfo{}, types2.ErrUserNotFund
	} else if err != nil {
		return types2.UserInfo{}, statuserr.InternalError(err)
	}
	return types2.EntToUser(record), nil
}

func (u UserHandler) ListUserByPage(ctx context.Context, page, size int, pattern string) (types2.UserSearchResult, error) {
	users, err := u.UserRepo.ListByPage(ctx, page, size, pattern)
	if err != nil {
		return types2.UserSearchResult{}, statuserr.InternalError(err)
	}
	toUsers := types2.EntsToUsers(users)
	return types2.UserSearchResult{Total: int64(len(users)), List: toUsers}, nil
}

func (u UserHandler) CreateUser(ctx context.Context, username string, email string, password string) (types2.UserInfo, error) {
	newUser, err := u.UserRepo.CreateNewUser(ctx, username, email, password)
	if err != nil {
		return types2.UserInfo{}, statuserr.InternalError(err)
	}
	return types2.EntToUser(newUser), nil
}

func (u UserHandler) RemoveUser(ctx context.Context, uid string) error {
	deleted, err := u.UserRepo.RemoveByUID(ctx, uid)
	if err != nil {
		return statuserr.InternalError(err)
	} else if deleted == 0 {
		return statuserr.Errorf("expected deleted 1, but got 0, remove user failed: uid %s", uid).SetStatus(status.InternalServerError)
	}
	return nil
}
