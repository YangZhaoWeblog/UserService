//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	diy_log "qqShuiHu/internal/log"
	"qqShuiHu/internal/user/biz"
	"qqShuiHu/internal/user/conf"
	"qqShuiHu/internal/user/data"
	"qqShuiHu/internal/user/server"
	"qqShuiHu/internal/user/service"

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
