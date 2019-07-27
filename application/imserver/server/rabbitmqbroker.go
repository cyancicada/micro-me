package server

import (
	"log"

	"github.com/micro/go-micro/broker"
)

type (
	RabbitMqBroker struct {
		topic          string
		rabbitMqBroker broker.Broker
	}
)

func NewRabbitMqBroker(topic string, rabbitMqBroker broker.Broker) (*RabbitMqBroker, error) {

	// 初始化
	if err := rabbitMqBroker.Init(); err != nil {
		return nil, err
	}
	if err := rabbitMqBroker.Connect(); err != nil {
		return nil, err
	}
	return &RabbitMqBroker{topic: topic, rabbitMqBroker: rabbitMqBroker}, nil
}

//发送者要写在rpc里，被网关调用
func (p *RabbitMqBroker) Publisher(msg *broker.Message) {
	if err := p.rabbitMqBroker.Publish(p.topic, msg); err != nil {
		log.Printf("[publisher %s err] : %+v", p.topic, err)
	}
	log.Printf("[publisher %s] : %s", p.topic, string(msg.Body))
}

//发送者要写在rpc里，被网关调用
func (p *RabbitMqBroker) Subscribe(handlerFunc func(msg []byte) error) {
	if _, err := p.rabbitMqBroker.Subscribe(p.topic, func(publication broker.Publication) error {
		if err := handlerFunc(publication.Message().Body);err != nil{
			log.Println("handlerFunc msg err %+v",err)
		}
		return nil
	}); err != nil {
		log.Printf("[Subscribe %s err] : %+v", p.topic, err)
	}
	log.Printf("[publisher err]")
}
