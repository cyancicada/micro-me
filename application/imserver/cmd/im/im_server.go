package main

import (
	"flag"
	"log"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/micro/go-plugins/registry/etcdv3"

	imConfig "micro-me/application/imserver/cmd/config"
	"micro-me/application/imserver/server"
)

func main() {
	imFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config_im.json",
		Usage: "please use xxx -f config_im.json",
	}
	configFile := flag.String(imFlag.Name, imFlag.Value, imFlag.Usage)
	flag.Parse()
	conf := new(imConfig.ImConfig)

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
	rabbitMqRegistry := rabbitmq.NewBroker(func(options *broker.Options) {
		options.Addrs = conf.RabbitMq.Address
	})
	service := micro.NewService(
		micro.Name(conf.Server.Name),
		micro.Registry(etcdRegisty),
		micro.Version(conf.Version),
		micro.Flags(imFlag),
	)

	log.Printf("has start listen topic %s",conf.RabbitMq.Topic)
	rabbitMqBroker, err := server.NewRabbitMqBroker(conf.RabbitMq.Topic, rabbitMqRegistry)
	if err != nil {
		log.Fatal(err)
	}
	imServer, err := server.NewImServer(rabbitMqBroker,
		func(im *server.ImServer) {
			im.Address = conf.Port
		})
	if err != nil {
		log.Fatal(err)
	}
	go imServer.Subscribe()
	go imServer.Run()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
