package http

import (
	"strconv"
	"log"
	"encoding/json"
	"github.com/valyala/fasthttp"
)


/* 提交交易 */
func deal(ctx *fasthttp.RequestCtx) {
	log.Println("deal")

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
	var action int
	_, ok = (*reqData)["action"].(float64)
	if !ok {
		respError(ctx, 9001, "need action")
		return
	} else {
		action = int((*reqData)["action"].(float64))	// 返回整数
	}
	assetsId, ok := (*reqData)["assets_id"].(string)
	if !ok {
		respError(ctx, 9002, "need assets_id")
		return
	}
	data, ok := (*reqData)["data"].(string)
	if !ok {
		respError(ctx, 9003, "need data")
		return
	}
	refer, _ := (*reqData)["refer"].(string)

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	// 提交交易
	respBytes, err := me.Deal(strconv.Itoa(int(action)), assetsId, data, refer)
	if err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respData map[string]interface{}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		respError(ctx, 9005, err.Error())
		return
	}

	resp := map[string] interface{} {
		"data" : respData,
	}

	respJson(ctx, &resp)

}
