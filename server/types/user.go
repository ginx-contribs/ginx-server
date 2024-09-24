package types

import "github.com/ginx-contribs/ginx-server/server/data/ent"

type UserSearchOption struct {
	Page   int    `form:"page" binding:"required,gt=0"`
	Size   int    `form:"size" binding:"required,gt=0"`
	Search string `form:"search"`
}

type ULIDOptions struct {
	Uid string `form:"uid" binding:"required"`
}

type UserInfo struct {
	Uid       string `json:"uid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}

type UserSearchResult struct {
	Total int64      `json:"total"`
	List  []UserInfo `json:"list"`
}

func EntToUser(user *ent.User) UserInfo {
	if user == nil {
		return UserInfo{}
	}

	return UserInfo{
		Uid:       user.UID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func EntsToUsers(users []*ent.User) []UserInfo {
	if users == nil {
		return []UserInfo{}
	}
	var us []UserInfo
	for _, u := range users {
		us = append(us, EntToUser(u))
	}
	return us
}
