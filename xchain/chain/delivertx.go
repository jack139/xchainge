package chain

/*
	交易上链处理
*/


import (
	"xchainge/types"

	"fmt"
	tmtypes "github.com/tendermint/tendermint/abci/types"
)


/*
	提交区块，主要业务逻辑放这里（没做实际事情， 其实检查已经在 checkTx 中做了）
*/
func (app *App) DeliverTx(req tmtypes.RequestDeliverTx) (rsp tmtypes.ResponseDeliverTx) {
	app.logger.Info("DeliverTx()", "para", req.Tx)

	var tx types.Transx
	cdc.UnmarshalJSON(req.Tx, &tx) //由于之前CheckTx中转换过，所以这里按道理不会有error

	// 数据上链
	deal, ok := tx.Payload.(*types.Deal)	// 交易
	if ok {
		switch {
		case deal.Action>0 && deal.Action<=3: // 交易链
			rsp.Log = actionMessage[deal.Action]
			// 业务逻辑放这里
		case deal.Action>10 && deal.Action<=20: // 为测试链保留
			rsp.Log = fmt.Sprintf("deal action=%d", deal.Action)
			// 业务逻辑放这里

		default:
			rsp.Log = "weird deal command"
			rsp.Code = 2
		}
	} else {
		auth, ok := tx.Payload.(*types.Auth)	// 授权
		if ok {
			switch auth.Action {
			case 0x04, 0x05:
				rsp.Log = actionMessage[auth.Action]
				// 业务逻辑放这里

			default:
				rsp.Log = "weird auth command"
				rsp.Code = 3
			}
		}
	}

	app.logger.Info("DeliverTx()", "action", rsp.Log)

	return
}

/*
	结束区块生成，此处更新链表
*/
func (app *App) EndBlock(req tmtypes.RequestEndBlock) (rsp tmtypes.ResponseEndBlock) {
	app.logger.Info("EndBlock()", "height", req.Height)

	db := app.state.db
	block := GetBlock(req.Height)

	if len(block.Data.Txs)==0 {
		return
	}

	var tx types.Transx
	cdc.UnmarshalJSON(block.Data.Txs[0], &tx)

	// 更新链表， 放这里，确保height的数据准确
	deal, ok := tx.Payload.(*types.Deal)	// 交易
	if ok {
		exchangeID, _ := cdc.MarshalJSON(deal.ExchangeID)

		// 完善链表
		AddToLink(db, "exchange", exchangeID, req.Height)
		AddToLink(db, "assets", deal.AssetsID, req.Height)
	} else {
		auth, ok := tx.Payload.(*types.Auth)	// 授权
		if ok {
			// 完善链表
			exchangeID, _ := cdc.MarshalJSON(auth.FromExchangeID)

			if auth.Action==0x05 {  // 授权时， 加入被授权人的链表
				exchangeID, _ = cdc.MarshalJSON(auth.ToExchangeID)
			}

			AddToLink(db, "exchange", exchangeID, req.Height)
		}
	}

	return
}
