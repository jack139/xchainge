package http

import (
	"xchainge/ipfs"

	"strings"
	"log"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

/* data字段是已序列化的json串，反序列化一下 */
func unmarshalData(respData *[]map[string]interface{}) error {
	// 处理data字段
	for _, item := range *respData {
		_, ok := item["data"]
		if !ok {
			continue
		}
		if !strings.HasPrefix(item["data"].(string), "{") {
			continue
		}
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(item["data"].(string)), &data); err != nil {
			return err
		}
		
		// 处理image 字段，从ipfs读取
		_, ok = data["image"]
		if ok && len(data["image"].(string))>0 {
			image_data, err := ipfs.Get(data["image"].(string))
			if err!=nil {
				return err
			}
			data["image"] = string(image_data)
		}		

		item["data"] = data
	}
	return nil
}

/* 查询交易， 只允许查询自己的 */
func queryDeals(ctx *fasthttp.RequestCtx) {
	log.Println("query_deals")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkey, ok := (*reqData)["userkey"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey")
		return
	}

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	// 只查询当前用户的交易
	respBytes, err := me.Query("deal", "_")
	if err!=nil {
		respError(ctx, 9001, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData []map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 处理data字段
	err = unmarshalData(&respData)
	if err!=nil{
		respError(ctx, 9014, err.Error())
		return		
	}

	// 返回结果
	resp := map[string] interface{} {
		"deals" : respData,
	}

	respJson(ctx, &resp)
}


/* 查询授权， 只允许查询自己的 */
func queryAuths(ctx *fasthttp.RequestCtx) {
	log.Println("query_auths")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkey, ok := (*reqData)["userkey"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey")
		return
	}

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	// 只查询当前用户的交易
	respBytes, err := me.Query("auth", "_")
	if err!=nil {
		respError(ctx, 9001, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData []map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	resp := map[string] interface{} {
		"auths" : respData,
	}

	respJson(ctx, &resp)
}


/* 按资产id查询交易 */
func queryByAsstes(ctx *fasthttp.RequestCtx) {
	log.Println("query_by_assets")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkey, ok := (*reqData)["userkey"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey")
		return
	}
	assetsId, ok := (*reqData)["assets_id"].(string)
	if !ok {
		respError(ctx, 9001, "need assets_id")
		return
	}

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	// 只查询当前用户的交易
	respBytes, err := me.Query("assets", assetsId)
	if err!=nil {
		respError(ctx, 9002, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData []map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 处理data字段
	err = unmarshalData(&respData)
	if err!=nil{
		respError(ctx, 9014, err.Error())
		return		
	}

	resp := map[string] interface{} {
		"deals" : respData,
	}

	respJson(ctx, &resp)
}


/* 按参考值查询交易 */
func queryByRefer(ctx *fasthttp.RequestCtx) {
	log.Println("query_by_refer")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkey, ok := (*reqData)["userkey"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey")
		return
	}
	refer, ok := (*reqData)["refer"].(string)
	if !ok {
		respError(ctx, 9001, "need refer")
		return
	}

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	// 只查询当前用户的交易
	respBytes, err := me.Query("refer", refer)
	if err!=nil {
		respError(ctx, 9002, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData []map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 处理data字段
	err = unmarshalData(&respData)
	if err!=nil{
		respError(ctx, 9014, err.Error())
		return		
	}

	resp := map[string] interface{} {
		"deals" : respData,
	}

	respJson(ctx, &resp)
}

/* 指定区块查询交易 */
func queryBlock(ctx *fasthttp.RequestCtx) {
	log.Println("query_block")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkey, ok := (*reqData)["userkey"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey")
		return
	}
	blockId, ok := (*reqData)["block_id"].(string)
	if !ok {
		respError(ctx, 9002, "need block_id")
		return
	}

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	respBytes, err := me.QueryTx(pubkey, blockId)
	if err!=nil {
		respError(ctx, 9003, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData map[string]interface{}
	if len(respBytes)>0 {
		if err := json.Unmarshal(respBytes, &respData); err != nil {
			respError(ctx, 9004, err.Error())
			return
		}
	}


	// 处理data字段
	temp := []map[string]interface{}{ respData }
	err = unmarshalData(&temp)
	if err!=nil{
		respError(ctx, 9014, err.Error())
		return		
	}

	resp := map[string] interface{} {
		"blcok" : respData,
	}

	respJson(ctx, &resp)
}


/* 指定区块查询交易 */
func queryRawBlock(ctx *fasthttp.RequestCtx) {
	log.Println("query_raw_block")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkey, ok := (*reqData)["userkey"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey")
		return
	}
	blockId, ok := (*reqData)["block_id"].(string)
	if !ok {
		respError(ctx, 9002, "need block_id")
		return
	}

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	respBytes, err := me.QueryRawBlock(pubkey, blockId)
	if err!=nil {
		respError(ctx, 9003, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData map[string]interface{}
	if len(respBytes)>0 {
		if err := json.Unmarshal(respBytes, &respData); err != nil {
			respError(ctx, 9004, err.Error())
			return
		}
	}
	resp := map[string] interface{} {
		"blcok" : respData,
	}

	respJson(ctx, &resp)
}
