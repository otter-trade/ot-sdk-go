package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	// "reflect"
	"sort"
	"testing"
	"time"

	// "github.com/otter-trade/coin-exchange-api/config"

	"github.com/otter-trade/coin-exchange-api/exchanges/request"
	"github.com/otter-trade/coin-exchange-api/exchanges/sharedtestvalues"

	configBinance "github.com/otter-trade/coin-exchange-api/config"
	"github.com/otter-trade/coin-exchange-api/exchanges/binance"
	"github.com/otter-trade/coin-exchange-api/exchanges/okx"
	// "github.com/otter-trade/spiders-coin-serve/spiders-rpc/model"
	// "github.com/otter-trade/spiders-coin-serve/spiders-rpc/internal/config"
	// "github.com/otter-trade/spiders-coin-serve/spiders-rpc/internal/server"
	// "github.com/otter-trade/spiders-coin-serve/spiders-rpc/internal/svc"
	// "github.com/otter-trade/spiders-coin-serve/spiders-rpc/pb"
	// "github.com/zeromicro/go-zero/core/conf"
)

const (
	apiKey                  = ""
	apiSecret               = ""
	canManipulateRealOrders = false
	BitcoinDonationAddress  = "bc1qk0jareu4jytc0cfrhr5wgshsq8282awpavfahc"
)

var b = &binance.Binance{}
var initFlagBinance = false
var initFlagOKX = false

func initBinance() {
	cfg := configBinance.GetConfig()
	if !initFlagBinance {
		err := cfg.LoadConfig("./configtest.json", true)
		if err != nil {
			log.Fatal("Binance load config error", err)
		} else {
			initFlagBinance = true
		}
	} else {
		return
	}

	binanceConfig, err := cfg.GetExchangeConfig("Binance")
	if err != nil {
		log.Fatal("Binance Setup() init error", err)
	}
	binanceConfig.API.AuthenticatedSupport = true
	binanceConfig.API.Credentials.Key = apiKey
	binanceConfig.API.Credentials.Secret = apiSecret
	b.SetDefaults()
	b.Websocket = sharedtestvalues.NewTestWebsocket()
	err = b.Setup(binanceConfig)
	if err != nil {
		log.Fatal("Binance setup error", err)
	}
	//b.setupOrderbookManager()
	request.MaxRequestJobs = 10
	b.Websocket.DataHandler = sharedtestvalues.GetWebsocketInterfaceChannelOverride()
	log.Printf(sharedtestvalues.LiveTesting, b.Name)
	err = b.UpdateTradablePairs(context.Background(), true)
	if err != nil {
		log.Fatal("Binance setup error", err)
	}
}

const (
	passphrase = ""
)

var ok = &okx.Okx{}

func initOkx() {
	cfg := configBinance.GetConfig()
	if !initFlagOKX {
		err := cfg.LoadConfig("./configtest.json", true)
		if err != nil {
			log.Fatal("okx load config error", err)
		} else {
			initFlagOKX = true
		}
	} else {
		return
	}

	exchCfg, err := cfg.GetExchangeConfig("Okx")
	if err != nil {
		log.Fatal(err)
	}
	exchCfg.API.Credentials.Key = apiKey
	exchCfg.API.Credentials.Secret = apiSecret
	exchCfg.API.Credentials.ClientID = passphrase
	ok.SetDefaults()
	if apiKey != "" && apiSecret != "" && passphrase != "" {
		exchCfg.API.AuthenticatedSupport = true
		exchCfg.API.AuthenticatedWebsocketSupport = true
	}
	ok.Websocket = sharedtestvalues.NewTestWebsocket()
	err = ok.Setup(exchCfg)
	if err != nil {
		log.Fatal(err)
	}
	request.MaxRequestJobs = 200
	ok.Websocket.DataHandler = sharedtestvalues.GetWebsocketInterfaceChannelOverride()
	ok.Websocket.TrafficAlert = sharedtestvalues.GetWebsocketStructChannelOverride()
}

type MarketCandlesReq struct {
	InstId string `json:"instId" validate:"endswith=-USDT"`
	Bar    string `json:"bar" validate:"oneof=1min 5min 10min 15min 30min"`
	Before int64  `json:"before,optional"` //稍后改为时间戳
	Limit  int64  `json:"limit,optional"`
}

type MarketCandlesResp struct {
	List []*Candle `json:"list"`
}
type Body struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type Candle struct {
	ID                      string `json:"_id,omitempty" json:"id,omitempty"`
	TimeStamp               int64  `json:"time_stamp,omitempty" json:"time_stamp,omitempty"`
	OpenTime                int64  `json:"open_time,omitempty"`
	OpenPrice               string `json:"open_price,omitempty"`
	HighestPrice            string `json:"highest_price,omitempty"`
	LowestPrice             string `json:"lowest_price,omitempty"`
	ClosePrice              string `json:"close_price,omitempty"`
	Volume                  string `json:"volume,omitempty"`
	QuoteAssetVolume        string `json:"quote_asset_volume,omitempty"`
	OKXTimeStamp            int64  `json:"OKX_time_stamp,omitempty"`
	OKXOpenTime             int64  `json:"OKX_open_time,omitempty"`
	OKXOpenPrice            string `json:"OKX_open_price,omitempty"`
	OKXHighestPrice         string `json:"OKX_highest_price,omitempty"`
	OKXLowestPrice          string `json:"okx_lowest_price,omitempty"`
	OKXClosePrice           string `json:"okx_close_price,omitempty"`
	OKXVolume               string `json:"okx_volume,omitempty"`
	OKXQuoteAssetVolume     string `json:"okx_quote_asset_volume,omitempty"`
	BinanceTimeStamp        int64  `json:"binance_time_stamp,omitempty"`
	BinanceOpenTime         int64  `json:"binance_open_time,omitempty"`
	BinanceOpenPrice        string `json:"binance_open_price,omitempty"`
	BinanceHighestPrice     string `json:"binance_highest_price,omitempty"`
	BinanceLowestPrice      string `json:"binance_lowest_price,omitempty"`
	BinanceClosePrice       string `json:"binance_close_price,omitempty"`
	BinanceVolume           string `json:"binance_volume,omitempty"`
	BinanceQuoteAssetVolume string `json:"binance_quote_asset_volume,omitempty"`
	UpdateAt                int64  `json:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt                int64  `json:"createAt,omitempty" json:"createAt,omitempty"`
}

func TestGetInstrument(t *testing.T) {
	initOkx()
	resp, err := ok.GetInstruments(context.Background(), &okx.InstrumentsFetchParams{
		InstrumentType: "OPTION",
		Underlying:     "SOL-USD",
	})
	if err != nil {
		fmt.Println("GetInstruments err ", err)
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error marshaling response to JSON:", err)
		return
	}

	fmt.Println("JSON response:", string(jsonResp))
}

func TestOKXGetTickers(t *testing.T) {
	initOkx()
	resp, err := ok.GetTickers(context.Background(), "OPTION", "", "SOL-USD")
	if err != nil {
		t.Error("Okx GetTickers() error", err)
	} else {
		fmt.Println("Okx GetTickers len: \n", len(resp))
		//fmt.Println("Okx GetTickers: \n", resp)
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error marshaling response to JSON:", err)
		return
	}

	fmt.Println("JSON response:", string(jsonResp))
}

func TestGetExchangeInfo(t *testing.T) {
	initBinance()
	resp, err := b.GetExchangeInfo(context.Background())
	if err != nil {
		t.Error(err)
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error marshaling response to JSON:", err)
		return
	}

	fmt.Println("JSON response:", string(jsonResp))

}

func TestBinanceGetTickers(t *testing.T) {
	initBinance()
	resp, err := b.GetTickers(context.Background())
	if err != nil {
		t.Error("Binance TestGetTickers error", err)
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error marshaling response to JSON:", err)
		return
	}

	fmt.Println("JSON response:", string(jsonResp))
}

func TestDelay(t *testing.T) {
	// //bianace延时测试
	// Start := time.Now().Add(-1 * time.Minute * time.Duration(100) * time.Duration(1))
	// pair := currency.NewPair(currency.BTC, currency.USDT)
	// startBinanceZQ := time.Now()
	// resp, err := b.GetSpotKline(context.Background(),
	// 	&binance.KlinesRequestParams{
	// 		Symbol:    pair,
	// 		Interval:  kline.OneMin.Short(),
	// 		Limit:     2,
	// 		StartTime: Start,
	// 		EndTime:   time.Now(),
	// 	})
	// endBinanceZQ := time.Now()
	// if err != nil {
	// 	fmt.Println("Binance获取最近1m的数据出错:", err)
	// } else {
	// 	fmt.Println("Binance最近1m的数据: ", len(resp))
	// }

	// //okx延时测试
	// startOKXZQ := time.Now()
	// respOkx, err := ok.GetCandlesticks(context.Background(), "BTC-USDT", kline.OneMin, Start, time.Now(), 2)
	// endOKXZQ := time.Now()
	// if err != nil {
	// 	fmt.Printf("connect redis failed! err : %v\n", err)
	// } else {
	// 	fmt.Println("okx resp: ", len(respOkx))
	// }

	// //okx http延时测试
	// before := time.Now().Add(-time.Hour * 24 * 500)
	// after := time.Now()
	// var rateLimit request.EndpointLimit = 121
	// var route string = "market/candles"

	// params := url.Values{}
	// params.Set("instId", "BTC-USDT")
	// var respHttp [][7]string
	// params.Set("limit", strconv.FormatInt(100, 10))
	// if !before.IsZero() {
	// 	params.Set("before", strconv.FormatInt(before.UnixMilli(), 10))
	// }
	// if !after.IsZero() {
	// 	params.Set("after", strconv.FormatInt(after.UnixMilli(), 10))
	// }
	// bar := ok.GetIntervalEnum(kline.OneHour, true)
	// if bar != "" {
	// 	params.Set("bar", bar)
	// }

	// startOKXHttp := time.Now()
	// err = ok.SendHTTPRequest(context.Background(), exchange.RestSpot, rateLimit, http.MethodGet, common.EncodeURLValues(route, params), nil, &respHttp, false)
	// endOKXHttp := time.Now()
	// if err != nil {
	// 	fmt.Printf("okx SendHTTPRequest error%+v", err)
	// }

	// //bianace http延时测试
	// symbol, err := b.FormatSymbol(pair, asset.Spot)
	// if err != nil {
	// 	fmt.Printf("bianace FormatSymbol error%+v", err)
	// }
	// params = url.Values{}
	// params.Set("symbol", symbol)
	// params.Set("interval", kline.OneMin.Short())
	// params.Set("limit", strconv.Itoa(100))
	// params.Set("startTime", strconv.FormatInt(Start.UnixMilli(), 10))
	// params.Set("endTime", strconv.FormatInt(time.Now().UnixMilli(), 10))
	// candleStick := "/api/v3/klines"
	// path := candleStick + "?" + params.Encode()
	// var respBinanceHttp interface{}
	// var spotDefaultRate request.EndpointLimit = 0

	// startBinanceHttp := time.Now()
	// err = b.SendHTTPRequest(context.Background(),
	// 	exchange.RestSpotSupplementary,
	// 	path,
	// 	spotDefaultRate,
	// 	&respBinanceHttp)
	// endBinanceHttp := time.Now()
	// if err != nil {
	// 	fmt.Printf("bianace SendHTTPRequest error%+v", err)
	// }

	//open-api延时测试
	data := &MarketCandlesReq{
		InstId: "BTC-USDT",
		Bar:    "1min",
		//Before: 0,
		Limit: 100,
	}

	bufs, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("json marshal failed:%+v\n", err)
	}

	// t.Logf("%s\n", bufs)
	startOpenAPI := time.Now()
	respOpenAPI, err := http.Post("http://test-api.ottertrade.com/open-api/market/candles", "application/json", bytes.NewBuffer(bufs))
	endOpenAPI := time.Now()
	if err != nil {
		fmt.Printf("HTTP POST request failed: %v\n", err)
		return
	}

	// 处理响应
	if respOpenAPI.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code: %d\n", respOpenAPI.StatusCode)
		return
	}
	defer respOpenAPI.Body.Close()

	byts, err := io.ReadAll(respOpenAPI.Body)
	if err != nil {
		t.Fatalf("read body  failed:%+v\n", err)
	}

	//t.Logf("%s\n", byts)

	var res *Body
	err = json.Unmarshal(byts, &res)
	if err != nil {
		t.Fatalf("json unMarshal failed:%+v\n", err)
	}

	var respList *MarketCandlesResp

	dataByts, err := json.Marshal(res.Data)
	if err != nil {
		t.Fatalf("json Marshal data  failed:%+v\n", err)
	}

	err = json.Unmarshal(dataByts, &respList)
	if err != nil {
		t.Fatalf("json Unmarshal data  failed:%+v\n", err)
	}

	//res 处理
	if res != nil && len(respList.List) > 0 {
		// 测试数据的连贯性
		fmt.Printf("has data \n")
		for key, val := range respList.List {
			preIndex := key - 1
			if preIndex < 0 {
				preIndex = 0
			}
			preItem := respList.List[preIndex]
			nowItem := respList.List[key]
			if key > 0 {
				if nowItem.OpenTime-preItem.OpenTime != -60 {
					fmt.Printf("opentime数据不连续: %+v, %+v", nowItem.OpenTime, preItem.OpenTime)
					fmt.Printf("opentime数据不连续: %+v, %+v", val.OpenTime, key)
					return
				}
				if nowItem.OKXOpenTime-preItem.OKXOpenTime != -60 {
					fmt.Printf("OKXOpenTime数据不连续: %+v, %+v", val.OpenTime, key)
					return
				}
				if nowItem.BinanceOpenTime-preItem.BinanceOpenTime != -60 {
					fmt.Printf("BinanceOpenTime数据不连续: %+v, %+v", val.OpenTime, key)
					return
				}
				if (nowItem.BinanceOpenTime != nowItem.OKXOpenTime) ||
					(nowItem.OpenTime != nowItem.OKXOpenTime) ||
					(nowItem.OpenTime != nowItem.BinanceOpenTime) {
					fmt.Printf("数据不相等: %+v, %+v", val.OpenTime, key)
					return
				}

				//fmt.Printf("数据没问题: %+v, %+v \n", val.OpenTime, key)
			}
		}
	} else {
		fmt.Printf("no data \n")
	}

	//计算延时时间
	// duration := endOKXZQ.Sub(startOKXZQ)
	// milliseconds := duration.Milliseconds()
	// // fmt.Printf("OKXZQ时间延迟为 %d 毫秒\n", milliseconds)

	// duration = endBinanceZQ.Sub(startBinanceZQ)
	// milliseconds = duration.Milliseconds()
	// // fmt.Printf("BinanceZQ时间延迟为 %d 毫秒\n", milliseconds)

	// duration = endOKXHttp.Sub(startOKXHttp)
	// milliseconds = duration.Milliseconds()
	// fmt.Printf("okx Http时间延迟为 %d 毫秒\n", milliseconds)

	// duration = endBinanceHttp.Sub(startBinanceHttp)
	// milliseconds = duration.Milliseconds()
	// fmt.Printf("binance Http时间延迟为 %d 毫秒\n", milliseconds)

	duration := endOpenAPI.Sub(startOpenAPI)
	milliseconds := duration.Milliseconds()
	fmt.Printf("OpenAPI Http时间延迟为 %d 毫秒\n", milliseconds)
}

func TestDelayLoop(t *testing.T) {
	// initOkx()
	// initBinance()
	for i := 0; i < 100; i++ {
		go TestDelay(t)
		time.Sleep(400)
	}
	select {}
}

// func TestFindLast(t *testing.T) {
// 	var c config.Config
// 	conf.MustLoad("../etc/spiders.yaml", &c)
// 	svcCtx := svc.NewServiceContext(c)
// 	resp, err := svcCtx.OkxItemModel1_BTC.FindLast(context.Background(), &model.OkxItem{})
// 	fmt.Printf("len(resp):%d\n", len(resp))
// 	if err != nil {
// 		fmt.Printf("err: %+v\n", err)
// 	} else if len(resp) > 0 {
// 		t := reflect.TypeOf(*resp[0])
// 		for i := 0; i < t.NumField(); i++ {
// 			field := t.Field(i)
// 			fmt.Printf("参数:%s, 类型:%s\n", field.Name, field.Type)
// 		}
// 	}

// 	byts, err := json.Marshal(resp)
// 	if err != nil {
// 		t.Fatalf("err:%+v\n", err)
// 	}
// 	t.Logf("%s\n", byts)
// }

//func TestLoop(t *testing.T) {
//	n := 1
//	setTime := "0 */" + strconv.Itoa(n) + " * * * ?"
//}

// func TestFindAll(t *testing.T) {

// 	var c config.Config
// 	conf.MustLoad("../etc/spiders.yaml", &c)
// 	svcCtx := svc.NewServiceContext(c)

// 	resp, err := svcCtx.OkxItemModel1_BTC.FindAll(context.Background(), &model.OkxItem{}, 1000)
// 	if err != nil {
// 		t.Fatalf("err:%+v\n", err)
// 	}
// 	// byts, err := json.Marshal(resp)
// 	// if err != nil {
// 	// 	t.Fatalf("err:%+v\n", err)
// 	// }
// 	// t.Logf("%s\n", byts)
// 	// 检查数据

// 	//检查数据
// 	for key, val := range resp {
// 		preIndex := key - 1
// 		if preIndex < 0 {
// 			preIndex = 0
// 		}
// 		preItem := resp[preIndex]
// 		nowItem := resp[key]
// 		if key > 0 {
// 			if nowItem.OpenTime-preItem.OpenTime != -60 {
// 				fmt.Printf("opentime数据不连续, nowItem.OpenTime: %+v, preItem.OpenTime: %+v \n", nowItem.OpenTime, preItem.OpenTime)
// 				fmt.Printf("opentime数据不连续, val.OpenTime: %+v,key: %+v \n", val.OpenTime, key)
// 				return
// 			}
// 			if nowItem.OKXOpenTime-preItem.OKXOpenTime != -60 {
// 				fmt.Printf("opentime数据不连续, val.OpenTime: %+v,key: %+v \n", val.OpenTime, key)
// 				return
// 			}
// 			if nowItem.BinanceOpenTime-preItem.BinanceOpenTime != -60 {
// 				fmt.Printf("opentime数据不连续, val.OpenTime: %+v,key: %+v \n", val.OpenTime, key)
// 				return
// 			}
// 			if (nowItem.BinanceOpenTime != nowItem.OKXOpenTime) ||
// 				(nowItem.OpenTime != nowItem.OKXOpenTime) ||
// 				(nowItem.OpenTime != nowItem.BinanceOpenTime) {
// 				fmt.Printf("opentime数据不连续, val.OpenTime: %+v,key: %+v \n", val.OpenTime, key)
// 				return
// 			}

// 			fmt.Printf("opentime数据不连续, val.OpenTime: %+v,key: %+v \n", val.OpenTime, key)
// 		}
// 	}
// }

// // 测试 rpc 方法
// func TestMongoAddBatch(t *testing.T) {

// 	var c config.Config
// 	conf.MustLoad("../etc/spiders.yaml", &c)
// 	svcCtx := svc.NewServiceContext(c)

// 	resp, err := server.NewSpidersServer(svcCtx).MongoAddBatch(context.Background(), &pb.MongoAddReqBatch{
// 		List: nil,
// 	})
// 	if err != nil {
// 		t.Fatalf("err:%+v\n", err)
// 	}
// 	byts, err := json.Marshal(resp)
// 	if err != nil {
// 		t.Fatalf("err:%+v\n", err)
// 	}
// 	t.Logf("%s\n", byts)

// }

type Person struct {
	Name string
	Age  int
	Tp   time.Time
}

func TestResverSlice(t *testing.T) {

	resp := []*Person{{
		Name: "a",
		Age:  1,
		Tp:   time.Now().Add(time.Hour * 2),
	}, {
		Name: "b",
		Age:  2,
		Tp:   time.Now().Add(time.Hour * 3),
	}, {
		Name: "c",
		Age:  3,
		Tp:   time.Now().Add(time.Hour * 4),
	}, {
		Name: "d",
		Age:  4,
		Tp:   time.Now().Add(time.Hour * 5),
	}}

	sort.Slice(resp, func(i, j int) bool {
		//return resp[i].Age > resp[j].Age
		return resp[i].Tp.After(resp[j].Tp)
	})
	byts, _ := json.Marshal(resp)

	t.Logf("%s\n", byts)
}
