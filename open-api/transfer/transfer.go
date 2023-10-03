package transfer

import "errors"

type Candle struct {
	TimeStamp           int64  `json:"time_stamp"`
	OkxOpenPrice        string `json:"okx_open_price"`
	OkxHighestPrice     string `json:"okx_highest_price"`
	OkxLowestPrice      string `json:"okx_lowest_price"`
	OkxClosePrice       string `json:"okx_close_price"`
	OkxVolume           string `json:"okx_volume"`
	BinanceOpenPrice    string `json:"binance_open_price"`
	BinanceHighestPrice string `json:"binance_highest_price"`
	BinanceLowestPrice  string `json:"binance_lowest_price"`
	BinanceClosePrice   string `json:"binance_close_price"`
	BinanceVolume       string `json:"binance_volume"`
}

type MarketResult struct {
	Code int64              `json:"code"`
	Data map[string][][]any `json:"data"`
	Msg  string             `json:"msg"`
}

func Transfer(openApiData [][]any) (list []Candle, err error) {

	for _, data := range openApiData {

		timeStamp, ok := data[0].(float64)
		if !ok {
			return nil, errors.New("timeStamp transfer error")
		}

		OkxData, ok := data[1].([]any)
		if !ok {
			return nil, errors.New("OkxData transfer error")
		}

		if len(OkxData) < 5 {
			return nil, errors.New("OkxData 数据丢失")
		}

		BinanceData, ok := data[2].([]any)
		if !ok {
			return nil, errors.New("BinanceData transfer error")
		}

		if len(BinanceData) < 5 {
			return nil, errors.New("BinanceData 数据丢失")
		}

		//Okx
		okxOpenPrice, ok := OkxData[0].(string)
		if !ok {
			return nil, errors.New("okxOpenPrice transfer error")
		}
		okxHighestPrice, ok := OkxData[1].(string)
		if !ok {
			return nil, errors.New("okxHighestPrice transfer error")
		}
		okxLowestPrice, ok := OkxData[2].(string)
		if !ok {
			return nil, errors.New("okxLowestPrice transfer error")
		}
		okxClosePrice, ok := OkxData[3].(string)
		if !ok {
			return nil, errors.New("okxClosePrice transfer error")
		}
		okxVolume, ok := OkxData[4].(string)
		if !ok {
			return nil, errors.New("okxVolume transfer error")
		}

		//binance
		binanceOpenPrice, ok := BinanceData[0].(string)
		if !ok {
			return nil, errors.New("binanceOpenPrice transfer error")
		}
		binanceHighestPrice, ok := BinanceData[1].(string)
		if !ok {
			return nil, errors.New("binanceHighestPrice transfer error")
		}
		binanceLowestPrice, ok := BinanceData[2].(string)
		if !ok {
			return nil, errors.New("binanceLowestPrice transfer error")
		}
		binanceClosePrice, ok := BinanceData[3].(string)
		if !ok {
			return nil, errors.New("binanceClosePrice transfer error")
		}
		binanceVolume, ok := BinanceData[4].(string)
		if !ok {
			return nil, errors.New("binanceVolume transfer error")
		}

		list = append(list, Candle{
			TimeStamp:           int64(timeStamp),
			OkxOpenPrice:        okxOpenPrice,
			OkxHighestPrice:     okxHighestPrice,
			OkxLowestPrice:      okxLowestPrice,
			OkxClosePrice:       okxClosePrice,
			OkxVolume:           okxVolume,
			BinanceOpenPrice:    binanceOpenPrice,
			BinanceHighestPrice: binanceHighestPrice,
			BinanceLowestPrice:  binanceLowestPrice,
			BinanceClosePrice:   binanceClosePrice,
			BinanceVolume:       binanceVolume,
		})

	}

	return
}
