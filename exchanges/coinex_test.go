package exchanges

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCoinexExchange_QueryPrice(t *testing.T) {
	exchange := CoinexExchange{addr: "https://api.coinex.com"}
	price, err := exchange.QueryPrice("spice")
	require.Nil(t, err)
	fmt.Println(price)
}
