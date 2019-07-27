package logic

import (
	"context"
	"time"

	"micro-me/application/common/baseerror"
	"micro-me/application/common/config"
	"micro-me/application/gateway/models"
	"micro-me/application/imserver/protos"
	"micro-me/application/userserver/protos"
)

type (
	GateWayLogic struct {
		userRpcModel user.UserService
		gateWayModel *models.GateWayModel
		imRpcModel   im.ImService

		imAddressList []*config.ImRpc
	}
	SendRequest struct {
		FromToken string    `json:"fromToken"  binding:"required"`
		ToToken   string    `json:"toToken"  binding:"required"`
		Body      string    `json:"body"  binding:"required"`
		Timestamp time.Time `json:"timestamp"`
	}

	SendResponse struct {
	}

	GetServerAddressRequest struct {
		Token string `json:"token" binding:"required"`
	}

	GetServerAddressResponse struct {
		Address string `json:"address"`
	}
)

var (
	SendMessageErr    = baseerror.NewBaseError("发送消息失败")
	UserNotFoundErr   = baseerror.NewBaseError("用户不存在")
	ImAddressErr      = baseerror.NewBaseError("请配置消息服务地址")
	AddDataErr        = baseerror.NewBaseError("维护关系错误")
	PublishMessageErr = baseerror.NewBaseError("发送消息到MQ失败")
	imRpcModelMapErr  = baseerror.NewBaseError("没有找到对应的RPC服务")
)

func NewGateWayLogic(userRpcModel user.UserService,
	gateWayModel *models.GateWayModel,
	imAddressList []*config.ImRpc,
	imRpcModel im.ImService,
) *GateWayLogic {

	return &GateWayLogic{
		userRpcModel:  userRpcModel,
		gateWayModel:  gateWayModel,
		imAddressList: imAddressList,
		imRpcModel:    imRpcModel,
	}
}

func (l *GateWayLogic) Send(r *SendRequest) (*SendResponse, error) {

	if _, err := l.userRpcModel.FindByToken(context.TODO(), &user.FindByTokenRequest{Token: r.ToToken}); err != nil {
		return nil, UserNotFoundErr
	}
	userGate, err := l.gateWayModel.FindByToken(r.ToToken)
	if err != nil {
		return nil, SendMessageErr
	}
	if userGate.Id < 0 {
		return nil, SendMessageErr
	}
	req := &im.PublishMessageRequest{
		FromToken:  r.FromToken,
		ToToken:    r.ToToken,
		Body:       r.Body,
		ServerName: userGate.ServerName,
		Topic:      userGate.Topic,
		Address:    userGate.ImAddress,
	}
	_, err = l.imRpcModel.PublishMessage(context.TODO(), req);
	// 发送消息逻辑
	if  err != nil {
		return nil, PublishMessageErr
	}
	// 发送消息逻辑结束
	return &SendResponse{}, nil
}

func (l *GateWayLogic) GetServerAddress(r *GetServerAddressRequest) (*GetServerAddressResponse, error) {

	u, err := l.userRpcModel.FindByToken(context.TODO(), &user.FindByTokenRequest{Token: r.Token});
	if err != nil {
		return nil, UserNotFoundErr
	}
	length := len(l.imAddressList)
	if length == 0 {
		return nil, ImAddressErr
	}
	i := u.Id % int64(length)
	imData := l.imAddressList[int(i)]

	if _, err := l.gateWayModel.Insert(&models.GateWay{
		Topic:      imData.Topic,
		Token:      r.Token,
		ImAddress:  imData.Address,
		ServerName: imData.ServerName,
	}); err != nil {
		return nil, AddDataErr
	}
	return &GetServerAddressResponse{
		Address: imData.Address,
	}, nil
}
