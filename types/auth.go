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
type Auth struct {
	ID             uuid.UUID
	AssetsID       []byte //资产ID
	FromExchangeID [32]byte //交易所的加密公钥
	ToExchangeID   [32]byte //被授权的交易所的加密公钥
	Refer          []byte // 参考字符串，用于索引
	Action         byte // 0x04 授权查询， 0x05 查询资产， 0x06 取消授权查询
}

// GetKey 获取实体键
func (auth *Auth) GetKey() []byte {
	return auth.ID.Bytes()
}

func (auth *Auth) getSignBytes() []byte {
	return auth.ID[:]
}