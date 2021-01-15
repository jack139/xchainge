package types

import (
	uuid "github.com/satori/go.uuid"
)

// 交易信息
type Deal struct {
	ID         uuid.UUID // 交易ID
	AssetsID   []byte //资产ID
	ExchangeID [32]byte //交易所的加密公钥
	Data       []byte // 加密交易数据（例如 ipfs hash）
	Refer      []byte // 参考字符串，用于索引
	Action     byte // 0x01 买入， 0x02 卖出， 0x03 变更所有权
}

// GetKey 获取实体键
func (deal *Deal) GetKey() []byte {
	return deal.ID.Bytes()
}

func (deal *Deal) getSignBytes() []byte {
	return deal.ID[:]
}
