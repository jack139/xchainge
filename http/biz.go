package http

import (
	"xchainge/client"

	"strconv"
	"log"
	"encoding/json"
	"encoding/base64"
	"github.com/valyala/fasthttp"
)

/* 企业链业务处理 */

/* 用户注册 
	action == 13
*/
func bizRegister(ctx *fasthttp.RequestCtx) {
	log.Println("biz_register")

	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	userName, ok := (*reqData)["user_name"].(string)
	if !ok {
		respError(ctx, 9001, "need user_name")
		return
	}
	userType, ok := (*reqData)["user_type"].(string)
	if !ok {
		respError(ctx, 9002, "need user_type")
		return
	}
	referrer, _ := (*reqData)["referrer"].(string)

	// 生成新用户密钥
	path := userKeyfilePath+"/"+userName
	me, err := client.GenUserKey(path)
	if err!=nil {
		respError(ctx, 9006, "fail to generate key")
		return
	}

	// 新密钥加入用户缓存
	pubkey := base64.StdEncoding.EncodeToString(me.CryptoPair.PubKey[:])
	SECRET_KEY[pubkey] = me // 保存用户信息

	// 准备数据
	var loadData = map[string]interface{}{
		"user_name" : userName,
		"user_type" : userType,
		"referrer"  : referrer,
	}
	loadBytes, err := json.Marshal(loadData)
	if err != nil {
		respError(ctx, 9008, err.Error())
		return
	}

	// 提交交易
	// data --> user_name,  refer --> user_type
	// user_type: "office" 事务所；"supplier" 供应商；"buyer" 企业用户
	respBytes, err := me.Deal("13", "*", string(loadBytes), "") 
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

	// 返回两个区块id
	resp := map[string] interface{} {
		"block" : respData,
		"userkey" : pubkey,
	}

	respJson(ctx, &resp)
}



/* 签合同 */
func bizContract(ctx *fasthttp.RequestCtx) {
	log.Println("biz_contract")
	doContractDelivery(ctx, 11)
}


/* 验收 */
func bizDelivery(ctx *fasthttp.RequestCtx) {
	log.Println("biz_delivery")
	doContractDelivery(ctx, 12)
}


/*  目前签合同和验收的操作一样，用同一个实现
	action： 11 前合同  12 验收 
*/
func doContractDelivery(ctx *fasthttp.RequestCtx, action int) {
	// POST 的数据
	content := ctx.PostBody()

	// 验签
	reqData, err := checkSign(content)
	if err!=nil {
		respError(ctx, 9000, err.Error())
		return
	}

	// 检查参数
	pubkeyA, ok := (*reqData)["userkey_a"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey_a")
		return
	}
	pubkeyB, ok := (*reqData)["userkey_b"].(string)
	if !ok {
		respError(ctx, 9009, "need userkey_b")
		return
	}
	assetsId, ok := (*reqData)["assets_id"].(string)
	if !ok {
		respError(ctx, 9001, "need assets_id")
		return
	}
	data, ok := (*reqData)["data"].(string)
	if !ok {
		respError(ctx, 9002, "need data")
		return
	}

	// 获取用户密钥
	meA, ok := SECRET_KEY[pubkeyA]
	if !ok {
		respError(ctx, 9011, "wrong userkey_a")
		return
	}

	// 获取用户密钥
	meB, ok := SECRET_KEY[pubkeyB]
	if !ok {
		respError(ctx, 9011, "wrong userkey_b")
		return
	}

	// 准备数据
	var loadData = map[string]interface{}{
		"image" : data,
	}
	loadBytes, err := json.Marshal(loadData)
	if err != nil {
		respError(ctx, 9008, err.Error())
		return
	}

	// 提交交易, A B 两个用户都提交
	respBytesA, err := meA.Deal(strconv.Itoa(action), assetsId, string(loadBytes), "") 
	if err != nil {
		respError(ctx, 9004, err.Error())
		return
	}
	respBytesB, err := meB.Deal(strconv.Itoa(action), assetsId, string(loadBytes), "") 
	if err != nil {
		respError(ctx, 9004, err.Error())
		return
	}

	// 转换成map, 生成返回数据
	var respDataA map[string]interface{}
	var respDataB map[string]interface{}

	if err := json.Unmarshal(respBytesA, &respDataA); err != nil {
		respError(ctx, 9005, err.Error())
		return
	}
	if err := json.Unmarshal(respBytesB, &respDataB); err != nil {
		respError(ctx, 9005, err.Error())
		return
	}

	// 返回两个区块id
	resp := map[string] interface{} {
		"block_a" : respDataA,
		"block_b" : respDataB,
	}

	respJson(ctx, &resp)

}
