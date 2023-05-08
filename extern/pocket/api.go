package pocket

import (
	"encoding/json"
	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/utils"
	"github.com/tendermint/tendermint/libs/bytes"
	"time"
)

func PoktApiGetHeight(url string) (uint64, error) {
	params := &HeightPoktParams{}
	apiUrl := utils.UrlJoin(url, "/query/height")
	response, err := web.HttpPostNoToken(apiUrl, params)
	if err != nil {
		GetLog().Error(err)
		return 0, err
	}

	poktRes := HeightData{}
	err = json.Unmarshal([]byte(response), &poktRes)
	if err != nil {
		GetLog().Error(err)
		return 0, err
	}

	return poktRes.Height, nil
}

// Info about the node's syncing state

type SyncInfo struct {
	LatestBlockHash   bytes.HexBytes `json:"latest_block_hash"`
	LatestAppHash     bytes.HexBytes `json:"latest_app_hash"`
	LatestBlockHeight string         `json:"latest_block_height"`
	LatestBlockTime   time.Time      `json:"latest_block_time"`

	EarliestBlockHash   bytes.HexBytes `json:"earliest_block_hash"`
	EarliestAppHash     bytes.HexBytes `json:"earliest_app_hash"`
	EarliestBlockHeight string         `json:"earliest_block_height"`
	EarliestBlockTime   time.Time      `json:"earliest_block_time"`

	CatchingUp bool `json:"catching_up"`
}

type ResultStatus struct {
	SyncInfo SyncInfo `json:"sync_info"`
}
type TmStatusResponse struct {
	JsonRpc string       `json:"jsonrpc"`
	Id      int          `json:"id"`
	Result  ResultStatus `json:"result"`
}

func PoktApiGetSync() (bool, error) {
	url := "http://127.0.0.1:26657/status"
	response, err := web.HttpGetNoToken(url, "")
	if err != nil {
		GetLog().Error(err)
		return false, err
	}

	GetLog().Debug("Status Response:", string(response))
	statusRes := TmStatusResponse{}
	err = json.Unmarshal([]byte(response), &statusRes)
	if err != nil {
		GetLog().Error(err)
		return false, err
	}

	diff := time.Since(statusRes.Result.SyncInfo.LatestBlockTime).Seconds()
	GetLog().Debug("Latest block was out ", diff, " seconds ago.")

	return diff < 60*30, nil
}
