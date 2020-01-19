package exchanges

type Exchange interface {
	QueryPrice(market string) (int, error)
}
