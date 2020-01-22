package exchanges

type Exchanges struct {
	coinex  *CoinexExchange
	binance *BinanceExchange
}

func NewExchanges(coinexURL, binaceURL string) *Exchanges {
	return &Exchanges{
		coinex:  &CoinexExchange{addr: coinexURL},
		binance: &BinanceExchange{addr: binaceURL},
	}
}

func (e *Exchanges) QueryPrice(market string) (string, error) {
	//switch strings.ToLower(market) {
	//case SPICE:
	//	return e.coinex.QueryPrice(market)
	//}
	//return e.binance.QueryPrice(market)
	return e.coinex.QueryPrice(market)
}
