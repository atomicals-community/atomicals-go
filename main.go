package main

import (
	"github.com/atomicals-core/atomicals"
	"github.com/atomicals-core/pkg/conf"
)

func main() {
	conf, err := conf.ReadJSONFromJSFile("./conf/config.json")
	if err != nil {
		panic(err)
	}
	var a *atomicals.Atomicals
	if conf.SelectDB == "memory" {
		a = atomicals.NewAtomicalsWithMemory(conf)
	} else {
		a = atomicals.NewAtomicalsWithSQL(conf)
	}
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
