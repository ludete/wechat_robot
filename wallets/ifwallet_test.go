package wallets

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIFWallet(t *testing.T) {
	ifwallet := NewWallet("http://47.52.172.230:8082",
		"bc9f49dde58089808ec2837c2efe7310",
		"6nU95IhOWpbymIWvVibiNyfTEblfqkDD0TbXIVEfD30=")
	walletID, err := ifwallet.CreateUserWallet()
	require.Nil(t, err)
	fmt.Println("walletID : ", walletID)

	toWalletID, err := ifwallet.CreateUserWallet()
	require.Nil(t, err)
	fmt.Println("toWalletID : ", toWalletID)

	addr, balance, err := ifwallet.GetAmountOfDenoms(strings.TrimSpace(toWalletID), "spice")
	require.Nil(t, err)
	fmt.Printf("addr : %s, balance : %d \n", addr, balance)
}

func TestIfWallet_GetAmountOfDenoms(t *testing.T) {
	//ad51d596fb56ca992624e182922c627b168974ca578a53adb32a2a6adde1e04d
	//10507cfa86dc0bfdf9e73b727630f86c5fa21f565f18c560dbdfa458e14b4a1a money
	//https://api2.ifwallet.com/h5_open/hd_wallet/slp/balance
	ifwallet := NewWallet("https://api2.ifwallet.com",
		"bb42979c22b37b9ba5f1def23f806051",
		"el5B8CnhPTxA-7mh4f15G_v9uEs6Uw5h3gyJxonnVH0=")
	addr, balance, err := ifwallet.GetAmountOfDenoms(
		"ecf301589599ab040106b1ecf00a8be41995fd01800cc7033f95204e556a09dc", "spice")
	require.Nil(t, err)
	fmt.Printf("addr : %s, balance : %d \n", addr, balance)

}

func TestIfWallet_SendMoney(t *testing.T) {
	ifwallet := NewWallet("http://47.52.172.230:8082",
		"bc9f49dde58089808ec2837c2efe7310",
		"6nU95IhOWpbymIWvVibiNyfTEblfqkDD0TbXIVEfD30=")
	txid, err := ifwallet.SendMoney("10507cfa86dc0bfdf9e73b727630f86c5fa21f565f18c560dbdfa458e14b4a1a", []TransferNews{
		{
			Address: "simpleledger:qrjv94g0hwsdyanwltn8xunky84vcl9d0u2z6et0sl",
			Denom:   "spice",
			Amount:  10,
		},
	})
	require.Nil(t, err)
	fmt.Println(txid)
}

func TestSign(t *testing.T) {
	sign := generateSign("bc9f49dde58089808ec2837c2efe7310",
		"6nU95IhOWpbymIWvVibiNyfTEblfqkDD0TbXIVEfD30=",
		"1579507409")
	fmt.Println(sign)
}
