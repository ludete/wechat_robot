package storage

// 用户的最新账户余额;
// key : N + accountID + denom;  value : amount
// 用户的历史交易记录
// key : O + accountID + time(big endian) + denom; value : amount

type DB interface {
	BuyTokenRecord(weChatID string, denom string, amount int) error
	ReceiveRMB(weChatID string, amount int) error
	GetUserStoreRMB(weChatID string) (int, error)
	ClearUserStoreRMB(waChatID string)
	GetUserDeomAddr(weChatID string, denom string) (string, error)
	GetUserWalletKeyID(weChatID string) string
}
