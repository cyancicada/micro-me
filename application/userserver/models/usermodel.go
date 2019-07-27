package models

import (
	"time"

	"github.com/go-xorm/xorm"
)

type (
	Members struct {
		Id         int64
		Token      string    `json:"token" xorm:"varchar(11) notnull 'token'"`
		Username   string    `json:"username" xorm:"varchar(60) notnull 'username'"`
		Password   string    `json:"password" xorm:"varchar(60) notnull 'password'"`
		CreateTime time.Time `json:"createTime" xorm:"DateTime 'create_time'"`
		UpdateTime time.Time `json:"updateTime" xorm:"DateTime 'update_time'"`
	}
	MembersModel struct {
		mysql *xorm.Engine
	}
)

func NewMembersModel(mysql *xorm.Engine) *MembersModel {

	return &MembersModel{
		mysql: mysql,
	}
}

func (m *MembersModel) FindByToken(token string) (*Members, error) {
	member := new(Members)

	if _, err := m.mysql.Where("token=?", token).Get(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (m *MembersModel) FindById(id int64) (*Members, error) {
	member := new(Members)
	if _, err := m.mysql.Where("id=?", id).Get(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (m *MembersModel) FindByUserName(userName string) (*Members, error) {
	member := new(Members)
	if _, err := m.mysql.Where("username=?", userName).Get(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (m *MembersModel) InsertMember(member *Members) (*Members, error) {
	if _, err := m.mysql.Insert(member); err != nil {
		return nil, err
	}
	return member, nil
}
