package exchanges

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinanceExchange_QueryPrice(t *testing.T) {
	exchange := BinanceExchange{addr: "https://api.binance.com"}
	price, err := exchange.QueryPrice("btc")
	require.Nil(t, err)
	fmt.Println(price)
}
