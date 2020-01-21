package exchanges

type Exchange interface {
	QueryPrice(market string) (string, error)
}
