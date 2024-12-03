//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"JoeTiktok/internal/user/biz"
	"JoeTiktok/internal/user/conf"
	"JoeTiktok/internal/user/data"
	"JoeTiktok/internal/user/server"
	"JoeTiktok/internal/user/service"
	diy_log "JoeTiktok/pkg/log"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Log) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		newApp,
		diy_log.NewLogger,
	))
}
