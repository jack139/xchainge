package chain

/*
	交易上链处理
*/



import (
	"encoding/json"

	"github.com/tendermint/tendermint/abci/types"
)


/*
	type TxReq struct {
		UserId      string
		FileHash    string
		OldFileHash string
		FileName    string
		ReaderId    string
		Action      byte  
	}
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
	case 0x01: // 新建文件
		rsp.Log = "new file"

		// 生成文件key
		fileLinkKey := filePrefixKey(m.FileHash)
		fileLinkValue := Int64ToByteArray(app.state.Height+1)

		// 生成blcok链表key
		blockLinkKey := blockPrefixKey(app.state.Height+1)		
		blockLinkValue := FindKey(db, fileLinkKey)

		// 添加到 db
		AddKV(db, fileLinkKey, fileLinkValue) 
		AddKV(db, blockLinkKey, blockLinkValue) 

		// 生成用户key, 生成file_data
		userFileKey := userPrefixKey(m.UserId, m.FileHash)
		NewFileData(db, userFileKey, "", false)


	case 0x02: // 浏览文件
		rsp.Log = "view file"

		// 生成文件key
		fileLinkKey := filePrefixKey(m.FileHash)
		fileLinkValue := Int64ToByteArray(app.state.Height+1)

		// 生成blcok链表key
		blockLinkKey := blockPrefixKey(app.state.Height+1)		
		blockLinkValue := FindKey(db, fileLinkKey)

		// 添加到 db
		AddKV(db, fileLinkKey, fileLinkValue) 
		AddKV(db, blockLinkKey, blockLinkValue) 

	case 0x03: // 修改文件
		rsp.Log = "modify file"

		// 生成新文件key
		fileLinkKey := filePrefixKey(m.FileHash)
		fileLinkValue := Int64ToByteArray(app.state.Height+1)

		// 生成旧文件key
		oldFileLinkKey := filePrefixKey(m.OldFileHash)

		// 生成blcok链表key
		blockLinkKey := blockPrefixKey(app.state.Height+1)		
		blockLinkValue := FindKey(db, oldFileLinkKey)

		// 添加到 db
		AddKV(db, fileLinkKey, fileLinkValue) 
		AddKV(db, blockLinkKey, blockLinkValue) 

		// 生成用户key, 生成file_data
		userFileKey := userPrefixKey(m.UserId, m.FileHash)
		NewFileData(db, userFileKey, "", false)

		// 修改旧文件file_data
		oldUserFileKey := userPrefixKey(m.UserId, m.OldFileHash)
		ModifyFileData(db, oldUserFileKey, true)


	//case 0x04: // 删除文件
	//	rsp.Log = "remove file"
	default:
		rsp.Log = "weird command"
		rsp.Code = 2
	}

	app.logger.Info("DeliverTx()", "action", rsp.Log)

	return
}

