package transfer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestTransfer(t *testing.T) {

	req := struct {
		InstId string `json:"instId" `
		Bar    string `json:"bar"`
		Before int64  `json:"before,optional"`
		Limit  int64  `json:"limit,optional"`
	}{
		InstId: "BTC-USDT",
		Bar:    "1min",
		Before: 0,
		Limit:  2,
	}

	reqByts, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json marshal err :%+v", err)
	}

	resp, err := http.Post("http://test-api.ottertrade.com/open-api/market/candles", "application/json", strings.NewReader(fmt.Sprintf("%s", reqByts)))
	if err != nil {
		t.Fatalf("request err :%+v", err)
	}
	defer resp.Body.Close()

	bodyByts, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body err :%+v", err)
	}

	res := &MarketResult{}

	err = json.Unmarshal(bodyByts, res)
	if err != nil {
		t.Fatalf("json marshal err :%+v", err)
	}

	candles, err := Transfer(res.Data["list"])
	if err != nil {
		t.Fatalf("open-api data transfer  err :%+v", err)
	}

	candlesByts, err := json.Marshal(candles)
	if err != nil {
		t.Fatalf("json marshal err :%+v", err)
	}

	t.Logf("%s\n", candlesByts)

}
