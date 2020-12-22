package types

// query返回数据结构
type RespQuery struct {
	Type           string // deal or auth
	ID             string //  v       v
	ExchangeId     string //  v       x
	AuthExchangeId string //  x       v
	AssetsId       string //  v       v
	Data           string //  v       x
	Refer          string //  ?       ?
	Action         byte   //  v       v
}
