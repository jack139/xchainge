package http

import (
	"fmt"
	//"log"
	//"time"
	"github.com/valyala/fasthttp"
)


/* 空接口, 只进行签名校验 */
func doNonthing(ctx *fasthttp.RequestCtx) {
	respJson := map[string] interface{} {
		"code" : 0,
		"msg" : "",
		"data" : "",
	}

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	data, err := checkSign(content)
	if err!=nil {
		respJson["code"] = 8000
		respJson["msg"] = err.Error()
		goto goodbye
	}
	fmt.Printf("%v\n", *data)

	respJson["data"] = *data
	respJson["msg"] = "success"

goodbye:
	doJSONWrite(ctx, fasthttp.StatusOK, respJson)
}
