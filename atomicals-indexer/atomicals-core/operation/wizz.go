package atomicals

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Response struct {
	Success  bool `json:"success"`
	Response struct {
		Height int    `json:"height"`
		Op     string `json:"op"`
		TxNum  int    `json:"tx_num"`
		Txid   string `json:"txid"`
	} `json:"response"`
}

func fetchTxFromWizz(txID string) (string, error) {
	encodedTxID := url.QueryEscape(fmt.Sprintf(`"%s"`, txID))
	url := fmt.Sprintf("https://ep.wizz.cash/proxy/blockchain.atomicals.transaction?params=[%s]", encodedTxID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Response.Op, nil
}
