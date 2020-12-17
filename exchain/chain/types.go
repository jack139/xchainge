package chain


import (
	"time"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// 一些参数
var (
	stateKey        = []byte("stateKey")
	assetsLinkPrefixKey = []byte("assetsLink:")
	blockLinkPrefixKey = []byte("blockLink:")

	ProtocolVersion uint64 = 0x1

	actionMessage = []string{
		"", 
		"buy in", 
		"buy out", 
		"change owner", 
		"authorization",
		"query",
	}
)

// 保存应用状态使用
type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

// 应用的结构
type App struct {
	types.BaseApplication

	state State
	RetainBlocks int64 // blocks to retain after commit (via ResponseCommit.RetainHeight)

	logger log.Logger
}


// 交易请求数据
type TxReq struct {
	ExchangerId string `json:"exid"` // 交易所id
	AssetsId    string `json:"aid"` // 资产id
	Category    string `json:"cat"` // 资产类别
	DataHash    string `json:"dhash"` // 数据实体hash
	UserId      string `json:"uid"` // 第三方用户id
	Action      byte   `json:"act"` // 0x01 买入， 0x02 卖出， 0x03 变更所有权， 
									// 0x04 授权查询， 0x05 查询资产
}

// 查询请求数据
type QueryReq struct {
	Query  string `json:"query"` // 视act，含义不同： 资产id（0x01）, 交易所id（0x02）, 
								 //第三方用户id（0x03）
	Action byte   `json:"act"`  // 0x01 查询资产历史, 0x02 查询交易所交易, 
								// 0x03 查询第三方用户交易
}


// query返回： 0x01 资产历史
type RespAssetsHistory struct {
	TxRequest TxReq `json:"tx_data"`
	BlockTime time.Time `json:"time"` 
}
