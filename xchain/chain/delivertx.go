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
	提交区块
*/
func (app *App) DeliverTx(req tmtypes.RequestDeliverTx) (rsp tmtypes.ResponseDeliverTx) {
	app.logger.Info("DeliverTx()", "para", req.Tx)

	db := app.state.db
	raw := req.Tx
	fmt.Println(string(raw))

	var tx types.Transx
	cdc.UnmarshalJSON(raw, &tx) //由于之前CheckTx中转换过，所以这里按道理不会有error

	// 数据上链
	deal, ok := tx.Payload.(*types.Deal)	// 交易
	if ok {
		switch deal.Action {
		case 0x01, 0x02, 0x03:
			rsp.Log = actionMessage[deal.Action]

			var exchangeID []byte
			copy(exchangeID[:], deal.ExchangeID[:])

			// 完善链表
			AddToLink(db, "exchange", exchangeID, app.state.Height+1)
			AddToLink(db, "assets", deal.AssetsID, app.state.Height+1)

		default:
			rsp.Log = "weird command"
			rsp.Code = 2
		}
	} else {
		auth, ok := tx.Payload.(*types.Auth)	// 授权
		if ok {
			switch auth.Action {
			case 0x04, 0x05, 0x06:
				rsp.Log = actionMessage[auth.Action]

				var exchangeID []byte
				copy(exchangeID[:], auth.FromExchangeID[:])

				// 完善链表
				AddToLink(db, "exchange", exchangeID, app.state.Height+1)
				AddToLink(db, "assets", auth.AssetsID, app.state.Height+1)

			default:
				rsp.Log = "weird command"
				rsp.Code = 2
			}
		}
	}

	app.logger.Info("DeliverTx()", "action", rsp.Log)

	return
}

