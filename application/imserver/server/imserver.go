package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/broker"

	"micro-me/application/common/baseerror"
)

type (
	ImServer struct {
		rabbitMqBroker *RabbitMqBroker
		clients        map[string]*websocket.Conn
		Address        string
		lock           sync.Mutex
		upgraer        *websocket.Upgrader
	}
	SendMsgRequest struct {
		FromToken     string `json:"fromToken"`
		ToToken       string `json:"toToken"`
		Body          string `json:"body"`
		TimeStamp     int64  `json:"timeStamp"`
		RemoteAddress string `json:"remoteAddress"`
	}

	LoginRequest struct {
		Token string `json:"token"`
	}
	SendMsgResponse struct {
		FromToken     string `json:"fromToken"`
		Body          string `json:"body"`
		RemoteAddress string `json:"remoteAddress"`
	}
	ImServerOptions func(im *ImServer)
)

var (
	DefaultAddress  = ":7272"
	WebSocketPrefix = "/ws"
	UserNoLoginErr  = baseerror.NewBaseError("此用户没有登录！")
	SendMessageErr  = baseerror.NewBaseError("发送消息失败！")
)

func NewImServer(rabbitMqBroker *RabbitMqBroker, opts ImServerOptions) (*ImServer, error) {
	// 初始化
	if err := broker.Init(); err != nil {
		return nil, err
	}
	if err := broker.Connect(); err != nil {
		return nil, err
	}
	imServer := &ImServer{
		rabbitMqBroker: rabbitMqBroker,
		clients:        make(map[string]*websocket.Conn, 0),
		upgraer: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	if opts != nil {
		opts(imServer)
	}
	if imServer.Address == "" {
		imServer.Address = DefaultAddress
	}
	return imServer, nil
}

func (l *ImServer) SendMsg(r *SendMsgRequest) (*SendMsgResponse, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	log.Printf("send SendMsgRequest  %+v", r)
	conn := l.clients[r.ToToken]
	if conn == nil {
		return nil, UserNoLoginErr
	}
	r.TimeStamp = time.Now().Unix()
	r.RemoteAddress = conn.RemoteAddr().String()
	bodyMsg, err := json.Marshal(r)
	if err != nil {
		return nil, SendMessageErr
	}
	if err := conn.WriteMessage(websocket.TextMessage, bodyMsg); err != nil {
		log.Printf("send message err %v", err)
		l.clients[r.ToToken] = nil
		//log.Println(conn.Close())
		return nil, err
	}
	log.Printf("send message succes  %v", r.Body)
	return &SendMsgResponse{}, nil
}

func (l *ImServer) Subscribe() {
	l.rabbitMqBroker.Subscribe(func(msg []byte) error {
		r := new(SendMsgRequest)
		if err := json.Unmarshal(msg, r); err != nil {
			log.Printf("[Unmarshal msg err] : %+v", err)
			return err
		}
		if _, err := l.SendMsg(r); err != nil {
			log.Printf("[SendMsg err] : %+v", err)
			return err
		}
		log.Printf("has Subscribe msg %+v", string(msg))
		return nil
	})
}

func (l *ImServer) Run() {
	log.Printf("websocket has listens at %s", l.Address)
	http.HandleFunc(WebSocketPrefix, l.login)
	log.Fatal(http.ListenAndServe(l.Address, nil))
}

func (l *ImServer) login(w http.ResponseWriter, r *http.Request) {
	conn, err := l.upgraer.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	msgType, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("read login message err %+v", err)
		return
	}
	if msgType != websocket.TextMessage {
		log.Printf("read login msgType err %+v", err)
		return
	}
	fmt.Println(string(message))
	loginMsgRequest := new(LoginRequest)

	if err := json.Unmarshal(message, loginMsgRequest); err != nil {
		log.Printf("json.Unmarshal msg err %+v", err)
		return
	}
	l.clients[loginMsgRequest.Token] = conn
	return
}
