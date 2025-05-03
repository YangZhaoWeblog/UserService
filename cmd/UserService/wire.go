//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/YangZhaoWeblog/UserService/internal/applog"
	"github.com/YangZhaoWeblog/UserService/internal/biz"
	"github.com/YangZhaoWeblog/UserService/internal/conf"
	"github.com/YangZhaoWeblog/UserService/internal/data"
	"github.com/YangZhaoWeblog/UserService/internal/pkg"
	"github.com/YangZhaoWeblog/UserService/internal/server"
	"github.com/YangZhaoWeblog/UserService/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.App, *conf.Log, *conf.Data) (*kratos.App, func(), error) {
	panic(wire.Build(
		pkg.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		applog.ProviderSet,
		server.ProviderSet,
		newApp,
	))
}
