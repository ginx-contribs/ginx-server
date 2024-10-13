//go:build wireinject

package wirex

import (
	"github.com/ginx-contribs/ginx-server/internal/common/types"
	"github.com/ginx-contribs/ginx-server/internal/modules"
	"github.com/google/wire"
)

func Inject(injector types.Injector) (modules.Modules, error) {
	panic(wire.Build(types.Provider, modules.Provider))
}
