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
	fileLinkPrefixKey = []byte("fileLink:")
	blockLinkPrefixKey = []byte("blockLink:")
	userFilePrefixKey = []byte("userFile:")

	ProtocolVersion uint64 = 0x1
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


// 文件信息 - 用户文件表使用
type FileData struct {
	FileName    string `json:"fn"` // 文件名，可为空
	Modified    bool `json:"is_mod"` // 文件是否已修改
}

// 交易请求数据
type TxReq struct {
	UserId      string `json:"uid"` // 文件主的用户id
	FileHash    string `json:"fhash"` // 文件hash
	OldFileHash string `json:"ofhash"` // 旧文件的hash (如果 action==修改 需提供)
	FileName    string `json:"fn"` // 文件名，可为空
	ReaderId    string `json:"rid"` // 浏览文件的用户id （如果 action==浏览 需提供）
	Action      byte   `json:"act"` // 0x01 文件建立， 0x02 文件浏览， 0x03 文件修改， 0x04 文件删除
}

// 查询请求数据
type QueryReq struct {
	UserId      string `json:"uid"` // 文件主的用户id，action==2时提供
	FileHash    string `json:"fhash"` // 文件hash，action==1时提供
	Action      byte   `json:"act"` // 0x01 查询文件历史, 0x02 查询用户的文件列表
}


// query返回： 文件历史
type RespFileHistory struct {
	TxRequest TxReq `json:"tx_data"`
	BlockTime time.Time `json:"time"` 
}

// query返回： 用户文件
type RespUserFile struct {
	UserId      string `json:"user_id"` 
	FileName    string `json:"filename"`
	FileHash    string `json:"file_hash"` 
	Modified    bool `json:"is_modified"` 
}
