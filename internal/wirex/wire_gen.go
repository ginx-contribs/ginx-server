// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wirex

import (
	"github.com/ginx-contribs/ginx-server/internal/common/data/cache"
	"github.com/ginx-contribs/ginx-server/internal/common/types"
	"github.com/ginx-contribs/ginx-server/internal/modules"
	"github.com/ginx-contribs/ginx-server/internal/modules/system"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/api"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/handler"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/repo"
)

// Injectors from wire.go:

func Inject(injector types.Injector) (modules.Modules, error) {
	resolver := injector.Token
	client := injector.EntDB
	userRepo := repo.UserRepo{
		DB: client,
	}
	redisClient := injector.Redis
	redisCodeCache := cache.NewRedisCaptchaCache(redisClient)
	app := injector.Config
	email := app.Email
	sender := injector.Email
	queue := injector.MQ
	emailHandler, err := handler.NewEmailHandler(email, sender, queue)
	if err != nil {
		return modules.Modules{}, err
	}
	metaInfo := app.Meta
	captchaHandler := handler.CaptchaHandler{
		CaptchaCache: redisCodeCache,
		EmailHandler: emailHandler,
		MetaInfo:     metaInfo,
	}
	authHandler := handler.AuthHandler{
		Token:          resolver,
		UserRepo:       userRepo,
		CaptchaHandler: captchaHandler,
	}
	authAPI := api.AuthAPI{
		TokenResolver:  resolver,
		AuthHandler:    authHandler,
		CaptchaHandler: captchaHandler,
	}
	repoUserRepo := &repo.UserRepo{
		DB: client,
	}
	userHandler := handler.UserHandler{
		UserRepo: repoUserRepo,
	}
	userAPI := api.UserAPI{
		UserHandler: userHandler,
	}
	healthHandler := handler.HealthHandler{}
	healthAPI := api.HealthAPI{
		HealthHandler: healthHandler,
	}
	module := system.Module{
		AuthAPI:       authAPI,
		UserAPI:       userAPI,
		HealthAPI:     healthAPI,
		AuthHandler:   authHandler,
		CodeHandler:   captchaHandler,
		EmailHandler:  emailHandler,
		UserHandler:   userHandler,
		HealthHandler: healthHandler,
		UserRepo:      userRepo,
	}
	modulesModules := modules.Modules{
		System: module,
	}
	return modulesModules, nil
}
