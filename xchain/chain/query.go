package chain

/*
	链上查询
*/

import (
	"xchainge/types"

	"fmt"
	"regexp"
	"bytes"
	tmtypes "github.com/tendermint/tendermint/abci/types"
	tmtypes2 "github.com/tendermint/tendermint/types"
)


const (
	// 匹配如下格式
	// /userpubkey/query/category
	queryPathPattern string = `^/((?P<uk>\S+)/query/(?P<cate>\S+)?)$`
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
 	按 refer 检索		/"<pubkey>"/query/refer
	按 assets 检索	/"<pubkey>"/query/assets
	检索 所有 auth 块		/"<pubkey>"/query/auth
	检索 所有 deal 块		/"<pubkey>"/query/deal
	检索 指定 deal/auth 块	/"<pubkey>"/query/tx
	检索 指定 block raw data	/"<pubkey>"/query/raw
*/


/*
查询资产历史
/"qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs="/query/assets
curl -g 'http://localhost:26657/abci_query?data="123"&path="/\"qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=\"/query/assets"'

查询交易所历史
/"qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs="/query/deal
curl -g 'http://localhost:26657/abci_query?data="_"&path="/\"qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=\"/query/deal"'

查询refer历史
/"qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs="/query/refer
curl -g 'http://localhost:26657/abci_query?data="yyy"&path="/\"qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=\"/query/refer"'
*/
func (app *App) Query(req tmtypes.RequestQuery) (rsp tmtypes.ResponseQuery) {
	app.logger.Info("Query()", "para", req.Data)

	db := app.state.db

	fmt.Println(req.Path)
	reg := regexp.MustCompile(queryPathPattern)
	submatches := reg.FindStringSubmatch(req.Path)
	groupNames := reg.SubexpNames()
	//fmt.Println(submatches, groupNames)
	if len(submatches)!=len(groupNames) {
		rsp.Log = "path error"
		return		
	}
	matchmap := getMatchMap(submatches, groupNames)

	// 解码 exchangeId (公钥)，序列化文本
	exchangeId := []byte(matchmap["uk"])

	if matchmap["cate"] == "" {
		rsp.Log = "no category"
		return
	}

	switch matchmap["cate"] {
	case "assets", "deal", "auth", "refer", "credit": // 资产交易历史， deal交易历史, auth历史
		var respHistory []types.Transx
		var linkKey []byte
		var linkType string
		var qData *[]byte

		if string(req.Data)=="_" {  //  查询自己的交易记录
			qData = &exchangeId
		} else {
			// TODO 检查授权
			qData = &req.Data
		}

		fmt.Printf("--> %s\n", *qData)

		// 文件key, 找到链头
		switch matchmap["cate"] {
		case "assets":
			rsp.Log = "assets history"
			linkKey = assetsPrefixKey(*qData)
			linkType = "assets"
		case "refer":
			rsp.Log = "exhcange history"
			linkKey = exhcangePrefixKey(exchangeId)
			linkType = "exchange"
		default:
			rsp.Log = "exhcange history"
			linkKey = exhcangePrefixKey(*qData)
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

			deal, ok := tx.Payload.(*types.Deal)	// 交易块
			if ok {
				if matchmap["cate"]=="auth" || matchmap["cate"]=="credit" { // deal, assets, refer
					goto go_next
				}
				if matchmap["cate"]=="refer" { 
					fmt.Printf("-----> %s %s\n", deal.Refer, req.Data)
					res := bytes.Compare(deal.Refer, req.Data)
					if res!=0 {  // 不同的refer
						goto go_next
					}
				}
			} else {  // 授权块，没有 refer
				_, ok := tx.Payload.(*types.Credit)	// CR
				if ok {
					if matchmap["cate"]!="credit" { 
						goto go_next
					}					
				} else {
					if matchmap["cate"]!="auth" { // auth
						goto go_next
					}
				}
			}

			respHistory = append(respHistory, tx) // 添加到返回结果数组
			fmt.Printf(">> %d", heightInt)

		go_next:
			// 在blcok链上找下一个
			blockLinkKey := blockPrefixKey(linkType, heightInt)
			height = FindKey(db, blockLinkKey)
		}

		fmt.Println()

		respBytes, _ := cdc.MarshalJSON(respHistory)
		rsp.Value = respBytes

	case "___global_refer": // refer参考值的交易 （全库遍历）
		rsp.Log = "refer history"

		var respHistory []types.Transx

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
			deal, ok := tx.Payload.(*types.Deal)	// 交易块
			if ok {
				refer = deal.Refer
			} else {  // 授权块，没有 refer
				continue 
			}

			res := bytes.Compare(refer, req.Data)
			if res==0 {  // 是否相同refer
				// 添加到返回结果数组
				respHistory = append(respHistory, tx)

				fmt.Printf(">> %d ", i)
			}

		}

		fmt.Println()

		respBytes, _ := cdc.MarshalJSON(respHistory)
		rsp.Value = respBytes

	case "tx": // 指定ID的 deal 或 auth 
		var qData [2][]byte

		// req.Data 格式： ["用户公钥", "DealID"]
		cdc.UnmarshalJSON(req.Data, &qData)

		if string(qData[0])=="_" {  //  查询自己的交易记录
			qData[0] = exchangeId
		}

		fmt.Printf("--> %s %s\n", qData[0], qData[1])

		// 文件key, 找到链头
		rsp.Log = "query TX"

		respTx := queryTx(app, qData[0], qData[1])

		if respTx!=nil {
			respBytes, _ := cdc.MarshalJSON(*respTx)
			rsp.Value = respBytes
		} else {
			rsp.Value = nil
		}

	case "raw": // 制定ID的 block raw data
		var qData [2][]byte

		// req.Data 格式： ["用户公钥", "DealID"]
		cdc.UnmarshalJSON(req.Data, &qData)

		if string(qData[0])=="_" {  //  查询自己的交易记录
			qData[0] = exchangeId
		}

		fmt.Printf("--> %s %s\n", qData[0], qData[1])

		// 文件key, 找到链头
		rsp.Log = "query raw data"

		block := queryRawBlock(app, qData[0], qData[1])

		if block!=nil {
			respBytes, _ := cdc.MarshalJSON(*block)
			rsp.Value = respBytes
		} else {
			rsp.Value = nil
		}

	case "check_auth_resp": // 检查 授权请求（authID） 是否已进行响应
		var qData [2][]byte

		// req.Data 格式： ["用户公钥", "DealID"]
		cdc.UnmarshalJSON(req.Data, &qData)

		if string(qData[0])=="_" {  //  查询自己的交易记录
			qData[0] = exchangeId
		}

		fmt.Printf("--> %s %s\n", qData[0], qData[1])

		rsp.Log = "check auth response"
		linkKey := exhcangePrefixKey(qData[0])

		rsp.Value = []byte{0}

		height := FindKey(db, linkKey)  // 这里 height 返回是 []byte
		for ;len(height)!=0; {
			// 高度转换为int64
			heightInt := ByteArrayToInt64(height)
			// 获取区块内容
			block := GetBlock(heightInt)

			var tx types.Transx
			cdc.UnmarshalJSON(block.Data.Txs[0], &tx)

			auth, ok := tx.Payload.(*types.Auth)	// 授权块
			if ok {
				// DealID相同，说明已回复
				if auth.Action==0x05 && auth.DealID.String()==string(qData[1]) {
					rsp.Value = []byte{1}
					break
				}
			}

			// 在blcok链上找下一个
			blockLinkKey := blockPrefixKey("exchange", heightInt)
			height = FindKey(db, blockLinkKey)
		}

	default:
		rsp.Log = "weird command"
		rsp.Code = 2
	}

	return
}


func queryTx(app *App, exchangeId, dealId []byte) (respTx *types.Transx) {
	db := app.state.db

	// 找到链头
	linkKey := exhcangePrefixKey(exchangeId)
	height := FindKey(db, linkKey)  // 这里 height 返回是 []byte

	for ;len(height)!=0; {
		// 高度转换为int64
		heightInt := ByteArrayToInt64(height)
		// 获取区块内容
		block := GetBlock(heightInt)

		var tx types.Transx
		cdc.UnmarshalJSON(block.Data.Txs[0], &tx)

		deal, ok := tx.Payload.(*types.Deal)	// 交易块
		if ok {
			if deal.ID.String()==string(dealId) {
				respTx = &tx
				return
			}
		} else {  // 授权块，没有 refer
			auth, ok := tx.Payload.(*types.Auth)	// 授权块
			if ok {
				if auth.ID.String()==string(dealId) {
					respTx = &tx
					return
				}
			}
		}

		// 在blcok链上找下一个
		blockLinkKey := blockPrefixKey("exchange", heightInt)
		height = FindKey(db, blockLinkKey)
	}

	return nil
}

func queryRawBlock(app *App, exchangeId, dealId []byte) (block *tmtypes2.Block) {
	db := app.state.db

	// 找到链头
	linkKey := exhcangePrefixKey(exchangeId)
	height := FindKey(db, linkKey)  // 这里 height 返回是 []byte

	for ;len(height)!=0; {
		// 高度转换为int64
		heightInt := ByteArrayToInt64(height)
		// 获取区块内容
		block = GetBlock(heightInt)

		var tx types.Transx
		cdc.UnmarshalJSON(block.Data.Txs[0], &tx)

		deal, ok := tx.Payload.(*types.Deal)	// 交易块
		if ok {
			if deal.ID.String()==string(dealId) {
				return
			}
		} else {  // 授权块，没有 refer
			auth, ok := tx.Payload.(*types.Auth)	// 授权块
			if ok {
				if auth.ID.String()==string(dealId) {
					return
				}
			}
		}

		// 在blcok链上找下一个
		blockLinkKey := blockPrefixKey("exchange", heightInt)
		height = FindKey(db, blockLinkKey)
	}

	return nil
}

