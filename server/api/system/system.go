package system

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx/pkg/resp"
)

func NewSystemAPI() *SystemAPI {
	return &SystemAPI{}
}

type SystemAPI struct {
}

// Ping
// @Summary      Ping
// @Description  test server if is available
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.Response{data=string}
// @Router       /ping [GET]
func (a *SystemAPI) Ping(ctx *gin.Context) {
	resp.Ok(ctx).Msg("pong").JSON()
}

// Pong
// @Summary      Pong
// @Description  test if server authentication is working
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.Response{data=string}
// @Router       /pong [GET]
func (a *SystemAPI) Pong(ctx *gin.Context) {
	resp.Ok(ctx).Msg("ping").JSON()
}
