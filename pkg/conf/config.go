package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	BtcRpcURL            string `json:"btc_rpc_url"`
	BtcRpcUser           string `json:"btc_rpc_user"`
	BtcRpcPassword       string `json:"btc_rpc_password"`
	SqlDNS               string `json:"sql_dns"`
	AtomicalsStartHeight int64  `json:"atomicals_start_height"`
}

func ReadJSONFromJSFile(filePath string) (*Config, error) {
	// Open the .js file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the struct
	var myData Config
	err = json.Unmarshal(byteValue, &myData)
	if err != nil {
		return nil, err
	}

	return &myData, nil
}
