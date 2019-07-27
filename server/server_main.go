package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	rl "github.com/juju/ratelimit"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/transport/grpc"
	"github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"

	hello "micro-me/protos"
)

type Say struct {
	Tag string
}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	time.Sleep(5 * time.Second)
	rsp.Msg = "Hello " + req.Name + " [FROM " + s.Tag + "]"
	return nil
}

type (
	//"Version": "0.0.1",
	//  "Hello" : {
	//    "Name": "hello"
	//  },
	//  "Etcd": {
	//    "Addrs": ["192.168.5.100:2379"],
	//    "UserName": "",
	//    "Password": ""
	//  }
	Config struct {
		Version string
		Hello   struct {
			Name string
		}
		Etcd struct {
			Addrs    []string
			UserName string
			Password string
		}
	}
)

func main() {
	configFile := flag.String("f", "./config/config.json", "please use config.json")
	conf := new(Config)

	if err := config.LoadFile(*configFile); err != nil {
		log.Fatal(err)
	}
	if err := config.Scan(conf); err != nil {
		log.Fatal(err)
	}
	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = conf.Etcd.Addrs
			//etcdv3.Auth("root","1234")(options)
		});
	md := map[string]string{
		"vector": "yang",
	}
	limit := 2
	b := rl.NewBucketWithRate(float64(limit), int64(limit))
	micro.Metadata(md)
	service := micro.NewService(
		micro.Name(conf.Hello.Name),
		micro.Registry(etcdRegisty),
		micro.Version(conf.Version),
		micro.Metadata(md),
		micro.Transport(grpc.NewTransport()),
		micro.WrapHandler(ratelimit.NewHandlerWrapper(b, false)),
	)

	service.Init()

	say := &Say{
		Tag: strconv.Itoa(rand.Int()),
	}
	fmt.Println("当前rpc服务的TAG 为" + say.Tag)
	hello.RegisterGreeterHandler(service.Server(), say)

	// 初始化
	if err := broker.Init(); err != nil {
		log.Fatal(err)
	}
	if err := broker.Connect(); err != nil {
		log.Fatal(err)
	}
	go publisher()
	go subscribe()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

var topic = "demo.topic"

func publisher() {
	t := time.NewTicker(time.Second)
	for e := range t.C {
		msg := &broker.Message{
			Header: map[string]string{
				"Tag": strconv.Itoa(rand.Int()),
			},
			Body: []byte(e.String()),
		}
		if err := broker.Publish(topic, msg); err != nil {

			log.Printf("[publisher err] : %+v", err)
		}
	}
}

func subscribe() {

	if _, err := broker.Subscribe(topic, func(publication broker.Publication) error {

		fmt.Printf("subscribe received msg : %s,Header is %+v",
			string(publication.Message().Body),
			publication.Message().Header,
		)
		fmt.Println()
		return nil
	}); err != nil {
		log.Printf("[subscribe err] : %+v", err)
	}
}
