package main

import (
	_ "net/http/pprof"

	"github.com/atomicals-core/api"
	"github.com/atomicals-core/atomicals"
	"github.com/atomicals-core/pkg/conf"
)

func main() {
	conf, err := conf.ReadJSONFromJSFile("./conf/config.json")
	if err != nil {
		panic(err)
	}

	// Run atomicals api with port:9000
	api.Run(conf)

	// Run atomicals indexer
	a := atomicals.NewAtomicalsWithSQL(conf)
	for {
		a.TraceBlock()
	}
}
