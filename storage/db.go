package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/syndtr/goleveldb/leveldb"
)

// 用户的最新账户余额;
// key : N + accountID + denom;  value : amount
// 用户的历史交易记录
// key : O + accountID + denom + time(big endian); value : amount

type GoLevelDB struct {
	db    *leveldb.DB
	batch *leveldb.Batch
}

func NewDB(config string) DB {
	db, err := leveldb.OpenFile(config, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to open db %s\n", config))
	}
	return &GoLevelDB{
		db:    db,
		batch: new(leveldb.Batch),
	}
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

func (g *GoLevelDB) BuyToken(addr string, denom string, amount int) error {
	balanceKey := generateBalanceKey(addr, denom)
	balance, err := g.getUserBalance(balanceKey)
	if err != nil {
		return err
	}
	total := calCurrBalance(balance, amount)
	balanceVal := make([]byte, 8)
	binary.BigEndian.PutUint64(balanceVal, total)
	g.batch.Put(balanceKey, balanceVal)

	recordVal := make([]byte, 8)
	recodeKey := generateRecordKey(addr, BUYRECORD, denom)
	binary.BigEndian.PutUint64(recordVal, uint64(amount))
	g.batch.Put(recodeKey, recordVal)
	err = g.db.Write(g.batch, nil)
	return err
}

//key: R + addr +
func (g *GoLevelDB) ReceiveRMB(addr string, amount int) error {
	key := generateReceiveRMBKey(addr)
	if len(key) == 0 {
		log.Error("generate receiveRMBKey failed")
		panic("generate receiveRMBKey failed")
	}

	buf := bytes.NewBuffer(nil)
	amountVal := make([]byte, 4)
	binary.BigEndian.PutUint64(amountVal, uint64(amount))
	buf.Write(amountVal)
	buf.Write([]byte{BUYOPEN})
	val := buf.Bytes()
	err := g.db.Put(key, val, nil)
	return err
}

func generateReceiveRMBKey(addr string) []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte{RECEIVERMB})
	buf.Write([]byte(addr))
	return buf.Bytes()
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

func generateRecordKey(addr string, keyType byte, denom string) []byte {
	now := time.Now()
	buf := bytes.NewBuffer(nil)
	if length, err := buf.Write([]byte{keyType}); err != nil || length != 1 {
		panic(fmt.Sprintf("generateBalanceKey failed; identify : %c time : %d\n", keyType, now.Unix()))
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
