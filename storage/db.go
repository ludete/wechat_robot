package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

// 用户的最新账户余额;
// key : N + accountID + denom;  value : amount
// 用户的历史交易记录
// key : O + accountID + denom + time(big endian) + denom; value : amount

type GoLevelDB struct {
	db *leveldb.DB
}

func NewDB(config string) DB {
	db, err := leveldb.OpenFile(config, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to open db %s\n", config))
	}
	return &GoLevelDB{db: db}
}

func (g *GoLevelDB) getUserBalance(key []byte) (uint64, error) {
	balanceVal, err := g.db.Get(key, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return 0, err
	}

	balance := uint64(0)
	if len(balanceVal) > 0 {
		balance = binary.BigEndian.Uint64(balanceVal)
	}
	return balance, nil
}

func calCurrBalance(remain uint64, delta int) uint64 {
	var balance uint64
	if delta < 0 {
		if remain < uint64(-delta) {
			panic(fmt.Sprintf(""))
		}
		balance = remain - uint64(-delta)
	} else {
		if uint64(delta) > (math.MaxUint64 - remain) {
			panic("overflow uint64 amount")
		}
		balance = remain + uint64(delta)
	}
	return balance
}

func (g *GoLevelDB) UpdateTokenAmount(addr string, denom string, amount int) error {
	key := generateBalanceKey(addr, denom)
	balance, err := g.getUserBalance(key)
	if err != nil {
		return err
	}
	total := calCurrBalance(balance, amount)
	balanceVal := make([]byte, 8)
	binary.BigEndian.PutUint64(balanceVal, total)
	return g.db.Put(key, balanceVal, nil)
}

func generateBalanceKey(addr string, denom string) []byte {
	now := time.Now()
	buf := bytes.NewBuffer(nil)
	if length, err := buf.Write([]byte{BALANCE}); err != nil || length != 1 {
		panic(fmt.Sprintf("generateBalanceKey failed; identify : %c time : %d\n", BALANCE, now.Unix()))
	}
	if length, err := buf.Write([]byte(addr)); err != nil || length != len(addr) {
		panic(fmt.Sprintf("generateBalanceKey failed; addr : %s, time : %d\n", addr, now.Unix()))
	}
	if length, err := buf.Write([]byte(denom)); err != nil || length != len(denom) {
		panic(fmt.Sprintf("generateBalanceKey failed; denom : %s, time : %d\n", denom, now.Unix()))
	}
	return buf.Bytes()
}

func generateRecordKey(addr string, denom string) []byte {
	now := time.Now()
	buf := bytes.NewBuffer(nil)
	if length, err := buf.Write([]byte{RECORD}); err != nil || length != 1 {
		panic(fmt.Sprintf("generateBalanceKey failed; identify : %c time : %d\n", RECORD, now.Unix()))
	}
	if length, err := buf.Write([]byte(addr)); err != nil || length != len(addr) {
		panic(fmt.Sprintf("generateBalanceKey failed; addr : %s, time : %d\n", addr, now.Unix()))
	}
	if length, err := buf.Write([]byte(denom)); err != nil || length != len(denom) {
		panic(fmt.Sprintf("generateBalanceKey failed; denom : %s, time : %d\n", denom, now.Unix()))
	}
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(now.Unix()))
	if length, err := buf.Write(bz); err != nil || length != len(bz) {
		panic(fmt.Sprintf("generateBalanceKey failed; time : %d\n", now.Unix()))
	}
	return buf.Bytes()
}
