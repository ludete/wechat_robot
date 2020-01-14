package exchanges

type QueryMarketPrice interface {
	QueryPrice(market string)int
}