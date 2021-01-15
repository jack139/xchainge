package types

import (
	uuid "github.com/satori/go.uuid"
)

// token 区块： Credit - CR
type Credit struct {
	ID         uuid.UUID // ID
	UserID     [32]byte //用户的加密公钥
	Data       []byte // 交易数据：不加密， 当 区块增减时，记录原因
	Action     byte // 0x01 空区块（占位用） 0x02 减少  0x03 增加
}

// GetKey 获取实体键
func (deal *Credit) GetKey() []byte {
	return deal.ID.Bytes()
}

func (deal *Credit) getSignBytes() []byte {
	return deal.ID[:]
}
