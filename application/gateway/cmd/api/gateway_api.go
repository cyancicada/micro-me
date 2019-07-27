package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"

	"micro-me/application/common/middleware"
	gateWayConfig "micro-me/application/gateway/cmd/config"
	"micro-me/application/gateway/controller"
	"micro-me/application/gateway/logic"
	"micro-me/application/gateway/models"
	imProto "micro-me/application/imserver/protos"
	"micro-me/application/userserver/protos"
)

func main() {
	userRpcFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config_api.json",
		Usage: "please use xxx -f config_rpc.json",
	}
	configFile := flag.String(userRpcFlag.Name, userRpcFlag.Value, userRpcFlag.Usage)
	flag.Parse()
	conf := new(gateWayConfig.ApiConfig)

	if err := config.LoadFile(*configFile); err != nil {
		log.Fatal(err)
	}
	if err := config.Scan(conf); err != nil {
		log.Fatal(err)
	}
	engineGateWay, err := xorm.NewEngine(conf.Engine.Name, conf.Engine.DataSource)
	if err != nil {
		log.Fatal(err)
	}
	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = conf.Etcd.Address
		});

	// Create a new service. Optionally include some options here.
	rpcService := micro.NewService(
		micro.Name(conf.Server.Name),
		micro.Registry(etcdRegisty),
		micro.Transport(grpc.NewTransport()),
		micro.WrapClient(hystrix.NewClientWrapper()),
		micro.Flags(userRpcFlag),
	)
	rpcService.Init()
	userRpcModel := user.NewUserService(conf.UserRpcServer.ServerName, rpcService.Client())

	imRpcModel := imProto.NewImService(conf.ImRpcServer.ServerName, rpcService.Client())
	gateWayModel := models.NewGateWayModel(engineGateWay)
	gateLogic := logic.NewGateWayLogic(userRpcModel, gateWayModel, conf.ImRpcServer.ImServerList,imRpcModel)
	gateWayController := controller.NewGateController(gateLogic)
	service := web.NewService(
		web.Name(conf.Server.Name),
		web.Registry(etcdRegisty),
		web.Version(conf.Version),
		web.Flags(userRpcFlag),
		web.Address(conf.Port),
		web.Flags(userRpcFlag),
	)
	router := gin.Default()

	userRouterGroup := router.Group("/gateway")
	userRouterGroup.Use(middleware.ValidAccessToken)
	{
		userRouterGroup.POST("/send", gateWayController.Send)
		userRouterGroup.POST("/address", gateWayController.GetServerAddress)
	}
	service.Handle("/", router)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
