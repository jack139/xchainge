package http

import (
	"fmt"
	"sort"
	"strconv"
	"log"
	"time"
	"encoding/json"
	"encoding/base64"
	"crypto/sha256"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/AubSs/fasthttplogger"
)

/* 返回值的 content-type */
var (
	strContentType = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

var SECRET_KEY = map[string]string{
	"19E179E5DC29C05E65B90CDE57A1C7E5" : "D91CEB11EE62219CD91CEB11EE62219C",
	"66A095861BAE55F8735199DBC45D3E8E" : "43E554621FF7BF4756F8C1ADF17F209C",
	"75C50F018B34AC0240915EC685F5961B" : "BCB3DF17A794368E1BB0352D3D2D5F50",
	"3EA25569454745D01219080B779F021F" : "41DF0E6AE27B5282C07EF5124642A352",
}


/* 入口 */
func RunServer(port string) {
	/* router */
	r := router.New()
	r.GET("/", index)
	r.POST("/test", doNonthing)
	//r.POST("/qa", onQA)

	/* 启动server */
	s := &fasthttp.Server{
		Handler: fasthttplogger.Combined(r.Handler),
		Name: "FastHttpLogger",
	}
	log.Fatal(s.ListenAndServe(":"+port))
}


/* 根返回 */
func index(ctx *fasthttp.RequestCtx) {
	log.Printf("%v", ctx.RemoteAddr())
	ctx.WriteString("Hello world.")
}


/* 处理返回值，返回json */
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

	// 检查参数
	if _, ok := fields["appId"]; !ok {
		return nil, fmt.Errorf("need appId")
	}	
	if _, ok := fields["version"]; !ok {
		return nil, fmt.Errorf("need version")
	}	
	if _, ok := fields["signType"]; !ok {
		return nil, fmt.Errorf("need signType")
	}	
	if _, ok := fields["signData"]; !ok {
		return nil, fmt.Errorf("need signData")
	}	
	if _, ok := fields["timestamp"]; !ok {
		return nil, fmt.Errorf("need timestamp")
	}	
	if _, ok := fields["data"]; !ok {
		return nil, fmt.Errorf("need data")
	}	

	// 获取参数
	appId     := fields["appId"].(string)
	version   := fields["version"].(string)
	signType  := fields["signType"].(string)
	signData  := fields["signData"].(string)
	timestamp := int(fields["timestamp"].(float64))
	data := fields["data"].(map[string]interface{})

	// 取得 secret
	secret, ok := SECRET_KEY[appId]
	if !ok {
		return nil, fmt.Errorf("wrong appId")
	}

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
		if k=="signData" {
			continue
		}
		if k=="data" {
			signString += k + "=" + string(dataStr) + "&"
		} else if k=="timestamp" {
			signString += k + "=" + strconv.Itoa(timestamp) + "&"
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
	fmt.Println(signStr)
	fmt.Println(signData)

	if signStr!=signData {
		return nil, fmt.Errorf("wrong signature")
	}

	return &data, nil
}


func getMapKeys(m map[string]interface{}) *[]string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return &keys
}
