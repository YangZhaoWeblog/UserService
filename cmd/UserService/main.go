package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/YangZhaoWeblog/UserService/internal/conf"

	"github.com/go-kratos/kratos/v2" // 确保 kratos v2 核心包导入
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
)

const startPrint = `

           .___________. __   __  ___ .___________.  ______    __  ___
           |           ||  | |  |/  / |           | /  __  \  |  |/  /
           .---|  |----.|  | |  '  /  .---|  |----.|  |  |  | |  '  /
               |  |     |  | |    <       |  |     |  |  |  | |    <
               |  |     |  | |  .  \      |  |     |  .--'  | |  .  \
               |__|     |__| |__|\__\     |__|      \______/  |__|\__\


                __    __       _______. _______ .______
               |  |  |  |     /       ||   ____||   _  \
               |  |  |  |    |   (----.|  |__   |  |_)  |
               |  |  |  |     \   \    |   __|  |      /
               |  .--'  | .----)   |   |  |____ |  |\  \----.
                \______/  |_______/    |_______|| _| .._____|
`

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string

	configPath string
	id, _      = os.Hostname()
)

func init() {
	// 从环境变量加载配置
	configMode := os.Getenv("CONFIG_MODE")
	if configMode == "" {
		os.Exit(1)
	}
	configPath = filepath.Join("configs", configMode+".user.config.yaml")
}

func newApp(gs *grpc.Server, hs *http.Server, logger log.Logger) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()

	_, err := filepath.Abs(configPath)
	if err != nil {
		panic(err)
	}

	c := config.New(
		config.WithSource(
			file.NewSource(configPath),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.App, bc.Log, bc.Data, bc.Trace)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	fmt.Print(startPrint)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
