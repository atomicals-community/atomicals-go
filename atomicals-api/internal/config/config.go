package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	SqlDNS string
}
