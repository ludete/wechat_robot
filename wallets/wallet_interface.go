package wallets

type WalletInterface interface {
	SendMoney(walletID string, news []TransferNews) (string, error)
	GetAmountOfDenoms(credentialID string, denom string) (string, int, error)
	CreateUserWallet() (string, error)
}
