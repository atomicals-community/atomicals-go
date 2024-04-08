package main

import (
	"net"
	_ "net/http/pprof"

	"github.com/atomicals-core/atomicals"
	"github.com/atomicals-core/pkg/conf"
	"github.com/atomicals-core/pkg/log"
)

const port = "127.0.0.1:9999"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Log.Panicf("......net.Listen error:%v", err)
	}
	defer listener.Close()

	conf, err := conf.ReadJSONFromJSFile("./conf/config.json")
	if err != nil {
		panic(err)
	}
	a := atomicals.NewAtomicalsWithSQL(conf)
	for {
		a.TraceBlock()
	}
}
