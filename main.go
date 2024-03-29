package main

import (
	"net"
	_ "net/http/pprof"

	"github.com/atomicals-core/atomicals"
	"github.com/atomicals-core/pkg/conf"
	"github.com/atomicals-core/pkg/log"
)

const port = "127.0.0.1.9999"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Log.Panicf("Port is in use...")
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

	// traced Specific Tx
	// tx, err := b.GetTransaction("b28f089b5a96c4803db73d51ed801aec4efec997761ee8dc914e0f934b6fcd59")
	// if err != nil {
	// 	log.Log.Panicf("GetTransaction err:%v", err)
	// }
	// a.TraceTx(*tx, height)
}
