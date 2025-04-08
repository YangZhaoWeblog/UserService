//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/YangZhaoWeblog/UserService/internal/biz"
	"github.com/YangZhaoWeblog/UserService/internal/conf"
	"github.com/YangZhaoWeblog/UserService/internal/data"
	"github.com/YangZhaoWeblog/UserService/internal/server"
	"github.com/YangZhaoWeblog/UserService/internal/service"

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
		//applog.ProviderSet,
		newApp,
	))
}
