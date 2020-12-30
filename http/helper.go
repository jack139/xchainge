package http

import (
	"xchainge/client"

	"fmt"
	"sort"
	"strconv"
	"log"
	"time"
	"encoding/json"
	"encoding/base64"
	"crypto/sha256"
	"github.com/valyala/fasthttp"
)

const (
	keyFileRoot = "./"
)

/* appid : 密钥文件路径， secret 为 密钥文件签名 */
var SECRET_KEY = map[string]string{
	"19E179E5DC29C05E65B90CDE57A1C7E5" : "user1",
	"66A095861BAE55F8735199DBC45D3E8E" : "user2",
	"75C50F018B34AC0240915EC685F5961B" : "user3",
	"3EA25569454745D01219080B779F021F" : "user4",
}

/* 返回值的 content-type */
var (
	strContentType = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

/* 处理返回值，返回json */
func respJson(ctx *fasthttp.RequestCtx, data *map[string] interface{}) {
	respJson := map[string] interface{} {
		"code" : 0,
		"msg"  : "success",
		"data" : *data,
	}
	doJSONWrite(ctx, fasthttp.StatusOK, respJson)
}

func respError(ctx *fasthttp.RequestCtx, code int, msg string) {
	log.Println("Error: ", msg)
	respJson := map[string] interface{} {
		"code" : code,
		"msg"  : msg,
		"data" : "",
	}
	doJSONWrite(ctx, fasthttp.StatusOK, respJson)
}

func doJSONWrite(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	start := time.Now()
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		elapsed := time.Since(start)
		log.Printf("", elapsed, err.Error(), obj)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}



/*
	接口验签，返回data数据
*/
func checkSign(content []byte) (*map[string]interface{}, *client.User, error) {
	fields := make(map[string]interface{})
	if err := json.Unmarshal(content, &fields); err != nil {
		return nil, nil, err
	}

	var appId, version, signType, signData string
	var timestamp int64
	var data map[string]interface{}
	var ok bool

	// 检查参数
	if appId, ok = fields["appid"].(string); !ok {
		return nil, nil, fmt.Errorf("need appId")
	}	
	if version, ok = fields["version"].(string); !ok {
		return nil, nil, fmt.Errorf("need version")
	}	
	if signType, ok = fields["sign_type"].(string); !ok {
		return nil, nil, fmt.Errorf("need signType")
	}	
	if signData, ok = fields["sign_data"].(string); !ok {
		return nil, nil, fmt.Errorf("need signData")
	}	
	if _, ok = fields["timestamp"].(float64); !ok {
		return nil, nil, fmt.Errorf("need timestamp")
	} else {
		timestamp = int64(fields["timestamp"].(float64))	// 返回整数
	}
	if data, ok = fields["data"].(map[string]interface{}); !ok {
		return nil, nil, fmt.Errorf("need data")
	}	

	// 取得 密钥文件 路径
	secretPath, ok := SECRET_KEY[appId]
	if !ok {
		return nil, nil, fmt.Errorf("wrong appId")
	}

	// 取得用户信息
	me, err := client.GetMe(keyFileRoot+secretPath)
	if err!=nil {
		return nil, nil, err
	}

	// 获取 secret，用户密钥的签名串
	secret := base64.StdEncoding.EncodeToString(me.SignKey.Bytes())

	// 检查版本
	if version!="1" {
		return nil, nil, fmt.Errorf("wrong version")
	}

	// 检查签名类型
	if signType!="SHA256" {
		return nil, nil, fmt.Errorf("unknown signType")
	}

	// 生成参数的key，并排序
	keys := getMapKeys(fields)
	sort.Strings(*keys)
	//fmt.Println(*keys)

	// data 串，用于验签， map已按key排序
	dataStr, _ := json.Marshal(data)

	// 拼接验签串
	var signString = string("")
	for _,k:= range *keys {
		if k=="sign_data" {
			continue
		}
		if k=="data" {
			signString += k + "=" + string(dataStr) + "&"
		} else if k=="timestamp" {
			signString += k + "=" + strconv.FormatInt(timestamp, 10) + "&"
		} else {
			signString += k + "=" + fields[k].(string) + "&"
		}
	}
	signString += "key=" + secret
	//fmt.Println(signString)

	h := sha256.New()
	h.Write([]byte(signString))
	sum := h.Sum(nil)
	sha256Str := fmt.Sprintf("%x", sum)
	signStr := base64.StdEncoding.EncodeToString([]byte(sha256Str))
	//fmt.Println(sha256Str)

	if signStr!=signData {
		fmt.Println(signStr)
		fmt.Println(signData)		
		return nil, nil, fmt.Errorf("wrong signature")
	}

	return &data, me, nil
}

// 返回 map 所有 key
func getMapKeys(m map[string]interface{}) *[]string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return &keys
}
