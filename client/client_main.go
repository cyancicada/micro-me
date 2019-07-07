package main

import (
	"context"
	"fmt"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"

	proto "micro-me/protos"
)

type MyClientWrapper struct {
	client.Client
}

func (c *MyClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(e error) error {
		fmt.Println("这是一个备用的服务")
		return nil
	})
}

// NewClientWrapper returns a hystrix client Wrapper.
func NewMyClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &MyClientWrapper{c}
	}
}
func main() {

	//hystrix.DefaultTimeout = 6000
	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = []string{"192.168.5.100:2379"}
			//etcdv3.Auth("root","1234")(options)
		})

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("greeter.client"),
		micro.Registry(etcdRegisty),
		//micro.Transport(grpc.NewTransport()),
		micro.WrapClient(NewMyClientWrapper()),
	)
	service.Init()

	// Create new greeter client
	greeter := proto.NewGreeterService("greeter", service.Client())

	t := time.NewTicker(100 * time.Millisecond)

	for e := range t.C {
		// Call the greeter
		rsp, err := greeter.Hello(context.TODO(), &proto.Request{Name: "John"})
		if err != nil {
			fmt.Printf("err=>%+v [%+v]", err, e)
		} else {
			fmt.Printf("msg=>%+v [%+v]", rsp.Msg, e)
		}
		// Print response
		fmt.Println()
	}

}
