package types

import (
	"time"
)

// query返回数据结构
type RespQuery struct {
	Type           string // deal or auth
	ID             string //  v       v
	ExchangeId     string //  v       v
	AuthExchangeId string //  x       v
	AssetsId       string //  v       x
	Data           string //  v       ?   auth：响应授权时返回解密的 Data
	Refer          string //  ?       v   auth：这里返回的是对应的 DealID
	Action         byte   //  v       v
	SendTime    time.Time //  v       v
}
