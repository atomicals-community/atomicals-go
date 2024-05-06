package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/atomicals-go/apiService/logic"
	"github.com/atomicals-go/pkg/conf"
	"github.com/beego/beego/v2/server/web"
	webContext "github.com/beego/beego/v2/server/web/context"
)

func EnableCORS() web.FilterFunc {
	return func(ctx *webContext.Context) {
		origin := ctx.Input.Header("Origin")
		if origin != "" {
			ctx.Output.Header("Access-Control-Allow-Origin", origin)
			ctx.Output.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			ctx.Output.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			ctx.Output.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			ctx.Output.Header("Access-Control-Allow-Credentials", "true")
		}
		if ctx.Input.Method() == "OPTIONS" {
			ctx.Output.SetStatus(200)
			ctx.ResponseWriter.WriteHeader(200)
			return
		}
	}
}

func main() {
	conf, err := conf.ReadJSONFromJSFile("../conf/config.json")
	if err != nil {
		panic(err)
	}

	web.InsertFilter("*", web.BeforeRouter, EnableCORS())
	web.AutoPrefix("/api/v1", &logic.Controller{})
	web.ErrorHandler("404", func(writer http.ResponseWriter, request *http.Request) {
		respBytes, _ := json.Marshal("unsupported interface")

		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(respBytes)
	})
	web.BConfig.CopyRequestBody = true
	logic.InitController(conf)
	web.Run()
	fmt.Println("service is running...")
}
