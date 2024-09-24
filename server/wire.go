//go:build wireinject

// The build tag makes sure the stub is not built in the final build
package server

import (
	"github.com/ginx-contribs/ginx-server/server/api"
	"github.com/ginx-contribs/ginx-server/server/data"
	"github.com/ginx-contribs/ginx-server/server/handler"
	"github.com/ginx-contribs/ginx-server/server/svc"
	"github.com/ginx-contribs/ginx-server/server/types"
	"github.com/google/wire"
)

// initialize and setup app environment
func setup(ctx types.Context) (svc.Context, error) {
	panic(wire.Build(ContextProvider, data.Provider, handler.Provider, api.Provider, svc.Provider))
}
