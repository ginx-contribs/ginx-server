package test

import (
	"context"
	"github.com/ginx-contribs/ginx-server/server/types"
	"github.com/ginx-contribs/ginx-server/test/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserHandler_Create_Query_Delete(t *testing.T) {
	ctx := context.Background()
	testServer, err := testutil.NewTestServer(ctx)
	if !assert.NoError(t, err) {
		return
	}
	defer testServer.Cleanup()
	userHandler := testServer.Sc.UserHandler

	samples := []struct {
		Name, Email, Paswd string
	}{
		{"jack", "a@qq.com", "123456"},
		{"mike", "mike@gmail.com", "789123"},
		{"1=1", "test@gmail.com", "1=1"},
	}

	for _, sample := range samples {
		wantUser, err := userHandler.CreateUser(ctx, sample.Name, sample.Email, sample.Paswd)
		if !assert.NoError(t, err) {
			return
		}
		t.Log("created", wantUser)
		queryUser, err := userHandler.FindByUID(ctx, wantUser.Uid)
		// must not exist
		if !assert.NoError(t, err) {
			return
		}
		// must be equal
		if !assert.EqualValues(t, wantUser, queryUser) {
			return
		}
		t.Log("found", queryUser)
		// remove user
		err = userHandler.RemoveUser(ctx, queryUser.Uid)
		if !assert.NoError(t, err) {
			return
		}
		// query again
		queryUser2, err := userHandler.FindByUID(ctx, wantUser.Uid)
		// must be not found
		if !assert.ErrorIs(t, err, types.ErrUserNotFund) {
			return
		}
		// must be not equal
		if !assert.NotEqualValues(t, wantUser, queryUser2) {
			return
		}
		t.Log("deleted", queryUser.Uid)
	}
}
