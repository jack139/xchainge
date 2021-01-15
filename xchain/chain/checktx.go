package chain

/*
	交易检查
*/



import (
	"xchainge/types"

	"fmt"
	tmtypes "github.com/tendermint/tendermint/abci/types"
)


// 检查交易
func (app *App) isValidDeal(deal *types.Deal) error {
	m := *deal

	// 检查参数
	if len(m.ExchangeID)==0 || len(m.AssetsID)==0 || len(m.Data)==0 { 
		return fmt.Errorf("bad parameters") // 参数问题
	}

	switch {
	case m.Action>0 && m.Action<=3: // 交易链
		{}
	case m.Action>10 && m.Action<=20: // 为测试链保留
		{}
	default:
		return fmt.Errorf("weird deal command")
	}

	return nil
}

// 检查授权
func (app *App) isValidAuth(auth *types.Auth) error {
	m := *auth

	// 检查参数
	if len(m.FromExchangeID)==0 || len(m.DealID)==0 || len(m.ToExchangeID)==0 { 
		return fmt.Errorf("bad parameters") // 参数问题
	}

	switch m.Action {
	case 0x04, 0x05:
		// 检查 dealID 是否存在
		fromExchangeIdStr, _ := cdc.MarshalJSON(m.FromExchangeID)
		respTx := queryTx(app, fromExchangeIdStr, []byte(m.DealID.String()))
		if respTx==nil {
			return fmt.Errorf("DealID not exist.")	
		}
	default:
		return fmt.Errorf("weird auth command")
	}

	return nil
}


// 检查CR
func (app *App) isValidCredit(credit *types.Credit) error {
	m := *credit

	// 检查参数
	if len(m.UserID)==0 { 
		return fmt.Errorf("bad parameters") // 参数问题
	}

	switch m.Action {
	case 0x01:
		{}
	case 0x02, 0x03:
		if len(m.Data)==0 { 
			return fmt.Errorf("bad parameters in Data") // 参数问题
		}
	default:
		return fmt.Errorf("weird auth command")
	}

	return nil
}

/*
	检查交易
*/
func (app *App) CheckTx(req tmtypes.RequestCheckTx) (rsp tmtypes.ResponseCheckTx) {
	app.logger.Info("CheckTx()", "para", req.Tx)

	var err error
	var tx types.Transx

	err = cdc.UnmarshalJSON(req.Tx, &tx)
	if err != nil {
		rsp.Code = 1
		rsp.Log = "error occured in decoding when CheckTx"
		return
	}
	if !tx.Verify() {
		rsp.Code = 2
		rsp.Log = "CheckTx failed"
		return
	}

	// 业务校验
	deal, ok := tx.Payload.(*types.Deal)
	if ok {
		err = app.isValidDeal(deal) // 检查 交易 合法性
	} else {
		auth, ok := tx.Payload.(*types.Auth)
		if ok {
			err = app.isValidAuth(auth)  // 检查 授权 合法性
		} else {
			credit, ok := tx.Payload.(*types.Credit)
			if ok {
				err = app.isValidCredit(credit)  // 检查 CR 合法性
			} else {
				rsp.Code = 3
				rsp.Log = "CheckTx unknown type"
				app.logger.Info("CheckTx() fail", "unknown type", rsp.Log)
				return
			}
		}
	}

	if err!=nil {
		rsp.Log = err.Error()
		rsp.Code = 4
		app.logger.Info("CheckTx() fail", "reason", rsp.Log)
	} else {
		rsp.GasWanted = 1
	}

	return 
}


