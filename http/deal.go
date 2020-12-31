package http

import (
	"strconv"
	"log"
	"github.com/valyala/fasthttp"
)


/* 提交交易 */
func deal(ctx *fasthttp.RequestCtx) {
	log.Println("deal")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, me, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	var action int
	_, ok := (*reqData)["action"].(float64)
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
	refer, _ := (*reqData)["data"].(string)

	// 提交交易
	err = me.Deal(strconv.Itoa(int(action)), assetsId, data, refer)
	if err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 正常 返回空
	resp := map[string] interface{} {
		"data" : nil,
	}
	respJson(ctx, &resp)
}
