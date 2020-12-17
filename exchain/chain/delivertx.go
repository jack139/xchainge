package chain

/*
	交易上链处理
*/



import (
	"encoding/json"

	"github.com/tendermint/tendermint/abci/types"
)


/*
	提交区块
*/
func (app *App) DeliverTx(req types.RequestDeliverTx) (rsp types.ResponseDeliverTx) {
	app.logger.Info("DeliverTx()", "para", req.Tx)

	db := app.state.db
	var m TxReq

	err := json.Unmarshal(req.Tx, &m)
	if err != nil {
		rsp.Log = "bad json format"
		rsp.Code = 1
		return
	}

	switch m.Action {
	case 0x01, 0x02, 0x03, 0x04, 0x05: // 
		rsp.Log = actionMessage[m.Action]

		// 生成资产key
		assetsLinkKey := assetsPrefixKey(m.AssetsId)
		assetsLinkValue := Int64ToByteArray(app.state.Height+1)

		// 生成blcok链表key
		blockLinkKey := blockPrefixKey(app.state.Height+1)		
		blockLinkValue := FindKey(db, assetsLinkKey)

		app.logger.Info("?", string(assetsLinkKey), assetsLinkValue)
		app.logger.Info("?", string(blockLinkKey), blockLinkValue)

		// 添加到 db
		AddKV(db, assetsLinkKey, assetsLinkValue) 
		AddKV(db, blockLinkKey, blockLinkValue) 

	default:
		rsp.Log = "weird command"
		rsp.Code = 2
	}

	app.logger.Info("DeliverTx()", "action", rsp.Log)

	return
}

