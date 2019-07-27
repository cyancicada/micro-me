package models

import (
	"time"

	"github.com/go-xorm/xorm"
)

type (
	GateWay struct {
		Id         int64
		Token      string    `json:"token" xorm:"varchar(11) notnull 'token'"`
		ImAddress  string    `json:"imAddress" xorm:"varchar(60) notnull 'im_address'"`
		ServerName string    `json:"server_name" xorm:"varchar(60) notnull 'server_name'"`
		Topic      string    `json:"topic" xorm:"varchar(60) notnull 'topic'"`
		CreateTime time.Time `json:"createTime" xorm:"DateTime 'create_time'"`
		UpdateTime time.Time `json:"updateTime" xorm:"DateTime 'update_time'"`
	}
	GateWayModel struct {
		mysql *xorm.Engine
	}
)

func (g *GateWay) TableName() string {
	return "gateway"
}
func NewGateWayModel(mysql *xorm.Engine) *GateWayModel {

	return &GateWayModel{mysql: mysql}
}

func (m *GateWayModel) FindByToken(token string) (*GateWay, error) {
	g := new(GateWay)

	if _, err := m.mysql.Where("token = ?", token).Get(g); err != nil {
		return nil, err
	}
	return g, nil
}

func (m *GateWayModel) FindByServerNameTokenAddressTopic(serverName, topic, token, address string) (*GateWay, error) {
	g := new(GateWay)

	if _, err := m.mysql.Where(
		"token = ? and im_address =? and topic = ? and server_name=?",
		token,
		address,
		topic,
		serverName,
	).Get(g); err != nil {
		return nil, err
	}
	return g, nil
}
func (m *GateWayModel) Insert(g *GateWay) (*GateWay, error) {

	has, err := m.FindByServerNameTokenAddressTopic(g.ServerName, g.Topic, g.Token, g.ImAddress)
	if has != nil && has.Id > 0 && err == nil {
		return has,nil
	}
	if _, err := m.mysql.Insert(g); err != nil {
		return nil, err
	}
	return g, nil
}
func (m *GateWayModel) FindByImAddress(imAddress string) ([]*GateWay, error) {
	gs := []*GateWay(nil)

	if err := m.mysql.Where("im_address = ?", imAddress).Find(&gs); err != nil {
		return nil, err
	}
	return gs, nil
}
