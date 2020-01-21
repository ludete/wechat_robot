package storage

// 用户的最新账户余额;
// key : N + accountID + denom;  value : amount
// 用户的历史交易记录
// key : O + accountID + time(big endian) + denom; value : amount

type DB interface {
	GetUserWalletKeyID(weChatID string) (string, error)
	PutUserWalletKeyID(weChatID string, walletID string) error
	GetUserDenomAddrInWallet(weChatID, walletID, denom string) (string, error)
	PutUserDenomAddrInWallet(weChatID, walletID, denom string, addr string) error
	BuyTokenRecord(weChatID string, denom string, amount int) error
	ReceiveRMB(weChatID string, amount int) error
	GetUserStoreRMB(weChatID string) (int, error)
	ClearUserStoreRMB(waChatID string)
}
