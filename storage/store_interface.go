package storage

// 用户的最新账户余额;
// key : N + accountID + denom;  value : amount
// 用户的历史交易记录
// key : O + accountID + time(big endian) + denom; value : amount

type DB interface {
	UpdateTokenAmount(addr string, denom string, amount int) error
}
