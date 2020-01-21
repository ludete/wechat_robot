package exchanges

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type BinanceExchange struct {
	addr string
}

func (b BinanceExchange) QueryPrice(market string) (string, error) {
	symbol := strings.ToUpper(market + "USDT")
	route := b.addr + "/api/v3/ticker/price?symbol=" + symbol
	res, err := http.Get(route)
	if err != nil {
		return "", fmt.Errorf("查询价格失败")
	}
	defer res.Body.Close()

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("查询价格失败")
	}
	var data map[string]interface{}
	err = json.Unmarshal(bz, &data)
	if err != nil {
		return "", fmt.Errorf("response from biance exchange unmarshal json failed\n")
	}
	if price, ok := data["price"]; ok {
		return price.(string) + " usdt", nil
	}
	return "", fmt.Errorf("查询价格失败")
}
