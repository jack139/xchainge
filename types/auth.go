package types

import (
	uuid "github.com/satori/go.uuid"
)

const (
	AuthQuery byte = 0x04
	DoQuery byte = 0x05
	DeAuthQuery byte = 0x06
)

// 授权操作
// ToExchangeID 请求 FromExchangeID 授权，指定 DealID，进入 FromExchangeID 的链表
// FromExchangeID 加密数据 Data 后 返回 ToExchangeID，进入 ToExchangeID 的链表
type Auth struct {
	ID             uuid.UUID
	DealID         uuid.UUID // 交易ID
	FromExchangeID [32]byte // 交易所的加密公钥
	ToExchangeID   [32]byte // 被授权的交易所的加密公钥
	Data           []byte // FromExchange加密数据，被授权者ToExchangeID可以解密
	Action         byte // 0x04 请求授权， 0x05 响应授权
}

// GetKey 获取实体键
func (auth *Auth) GetKey() []byte {
	return auth.ID.Bytes()
}

func (auth *Auth) getSignBytes() []byte {
	return auth.ID[:]
}