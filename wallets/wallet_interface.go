package wallets

type WalletInterface interface {
	SendMoney(from, to string, denom string, amount int)bool
	GetAmountOfDenom(addr string, denom string)int
	GetAllAmounts(addr string)map[string]int
}
