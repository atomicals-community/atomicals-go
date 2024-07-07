package main

import (
	"fmt"
	"net"
	"os"

	atomicals "github.com/atomicals-go/atomicals-core/operation"
	"github.com/atomicals-go/pkg/conf"
)

func main() {
	// use a port for single-mode
	port := "8080"
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error:%v, Cannot obtain port %s, Atomicals-Core is probably already running.\n", err, port)
		os.Exit(1)
	}
	defer ln.Close()
	conf, err := conf.ReadJSONFromJSFile("./conf/config.json")
	if err != nil {
		panic(err)
	}
	a := atomicals.NewAtomicalsWithSQL(conf)
	for {
		a.Run()
	}

}
