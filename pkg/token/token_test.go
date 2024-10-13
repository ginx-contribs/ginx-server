package token

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenIssue_NoRefresh_OK(t *testing.T) {
	ctx := context.Background()
	resolver := NewResolver(Options{})
	pair, err := resolver.Issue(ctx, map[string]any{"a": "b"}, false)
	if !assert.NoError(t, err) {
		return
	}
	t.Logf("%+v", pair)
}

func TestTokenIssue_Payload(t *testing.T) {
	ctx := context.Background()
	resolver := NewResolver(Options{})
	pair, err := resolver.Issue(ctx, map[string]any{"a": "b"}, false)
	if !assert.NoError(t, err) {
		return
	}
	m := pair.Access.Claims.Payload
	if !assert.EqualValues(t, m["a"], "b") {
		return
	}
	t.Logf("%+v", m)
}

func TestTokenIssue_WithRefresh_OK(t *testing.T) {
	ctx := context.Background()
	resolver := NewResolver(Options{})
	pair, err := resolver.Issue(ctx, map[string]any{"a": "b"}, true)
	if !assert.NoError(t, err) {
		return
	}
	t.Logf("%+v", pair)
}

func TestResolver_VerifyAccess(t *testing.T) {
	ctx := context.Background()
	resolver := NewResolver(Options{})
	pair, err := resolver.Issue(ctx, map[string]any{"a": "b"}, false)
	if !assert.NoError(t, err) {
		return
	}
	access, err := resolver.VerifyAccess(ctx, pair.Access.Raw)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.EqualValues(t, access.Raw, pair.Access.Raw) {
		return
	}
	t.Logf("%+v", access)
}

func TestResolver_VerifyRefresh(t *testing.T) {
	ctx := context.Background()
	resolver := NewResolver(Options{})
	pair, err := resolver.Issue(ctx, map[string]any{"a": "b"}, true)
	if !assert.NoError(t, err) {
		return
	}
	access, err := resolver.VerifyAccess(ctx, pair.Access.Raw)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.EqualValues(t, access.Raw, pair.Access.Raw) {
		return
	}
	refresh, err := resolver.VerifyRefresh(ctx, pair.Refresh.Raw)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.EqualValues(t, refresh.Raw, pair.Refresh.Raw) {
		return
	}
	t.Logf("%+v", access)
	t.Logf("%+v", refresh)
}

func TestResolver_Refresh(t *testing.T) {
	ctx := context.Background()
	resolver := NewResolver(Options{})
	pair, err := resolver.Issue(ctx, map[string]any{"a": "b"}, true)
	if !assert.NoError(t, err) {
		return
	}
	access, err := resolver.VerifyAccess(ctx, pair.Access.Raw)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.EqualValues(t, access.Raw, pair.Access.Raw) {
		return
	}
	refresh, err := resolver.VerifyRefresh(ctx, pair.Refresh.Raw)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.EqualValues(t, refresh.Raw, pair.Refresh.Raw) {
		return
	}

	for range 4 {
		// refresh pair
		newPair, err := resolver.Refresh(ctx, pair.Access.Raw, pair.Refresh.Raw)
		if !assert.NoError(t, err) {
			return
		}

		newAccess, err := resolver.VerifyAccess(ctx, newPair.Access.Raw)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.EqualValues(t, newAccess.Raw, newPair.Access.Raw) {
			return
		}
		newRefresh, err := resolver.VerifyRefresh(ctx, newPair.Refresh.Raw)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.EqualValues(t, newRefresh.Raw, newPair.Refresh.Raw) {
			return
		}

		pair = newPair
		access = newAccess
		refresh = newRefresh
	}
}
