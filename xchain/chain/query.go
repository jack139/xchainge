package chain

/*
	链上查询
*/



import (
	"xchainge/types"

	"encoding/json"
	"fmt"
	"regexp"
	"bytes"
	tmtypes "github.com/tendermint/tendermint/abci/types"
)


const (
	// 匹配如下格式
	// /userpubkey/query/category/content
	queryPathPattern string = `^/((?P<uk>\S+)/query/(?P<cate>\S+)/(?P<content>\S+)?)$`
)

func getMatchMap(submatches []string, groupNames []string) map[string]string {
	result := make(map[string]string)
	for i, name := range groupNames {
		if i != 0 && name != "" {
			result[name] = submatches[i]
		}
	}
	return result
}

/*
查询资产历史
curl -g 'http://localhost:26657/abci_query?data="{\"query\":\"abc\",\"act\":1}"'

查询交易所历史
curl -g 'http://localhost:26657/abci_query?data="{\"query\":\"1234\",\"act\":2}"'

查询资产历史
curl -g 'http://localhost:26657/abci_query?data="{\"query\":\"xxx\",\"act\":3}"'

测试入口
curl -g 'http://localhost:26657/abci_query?data="{\"act\":255}"'
*/
func (app *App) Query(req tmtypes.RequestQuery) (rsp tmtypes.ResponseQuery) {
	app.logger.Info("Query()", "para", req.Data)

	db := app.state.db

	fmt.Println(req.Path)
	reg := regexp.MustCompile(queryPathPattern)
	submatches := reg.FindStringSubmatch(req.Path)
	groupNames := reg.SubexpNames()
	matchmap := getMatchMap(submatches, groupNames)

	// 解码 exchangeId (公钥)
	var exchangeId []byte
	err := cdc.UnmarshalJSON([]byte(matchmap["uk"]), &exchangeId)
	if err != nil {
		rsp.Code = 1
		rsp.Log = err.Error()
		return
	}

	if matchmap["cate"] == "" {
		rsp.Log = "no category"
		return
	}

	switch matchmap["cate"] {
	case "assets", "exchange": // 资产交易历史， 交易所交易历史
		var respHistory []RespAssetsHistory
		var linkKey []byte
		var linkType string

		// 文件key, 找到链头
		if matchmap["cate"]=="assets" {
			rsp.Log = "assets history"
			linkKey = assetsPrefixKey(req.Data)
			linkType = "assets"
		} else {
			rsp.Log = "exhcange history"
			linkKey = exhcangePrefixKey(req.Data)
			linkType = "exchange"
		}
		height := FindKey(db, linkKey)  // 这里 height 返回是 []byte
		for ;len(height)!=0; {
			// 高度转换为int64
			heightInt := ByteArrayToInt64(height)
			// 获取区块内容
			block := GetBlock(heightInt)

			var tx types.Transx
			cdc.UnmarshalJSON(block.Data.Txs[0], &tx)

			// 添加到返回结果数组
			respHistory = append(respHistory, RespAssetsHistory{
				TxRequest: tx,
				BlockTime: block.Header.Time,
			})

			fmt.Println(heightInt, string(block.Data.Txs[0]))

			// 在blcok链上找下一个
			blockLinkKey := blockPrefixKey(linkType, heightInt)
			height = FindKey(db, blockLinkKey)
		}

		// 返回结果转为json
		respBytes, err := json.Marshal(respHistory)
		if err != nil {
			panic(err)
		}
		rsp.Value = respBytes

	case "refer": // refer参考值的交易 （全库遍历）
		rsp.Log = "refer history"

		var respHistory []RespAssetsHistory

		high := app.state.Height // 链高度

		/* 遍历整个链 */
		for i:=high;i>0;i-- {  //  从顶向下遍历
			// 获取区块内容
			block := GetBlock(i)

			if len(block.Data.Txs)==0 {  // 忽略空块
				continue
			}

			// 交易请求转为struct
			var tx types.Transx
			cdc.UnmarshalJSON(block.Data.Txs[0], &tx)

			var refer []byte
			deal, ok := tx.Payload.(*types.Deal)	// 交易
			if ok {
				refer = deal.Refer
			} else {
				auth, _ := tx.Payload.(*types.Auth)	// 授权
				refer = auth.Refer
			}
			res := bytes.Compare(refer, req.Data)
			if res==0 {  // 是否相同refer
				// 添加到返回结果数组
				respHistory = append(respHistory, RespAssetsHistory{
					TxRequest: tx,
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

	default:
		rsp.Log = "weird command"
		rsp.Code = 2
	}

	return
}

