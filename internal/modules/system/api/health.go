package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/handler"
	"github.com/ginx-contribs/ginx/pkg/resp"
)

type HealthAPI struct {
	HealthHandler handler.HealthHandler
}

// Ping
// @Summary      Ping
// @Description  ping test web service if is available
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.Response{data=string}
// @Router       /health/ping [GET]
func (h HealthAPI) Ping(ctx *gin.Context) {
	resp.Ok(ctx).
		Data("pong").
		JSON()
}
