package chain

/*
	链上查询
*/



import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/abci/types"
)

/*
查询资产历史
curl -g 'http://localhost:26657/abci_query?data="{\"query\":\"acbadbdc\",\"act\":1}"'

查询交易所历史
curl -g 'http://localhost:26657/abci_query?data="{\"query\":\"1234\",\"act\":2}"'

测试入口
curl -g 'http://localhost:26657/abci_query?data="{\"act\":255}"'
*/
func (app *App) Query(req types.RequestQuery) (rsp types.ResponseQuery) {
	app.logger.Info("Query()", "para", req.Data)

	db := app.state.db

	var m QueryReq

	err := json.Unmarshal(req.Data, &m)
	if err != nil {
		rsp.Log = "bad json format"
		rsp.Code = 1
		return
	}

	switch m.Action {
	case 0x01: // 资产交易历史
		rsp.Log = "assets history"

		var respHistory []RespAssetsHistory

		// 文件key, 找到链头
		assetsLinkKey := assetsPrefixKey(m.Query)
		height := FindKey(db, assetsLinkKey)  // 这里 height 返回是 []byte
		for ;len(height)!=0; {
			// 高度转换为int64
			heightInt := ByteArrayToInt64(height)
			// 获取区块内容
			block := GetBlock(heightInt)

			// 交易请求转为struct
			var txReq TxReq
			err := json.Unmarshal(block.Data.Txs[0], &txReq)
			if err != nil {
				panic(err)
			}

			// 添加到返回结果数组
			respHistory = append(respHistory, RespAssetsHistory{
				TxRequest: txReq,
				BlockTime: block.Header.Time,
			})

			fmt.Println(heightInt, string(block.Data.Txs[0]))

			// 在blcok链上找下一个
			blockLinkKey := blockPrefixKey(heightInt)
			height = FindKey(db, blockLinkKey)
		}

		// 返回结果转为json
		respBytes, err := json.Marshal(respHistory)
		if err != nil {
			panic(err)
		}
		rsp.Value = respBytes

	case 0x02: // 查询交易所的交易
		rsp.Log = "exchanger history"

		var respHistory []RespAssetsHistory

		high := app.state.Height // 链高度

		/* 遍历整个链 */
		for i:=int64(1);i<=high;i++ {
			// 获取区块内容
			block := GetBlock(i)			

			if len(block.Data.Txs)==0 {  // 忽略空块
				continue
			}

			// 交易请求转为struct
			var txReq TxReq
			err := json.Unmarshal(block.Data.Txs[0], &txReq)
			if err != nil {
				panic(err)
			}

			if txReq.ExchangerId == m.Query {  // 是否相同交易所？
				// 添加到返回结果数组
				respHistory = append(respHistory, RespAssetsHistory{
					TxRequest: txReq,
					BlockTime: block.Header.Time,
				})

				fmt.Println(i, string(block.Data.Txs[0]))				
			}

		}

		// 返回结果转为json
		respBytes, err := json.Marshal(respHistory)
		if err != nil {
			panic(err)
		}
		rsp.Value = respBytes


	case 0xff: // 测试
		rsp.Log = "test"

	default:
		rsp.Log = "weird command"
		rsp.Code = 2
	}

	return
}

