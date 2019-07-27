package rpcserveriml

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/micro/go-micro/broker"

	"micro-me/application/common/baseerror"
	"micro-me/application/imserver/protos"
	"micro-me/application/imserver/server"
)

type (
	ImRpcServerIml struct {
		sync.Mutex
		publisherServerMap map[string]*server.RabbitMqBroker
	}
)

var (
	PublishMessageErr = baseerror.NewBaseError("发送消息失败")
)

func NewImRpcServerIml(publisherServerMap map[string]*server.RabbitMqBroker) *ImRpcServerIml {

	return &ImRpcServerIml{publisherServerMap: publisherServerMap}
}

func (s *ImRpcServerIml) PublishMessage(ctx context.Context, req *im.PublishMessageRequest, rsp *im.PublishMessageResponse) error {
	body, err := json.Marshal(req)
	if err != nil {
		return PublishMessageErr
	}
	key := req.ServerName + req.Topic
	publisher := s.publisherServerMap[key]
	publisher.Publisher(&broker.Message{
		Body: body,
	})
	return nil
}
