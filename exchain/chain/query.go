package chain

/*
	链上查询
*/



import (
	"encoding/json"
	"bytes"
	"fmt"

	"github.com/tendermint/tendermint/abci/types"
)

/*
	type QueryReq struct {
		UserId      string `json:"user_id"` // 文件主的用户id，action==2时提供
		FileHash    string `json:"file_hash"` // 文件hash，action==1时提供
		Action      byte   `json:"action"` // 0x01 查询文件历史, 0x02 查询用户的文件列表
	}
*/

/*
查询文件历史
curl -g 'http://localhost:26657/abci_query?data="{\"fhash\":\"5678\",\"act\":1}"'

查询用户文件
curl -g 'http://localhost:26657/abci_query?data="{\"uid\":\"abc\",\"act\":2}"'

测试入口
curl -g 'http://localhost:26657/abci_query?data="{\"act\":3}"'
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
	case 0x01: // 文件历史
		rsp.Log = "file history"

		var respHistory []RespFileHistory

		// 文件key, 找到链头
		fileLinkKey := filePrefixKey(m.FileHash)
		height := FindKey(db, fileLinkKey)  // 这里 height 返回是 []byte
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
			respHistory = append(respHistory, RespFileHistory{
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


	case 0x02: // 用户文件
		rsp.Log = "user file list"

		start := fmt.Sprintf("%s%s:", userFilePrefixKey, m.UserId)
		end := fmt.Sprintf("%s\xff", start)

		// 循环获取
		itr, err := db.Iterator([]byte(start), []byte(end))
		if err != nil {
			panic(err)
		}

		var respUserFile []RespUserFile
		for ; itr.Valid(); itr.Next() {
			// 从key分解出用户id和文件hash
			parts := bytes.Split(itr.Key(), []byte(":"))
			if len(parts) != 3 {
				panic("bad key format")
			}

			// 交易请求转为struct
			var fileData FileData
			err := json.Unmarshal(itr.Value(), &fileData)
			if err != nil {
				panic(err)
			}
			//fmt.Println(fileData)

			// 添加到返回结果数组
			respUserFile = append(respUserFile, RespUserFile{
				UserId: string(parts[1]),
				FileName: fileData.FileName,
				FileHash: string(parts[2]),
				Modified: fileData.Modified,
			})

			fmt.Println(string(itr.Key()), "=", string(itr.Value()))
		}

		// 返回结果转为json
		respBytes, err := json.Marshal(respUserFile)
		if err != nil {
			panic(err)
		}
		rsp.Value = respBytes


	case 0x03: // 测试
		rsp.Log = "test"

	default:
		rsp.Log = "weird command"
		rsp.Code = 2
	}

	return
}

