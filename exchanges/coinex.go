package exchanges

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type CoinexExchange struct {
	addr string
}

func (c CoinexExchange) QueryPrice(market string) (string, error) {
	symbol := strings.ToLower(market + "cet")
	route := c.addr + "/v1/market/ticker?market=" + symbol
	res, err := http.Get(route)
	if err != nil {
		return "", fmt.Errorf("查询价格失败")
	}
	defer res.Body.Close()

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("查询价格失败")
	}
	var data CoinexResult
	err = json.Unmarshal(bz, &data)
	if err != nil {
		return "", fmt.Errorf("查询价格失败")
	}
	if data.Code != 0 || data.Message != "OK" {
		return "", fmt.Errorf("查询价格失败")
	}
	return data.Data.Ticker.Last + " cet", nil
}

type CoinexResult struct {
	Code int
	Data struct {
		Ticker struct {
			Last string
		}
	}
	Message string
}
