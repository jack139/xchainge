package http

import (
	"log"
	"time"
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/AubSs/fasthttplogger"
)

/* 返回值的 content-type */
var (
	strContentType = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

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


/* 根返回 */
func index(ctx *fasthttp.RequestCtx) {
	log.Printf("%v", ctx.RemoteAddr())
	ctx.WriteString("Hello world.")
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
