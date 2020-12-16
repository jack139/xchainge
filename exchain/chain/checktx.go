package chain

/*
	交易检查
*/



import (
	"encoding/json"
	"fmt"

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

/*
新建
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"fhash\":\"1234\",\"uid\":\"abc\",\"act\":1}"'

浏览
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"fhash\":\"1234\",\"uid\":\"abc\",\"act\":2,\"rid\":\"xyz\",\"nonce\":123}"'

修改
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"fhash\":\"5678\",\"uid\":\"abc\",\"act\":3,\"ofhash\":\"1234\"}"'
*/

// 检查参数
func (app *App) isValid(tx []byte) error {
	db := app.state.db
	var m TxReq

	err := json.Unmarshal(tx, &m)
	if err != nil {
		return err // json 格式问题
	}

	// 检查参数
	if len(m.UserId)==0 || len(m.FileHash)==0 || m.Action==0 { 
		return fmt.Errorf("bad parameters") // 参数问题
	}

	switch m.Action {
	case 0x01: // 新建文件
		if FileHashExisted(db, m.FileHash) {
			return fmt.Errorf("file_hash existed")
		}
	case 0x02: // 浏览文件
		if len(m.ReaderId)==0 {
			return fmt.Errorf("reader_id needed")
		}
		if !FileHashExisted(db, m.FileHash) {
			return fmt.Errorf("file_hash not existed")
		}
	case 0x03: // 修改文件
		if len(m.OldFileHash)==0 {
			return fmt.Errorf("old_file_hash needed")
		}
		if !FileHashExisted(db, m.OldFileHash) {
			return fmt.Errorf("old_file_hash not existed")
		}
		if FileHashExisted(db, m.FileHash) {
			return fmt.Errorf("new file_hash existed")
		}
	//case 0x04: // 删除文件
	//	rsp.Log = "remove file"
	default:
		return fmt.Errorf("weird command")
	}

	return nil
}

func (app *App) CheckTx(req types.RequestCheckTx) (rsp types.ResponseCheckTx) {
	app.logger.Info("CheckTx()", "para", req.Tx)

	err := app.isValid(req.Tx)
	if err!=nil {
		rsp.Log = err.Error()
		rsp.Code = 1
		app.logger.Info("CheckTx() fail", "reason", rsp.Log)
	} else {
		rsp.GasWanted = 1
	}

	return 
}


