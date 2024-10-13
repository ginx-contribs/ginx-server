package system

import (
	"github.com/ginx-contribs/ginx-server/internal/common/types"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/api"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/handler"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/repo"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	// repo
	wire.Struct(new(repo.UserRepo), "*"),
	// handler
	handler.NewEmailHandler,
	wire.Struct(new(handler.AuthHandler), "*"),
	wire.Struct(new(handler.CaptchaHandler), "*"),
	wire.Struct(new(handler.UserHandler), "*"),
	wire.Struct(new(handler.HealthHandler), "*"),
	// api
	wire.Struct(new(api.AuthAPI), "*"),
	wire.Struct(new(api.UserAPI), "*"),
	wire.Struct(new(api.HealthAPI), "*"),

	// module
	wire.Struct(new(Module), "*"),
)

// Module is representation of system
type Module struct {
	// api
	AuthAPI   api.AuthAPI
	UserAPI   api.UserAPI
	HealthAPI api.HealthAPI

	// handler
	AuthHandler   handler.AuthHandler
	CodeHandler   handler.CaptchaHandler
	EmailHandler  handler.EmailHandler
	UserHandler   handler.UserHandler
	HealthHandler handler.HealthHandler

	// repo
	UserRepo repo.UserRepo
}

func (m Module) Name() string {
	return "system"
}

func (m Module) Init(injector types.Injector) error {
	m.RegisterRouter(injector)
	return nil
}

func (m Module) Close() error {
	return nil
}

func (m Module) RegisterRouter(injector types.Injector) {
	router := injector.Router
	// auth api
	authAPI := m.AuthAPI
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authAPI.Login)
		authGroup.POST("/register", authAPI.Register)
		authGroup.POST("/reset", authAPI.ResetPassword)
		authGroup.POST("/refresh", authAPI.Refresh)
		authGroup.POST("/captcha", authAPI.VerifyCode)
	}

	// user api
	userAPI := m.UserAPI
	userGroup := router.Group("/user")
	{
		userGroup.GET("/info", userAPI.Info)
		userGroup.GET("/me", userAPI.Me)
		userGroup.GET("/list", userAPI.List)
	}

	// health api
	healthAPI := m.HealthAPI
	healthGroup := router.Group("/health")
	{
		healthGroup.GET("/ping", healthAPI.Ping)
	}
}
