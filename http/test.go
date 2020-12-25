package http

import (
	"fmt"
	"sort"
	"strconv"
	//"log"
	//"time"
	"encoding/json"

	"github.com/valyala/fasthttp"
)


/*
type RespHTTP struct {
	appId    string
	code	int
	encType string
	success bool
	timestamp int32
	data    []byte
}
*/

var SECRET_KEY = map[string]string{
	"19E179E5DC29C05E65B90CDE57A1C7E5" : "D91CEB11EE62219CD91CEB11EE62219C",
	"66A095861BAE55F8735199DBC45D3E8E" : "43E554621FF7BF4756F8C1ADF17F209C",
	"75C50F018B34AC0240915EC685F5961B" : "BCB3DF17A794368E1BB0352D3D2D5F50",
	"3EA25569454745D01219080B779F021F" : "41DF0E6AE27B5282C07EF5124642A352",
}

func getMapKeys(m map[string]interface{}) *[]string {
	var keys []string
    for k := range m {
        keys = append(keys, k)
    }
    return &keys
}

/* 空接口, 只进行签名校验 */
func doNonthing(ctx *fasthttp.RequestCtx) {
	respJson := map[string] interface{} {
		"code" : 0,
		"msg" : "",
		"timestamp" : 0,
		"data" : "",
	}

	content := ctx.PostBody()

	fields := make(map[string]interface{})
	if err := json.Unmarshal(content, &fields); err != nil {
		fmt.Println(err)
		return
	}

	appId     := fields["appId"].(string)
	//version   := fields["version"].(string)
	//signType  := fields["signType"].(string)
	//signData  := fields["signData"].(string)
	timestamp := int(fields["timestamp"].(float64))
	data      := fields["data"].(map[string]interface{})

	dataStr, _ := json.Marshal(data)

	//fmt.Printf("%v\n%s\n", data, b)

	// 生成参数的key，并排序
	keys := getMapKeys(fields)
	sort.Strings(*keys)
	fmt.Println(*keys)

	// 取得 secret
	secret, ok := SECRET_KEY[appId]
	if !ok {
		respJson["code"] = 9000
		respJson["msg"] = "appid fail"
		//goto goodbye
	}
	fmt.Println(secret)

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
	signString = signString[:len(signString)-1]
	fmt.Println(signString)

	fmt.Printf("%v %d %v %v\n", fields, 
		int(fields["timestamp"].(float64)), 
		fields["data"].(map[string]interface{})["test"],
		fields["xxx"] )

//goodbye:

	doJSONWrite(ctx, fasthttp.StatusOK, respJson)
}

/* 
	使用 tf 进行问答，输入格式： 
		{"c":"背景知识", "q":"问题"}
*/
/*
func onQA(ctx *fasthttp.RequestCtx) {
	retJson := map[string] string {"code":"9000","msg":"error"}
	content := ctx.PostBody()
	fields := make(map[string]interface{})
	if err := json.Unmarshal(content, &fields); err != nil {
		retJson["code"] = "9001"; retJson["msg"] = "invalid json data"
		doJSONWrite(ctx, fasthttp.StatusOK, retJson)
		return
	}
	log.Printf("%v", fields)

	// 检查 c q 是否存在 
	corpus, ok := fields["c"]
	if !ok {
		retJson["code"] = "9010"; retJson["msg"] = "data error"
		doJSONWrite(ctx, fasthttp.StatusOK, retJson)
		return
	}
	question, ok := fields["q"]
	if !ok {
		retJson["code"] = "9020"; retJson["msg"] = "data error"
		doJSONWrite(ctx, fasthttp.StatusOK, retJson)
		return
	}

	// 检查类型是否是字符串 
	if _, ok := corpus.(string); !ok {
		retJson["code"] = "9030"; retJson["msg"] = "data type error"
		doJSONWrite(ctx, fasthttp.StatusOK, retJson)
		return		
	}
	if _, ok := question.(string); !ok {
		// not string 
		retJson["code"] = "9040"; retJson["msg"] = "data type error"
		doJSONWrite(ctx, fasthttp.StatusOK, retJson)
		return
	}

	// 调用问答 
	ans, err := bert.BertQA(corpus.(string), question.(string))
	if err!=nil {
		log.Printf("ERROR: %s", err)
		retJson["code"] = "9002"; retJson["msg"] = "backend error"
		doJSONWrite(ctx, fasthttp.StatusOK, retJson)		
	} else {
		log.Printf("ans: %v", ans)
		doJSONWrite(ctx, fasthttp.StatusOK, map[string] string {"code":"0","msg":"ok","data":ans})
	}
}

*/