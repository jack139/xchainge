package client

import (
	"xchainge/types"

	"bytes"
	"fmt"
	"encoding/json"
)


// 链上查询指定 ID 的交易数据
// xcli queryBlock dcfe656c-6c65-45e7-9e94-f082a068a93d
func (me *User) QueryTx(exchangeId, idStr string) ([]byte, error) {
	addr, _ := cdc.MarshalJSON(*me.CryptoPair.PubKey)

	tx, err := queryTx(addr, exchangeId, idStr)
	if err != nil {
		return nil, err
	}

	if tx==nil {  // 未找到
		return nil, nil
	}

	// 转换为返回的结构
	respQ := txToResp(me, tx)

	// 返回结果转为json
	respBytes, err := json.Marshal(*respQ)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("json => %s\n", respBytes)
	return respBytes, nil
}

func queryTx(addr []byte, exchangeId, idStr string) (*types.Transx, error) {
	var buf bytes.Buffer
	buf.WriteString("/")
	buf.Write(addr)
	buf.WriteString("/query/tx")
	//获得拼接后的字符串
	path := buf.String()
	if exchangeId!="_" {  // 用户公钥需要加双引号
		exchangeId = "\"" + exchangeId + "\""	
	}

	// req.Data 格式： ["用户公钥", "DealID"]
	reqBytes, _ := cdc.MarshalJSON([2][]byte{[]byte(exchangeId), []byte(idStr)})

	rsp, err := cli.ABCIQuery(ctx, path, reqBytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data := rsp.Response.Value
	//fmt.Printf("resp => %s\n", data)

	if data!=nil {
		var respTx types.Transx
		cdc.UnmarshalJSON(data, &respTx)
		return &respTx, nil
	} else {
		return nil, nil
	}
}


// 链上查询  category取值： deal, auth, assets, refer
// deal 和 auth 可以带公钥，查其他人的 
// xcli queryDeal _
// xcli queryDeal j9cIgmm17x0aLApf0i20UR7Pj34Ua/JwyWOuBGgYIFg=
func (me *User) Query(category, queryContent string) ([]byte, error) {
	addr, _ := cdc.MarshalJSON(*me.CryptoPair.PubKey)

	var respList []types.RespQuery
	txList, err := query(addr, category, queryContent)

	for _, tx := range *txList {
		respQ := txToResp(me, &tx)
		respList = append(respList, *respQ)
	}

	// 返回结果转为json
	respBytes, err := json.Marshal(respList)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("json => %s\n", respBytes)

	return respBytes, nil
}

func query(addr []byte, category, queryContent string) (*[]types.Transx, error) {
	var buf bytes.Buffer
	buf.WriteString("/")
	buf.Write(addr)
	buf.WriteString("/query/")
	buf.WriteString(category)
	//获得拼接后的字符串
	path := buf.String()
	if (category=="deal"||category=="auth") && queryContent!="_" {  
		if queryContent[0]!='"' { // 用户公钥需要加双引号
			queryContent = "\"" + queryContent + "\""
		}
	}
	rsp, err := cli.ABCIQuery(ctx, path, []byte(queryContent))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data := rsp.Response.Value
	//fmt.Printf("resp => %s\n", data)

	/*
		exchange 不解密
		assets 根据授权解密
		refer 不解密
	*/

	var txHistory, txResp []types.Transx
	cdc.UnmarshalJSON(data, &txHistory)

	for _, tx := range txHistory {
		if category=="auth" {
			_, ok := tx.Payload.(*types.Auth)	// 授权
			if ok {
				txResp = append(txResp, tx)
			}
		} else { // category == deal, assets, refer
			_, ok := tx.Payload.(*types.Deal) // 交易
			if ok {
				txResp = append(txResp, tx)
			}
		} 
	}

	return &txResp, nil

}


// 检查 授权请求的交易（dealID） 是否已进行响应
func checkAuthResp(addr []byte, toExchangeId, dealId string) (bool, error) {
	var buf bytes.Buffer
	buf.WriteString("/")
	buf.Write(addr)
	buf.WriteString("/query/check_auth_resp")
	//获得拼接后的字符串
	path := buf.String()

	// req.Data 格式： ["用户公钥", "DealID"]
	reqBytes, _ := cdc.MarshalJSON([2][]byte{[]byte(toExchangeId), []byte(dealId)})

	rsp, err := cli.ABCIQuery(ctx, path, reqBytes)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return rsp.Response.Value[0]==1, nil
}
