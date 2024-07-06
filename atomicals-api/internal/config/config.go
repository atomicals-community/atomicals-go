package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	BtcRpcURL      string `json:"btc_rpc_url"`
	BtcRpcUser     string `json:"btc_rpc_user"`
	BtcRpcPassword string `json:"btc_rpc_password"`
	SqlDNS         string `json:"sql_dns"`
}
