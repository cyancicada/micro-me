package main

import (
	"flag"
	"log"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/micro/go-plugins/registry/etcdv3"

	imConfig "micro-me/application/imserver/cmd/config"
	proto "micro-me/application/imserver/protos"
	"micro-me/application/imserver/rpcserveriml"
	"micro-me/application/imserver/server"
)

func main() {
	imFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config_rpc.json",
		Usage: "please use xxx -f config_rpc.json",
	}
	configFile := flag.String(imFlag.Name, imFlag.Value, imFlag.Usage)
	flag.Parse()
	conf := new(imConfig.ImRpcConfig)

	if err := config.LoadFile(*configFile); err != nil {
		log.Fatal(err)
	}
	if err := config.Scan(conf); err != nil {
		log.Fatal(err)
	}
	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = conf.Etcd.Address
		});
	service := micro.NewService(
		micro.Name(conf.Server.Name),
		micro.Registry(etcdRegisty),
		micro.Version(conf.Version),
		micro.Transport(grpc.NewTransport()),
		micro.Flags(imFlag),
	)
	publisherServerMap := make(map[string]*server.RabbitMqBroker)
	for _, item := range conf.ImServerList {
		amqbAddress := item.AmqbAddress
		p, err := server.NewRabbitMqBroker(
			item.Topic,
			rabbitmq.NewBroker(func(options *broker.Options) {
				options.Addrs = amqbAddress
			}),
		)
		if err != nil {
			log.Fatal(err)
		}
		publisherServerMap[item.ServerName+item.Topic] = p
	}
	imRpcServer := rpcserveriml.NewImRpcServerIml(publisherServerMap)
	if err := proto.RegisterImHandler(service.Server(), imRpcServer); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
