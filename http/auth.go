package http

import (
	"log"
	"encoding/json"
	"github.com/valyala/fasthttp"
)


/* 请求授权 */
func authRequest(ctx *fasthttp.RequestCtx) {
	log.Println("auth_request")

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
	dealId, ok := (*reqData)["deal_id"].(string)
	if !ok {
		respError(ctx, 9001, "need deal_id")
		return
	}
	fromExchangeId, ok := (*reqData)["from_exchange_id"].(string)
	if !ok {
		respError(ctx, 9002, "need from_exchange_id")
		return
	}

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	// 提交 授权请求
	respBytes, err := me.AuthRequest(fromExchangeId, dealId)
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


/* 请求授权 */
func authResponse(ctx *fasthttp.RequestCtx) {
	log.Println("auth_response")

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
	authId, ok := (*reqData)["auth_id"].(string)
	if !ok {
		respError(ctx, 9001, "need auth_id")
		return
	}

	// 获取用户密钥
	me, ok := SECRET_KEY[pubkey]
	if !ok {
		respError(ctx, 9011, "wrong userkey")
		return
	}

	// 提交 授权响应
	respBytes, err := me.AuthResponse(authId)
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
