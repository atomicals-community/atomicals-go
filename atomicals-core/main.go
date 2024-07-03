package main

import (
	atomicals "github.com/atomicals-go/atomicals-core/operation"
	"github.com/atomicals-go/pkg/conf"
)

func main() {
	conf, err := conf.ReadJSONFromJSFile("../conf/config.json")
	if err != nil {
		panic(err)
	}
	a := atomicals.NewAtomicalsWithSQL(conf)
	for {
		a.Run()
	}

	// a.TraceSpecificTx()
}
