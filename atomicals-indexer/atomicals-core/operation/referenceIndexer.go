package atomicals

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo"
)

func ReferenceIndexer(data *repo.AtomicaslData, operation *witness.WitnessAtomicalsOperation) {
	hasAtomicalsAsset, err := fetchAtomicalsAssetFromReferenceIndexer(operation.AtomicalsID)
	if err != nil {
		log.Log.Infof("RevealLocationTxID:%v", operation.RevealLocationTxID)
		panic(err)
	}

	switch operation.Op {
	case "dmt":
	case "dft":
		if hasAtomicalsAsset && data.NewGlobalDistributedFt == nil {
			log.Log.Infof("%v %v %v", operation.RevealLocationHeight, operation.AtomicalsID, operation.RevealLocationTxID)
			panic("hasAtomicalsAsset")
		}
		if !hasAtomicalsAsset && data.NewGlobalDistributedFt != nil {
			log.Log.Infof("%v %v %v", operation.RevealLocationHeight, operation.AtomicalsID, operation.RevealLocationTxID)
			panic("hasAtomicalsAsset")
		}
	case "ft":
		if hasAtomicalsAsset && data.NewGlobalDirectFt == nil {
			log.Log.Infof("%v %v %v", operation.RevealLocationHeight, operation.AtomicalsID, operation.RevealLocationTxID)
			panic("hasAtomicalsAsset")
		}
		if !hasAtomicalsAsset && data.NewGlobalDirectFt != nil {
			log.Log.Infof("%v %v %v", operation.RevealLocationHeight, operation.AtomicalsID, operation.RevealLocationTxID)
			panic("hasAtomicalsAsset")
		}
	case "nft":
		if hasAtomicalsAsset && data.NewUTXONftInfo == nil {
			if !(operation.Payload.Args.RequestContainer != "" || operation.Payload.Args.RequestDmitem != "") {
				return
			}
			log.Log.Infof("%v %v %v", operation.RevealLocationHeight, operation.AtomicalsID, operation.RevealLocationTxID)
			panic("hasAtomicalsAsset")
		}
		if !hasAtomicalsAsset && data.NewUTXONftInfo != nil {
			log.Log.Infof("%v %v %v", operation.RevealLocationHeight, operation.AtomicalsID, operation.RevealLocationTxID)
			panic("hasAtomicalsAsset")
		}
	case "evt":
		panic(operation.Payload)
	case "dat":
	case "sl":
		panic(operation.Payload)
	default:
	}

}

func fetchAtomicalsAssetFromReferenceIndexer(atomicalsID string) (bool, error) {
	encodedTxID := url.QueryEscape(fmt.Sprintf(`"%s"`, atomicalsID))
	// endpoint := "https://ep.atomicals.xyz/proxy"
	// endpoint := "https://ep.atomicalswallet.com/proxy"
	endpoint := "https://ep.wizz.cash/proxy"
	// endpoint := "https://atomindexer.satsx.io/proxy"
	url := fmt.Sprintf(endpoint+"/blockchain.atomicals.get?params=[%s]", encodedTxID)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	// fmt.Println(string(body))
	var response RespAtomicalsAssetFromReferenceIndexer
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	return response.Response.Result.Confirmed, nil
}

type RespAtomicalsAssetFromReferenceIndexer struct {
	Success  bool `json:"success"`
	Response struct {
		Result struct {
			AtomicalID     string `json:"atomical_id"`
			AtomicalNumber int    `json:"atomical_number"`
			AtomicalRef    string `json:"atomical_ref"`
			Type           string `json:"type"`
			Confirmed      bool   `json:"confirmed"`
			MintInfo       struct {
				CommitTxid               string `json:"commit_txid"`
				CommitIndex              int    `json:"commit_index"`
				CommitLocation           string `json:"commit_location"`
				CommitTxNum              int    `json:"commit_tx_num"`
				CommitHeight             int    `json:"commit_height"`
				RevealLocationTxid       string `json:"reveal_location_txid"`
				RevealLocationIndex      int    `json:"reveal_location_index"`
				RevealLocation           string `json:"reveal_location"`
				RevealLocationTxNum      int    `json:"reveal_location_tx_num"`
				RevealLocationHeight     int    `json:"reveal_location_height"`
				RevealLocationHeader     string `json:"reveal_location_header"`
				RevealLocationBlockhash  string `json:"reveal_location_blockhash"`
				RevealLocationScripthash string `json:"reveal_location_scripthash"`
				RevealLocationScript     string `json:"reveal_location_script"`
				RevealLocationValue      int    `json:"reveal_location_value"`
				Args                     struct {
					MintAmount    int    `json:"mint_amount"`
					MintHeight    int    `json:"mint_height"`
					MaxMints      int    `json:"max_mints"`
					RequestTicker string `json:"request_ticker"`
					Bitworkc      string `json:"bitworkc"`
					// Nonce         string `json:"nonce"`
					Time int `json:"time"`
				} `json:"args"`
				Meta struct {
				} `json:"meta"`
				Ctx struct {
				} `json:"ctx"`
				RequestTicker string `json:"$request_ticker"`
				Bitwork       struct {
					Bitworkc string      `json:"bitworkc"`
					Bitworkr interface{} `json:"bitworkr"`
				} `json:"$bitwork"`
			} `json:"mint_info"`
			Subtype    string `json:"subtype"`
			MintMode   string `json:"$mint_mode"`
			MaxSupply  int    `json:"$max_supply"`
			MintHeight int    `json:"$mint_height"`
			MintAmount int    `json:"$mint_amount"`
			MaxMints   int    `json:"$max_mints"`
			Bitwork    struct {
				Bitworkc string      `json:"bitworkc"`
				Bitworkr interface{} `json:"bitworkr"`
			} `json:"$bitwork"`
			TickerCandidates []struct {
				TxNum                int    `json:"tx_num"`
				AtomicalID           string `json:"atomical_id"`
				Txid                 string `json:"txid"`
				CommitHeight         int    `json:"commit_height"`
				RevealLocationHeight int    `json:"reveal_location_height"`
			} `json:"$ticker_candidates"`
			RequestTickerStatus struct {
				Status             string `json:"status"`
				VerifiedAtomicalID string `json:"verified_atomical_id"`
				Note               string `json:"note"`
			} `json:"$request_ticker_status"`
			RequestTicker string `json:"$request_ticker"`
			Ticker        string `json:"$ticker"`
			MintData      struct {
				Fields struct {
					ImageJpg struct {
						Ct string `json:"$ct"`
						B  struct {
							D    string `json:"$d"`
							B    string `json:"$b"`
							Len  int    `json:"$len"`
							Auto bool   `json:"$auto"`
						} `json:"$b"`
					} `json:"image.jpg"`
					Args struct {
						MintAmount    int    `json:"mint_amount"`
						MintHeight    int    `json:"mint_height"`
						MaxMints      int    `json:"max_mints"`
						RequestTicker string `json:"request_ticker"`
						Bitworkc      string `json:"bitworkc"`
						// Nonce         string `json:"nonce"`
						Time int `json:"time"`
					} `json:"args"`
				} `json:"fields"`
			} `json:"mint_data"`
			DftInfo struct {
				MintCount int `json:"mint_count"`
			} `json:"dft_info"`
			LocationSummary struct {
				UniqueHolders     int `json:"unique_holders"`
				CirculatingSupply int `json:"circulating_supply"`
			} `json:"location_summary"`
		} `json:"result"`
	} `json:"response"`
}
