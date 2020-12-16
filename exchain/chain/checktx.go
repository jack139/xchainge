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
		ExchangerId string
		AssetsId    string
		Category    string
		DataHash    string
		UserId      string
		Action      byte  
	}
*/

/*
新建
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"exid\":\"1234\",\"aid\":\"abc\",\"dhash\":\"abc\",\"act\":1}"'
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
	if len(m.ExchangerId)==0 || len(m.AssetsId)==0 || len(m.DataHash)==0 || m.Action==0 { 
		return fmt.Errorf("bad parameters") // 参数问题
	}

	switch m.Action {
	case 0x01, 0x02, 0x03, 0x04, 0x05: 
		{}
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


