package main

import (
	"flag"
	"fmt"

	"github.com/atomicals-go/atomicals-api/internal/config"
	"github.com/atomicals-go/atomicals-api/internal/handler"
	"github.com/atomicals-go/atomicals-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/main-api.yaml", "the config file")
var atomicalsConfigFilePath = "./conf/config.json"

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c, atomicalsConfigFilePath)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
