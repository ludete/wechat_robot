package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoLevelDB_PutUserWalletKeyID(t *testing.T) {
	db := NewDB("data")
	defer os.RemoveAll("data")
	weChatID := "hello"
	walletID := "nihao"
	err := db.PutUserWalletKeyID(weChatID, walletID)
	require.Nil(t, err)
	tmpWalletID, err := db.GetUserWalletKeyID(weChatID)
	require.Nil(t, err)
	require.EqualValues(t, walletID, tmpWalletID)

}

func TestGoLevelDB_PutUserDenomAddrInWallet(t *testing.T) {
	db := NewDB("data")
	defer os.RemoveAll("data")
	denom := "spice"
	weChatID := "nihao"
	addr := "djkwu"
	walletID := "sdwee"

	err := db.PutUserDenomAddrInWallet(weChatID, walletID, denom, addr)
	require.Nil(t, err)
	tmpAddr, err := db.GetUserDenomAddrInWallet(weChatID, walletID, denom)
	require.Nil(t, err)
	require.EqualValues(t, addr, tmpAddr)
}
