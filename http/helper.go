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
	//"crypto/md5"
	"io/ioutil"
	"github.com/valyala/fasthttp"
)

var (
	/* 接口验签使用 appid : appsecret (md5sum : sha1sum|base64) */
	APPID_SECRET = map[string]string{
		"bdecaa718f290152925e8d570c71adfe" : "YWQ2YjZjNmE3MTVjZTNlNzhiMjk2YjI2MGYyYzI2ZDllNGUyMjRiNyAgLQo=",
		"1ff3a3d2c1a8c236423ea3fe7bbdcff6" : "ZDlmZjk2YmNlMTEyNDYzN2E4ZGRlMWJhMTYyZDcxZDIxMjRkYTIwZiAgLQo=",
		"4fcf3871f4a023712bec9ed44ee4b709" : "MjdjNGQxNGU3NjA1OWI0MGVmODIyN2FkOTEwYTViNDQzYTNjNTIyNSAgLQo=",
	}

	/* userid 为 公钥base64 */
	SECRET_KEY = make(map[string]*client.User)
)

/* 装入用户密钥 */
func loadSecretKey(path string) error{
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		// 用户keyfile在子目录下
		if !f.IsDir(){
			continue
		}
		// 装入 keyfile
		u, err := client.GetMe(path + "/" + f.Name())
		if err!=nil {
			return err
		}

		pubkey := base64.StdEncoding.EncodeToString(u.CryptoPair.PubKey[:])
		SECRET_KEY[pubkey] = u // 保存用户信息
	}

	return nil
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
func checkSign(content []byte) (*map[string]interface{}, error) {
	fields := make(map[string]interface{})
	if err := json.Unmarshal(content, &fields); err != nil {
		return nil, err
	}

	var appId, version, signType, signData string
	var timestamp int64
	var data map[string]interface{}
	var ok bool

	// 检查参数
	if appId, ok = fields["appid"].(string); !ok {
		return nil, fmt.Errorf("need appid")
	}	
	//if userId, ok = fields["userkey"].(string); !ok { // 链用户公钥base64
	//	return nil, nil, fmt.Errorf("need userkey")
	//}	
	if version, ok = fields["version"].(string); !ok {
		return nil, fmt.Errorf("need version")
	}	
	if signType, ok = fields["sign_type"].(string); !ok {
		return nil, fmt.Errorf("need sign_type")
	}	
	if signData, ok = fields["sign_data"].(string); !ok {
		return nil, fmt.Errorf("need sign_data")
	}	
	if _, ok = fields["timestamp"].(float64); !ok {
		return nil, fmt.Errorf("need timestamp")
	} else {
		timestamp = int64(fields["timestamp"].(float64))	// 返回整数
	}
	if data, ok = fields["data"].(map[string]interface{}); !ok {
		return nil, fmt.Errorf("need data")
	}	

	// 获取 secret，用户密钥的签名串
	secret, ok := APPID_SECRET[appId]
	if !ok {
		return nil, fmt.Errorf("wrong appId")
	}

	// 取得用户信息
	//me, ok := SECRET_KEY[userId]
	//if !ok {
	//	return nil, nil, fmt.Errorf("wrong userId")
	//}

	// 检查版本
	if version!="1" {
		return nil, fmt.Errorf("wrong version")
	}

	// 检查签名类型
	if signType!="SHA256" {
		return nil, fmt.Errorf("unknown signType")
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
	fmt.Println(signString)

	h := sha256.New()
	h.Write([]byte(signString))
	sum := h.Sum(nil)
	sha256Str := fmt.Sprintf("%x", sum)
	signStr := base64.StdEncoding.EncodeToString([]byte(sha256Str))
	//fmt.Println(sha256Str)

	if signStr!=signData {
		fmt.Println(signStr)
		fmt.Println(signData)		
		return nil, fmt.Errorf("wrong signature")
	}

	return &data, nil
}

// 返回 map 所有 key
func getMapKeys(m map[string]interface{}) *[]string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return &keys
}
