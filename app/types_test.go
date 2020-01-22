package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRealMsg(t *testing.T) {
	msg := "帮助 [@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12]"
	msg = getRealMsg(msg)
	require.EqualValues(t, "帮助", msg)

	msg = "[@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12] 帮助"
	msg = getRealMsg(msg)
	require.EqualValues(t, "帮助", msg)

	msg = "[@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12] [@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12] 帮助[@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12]"
	msg = getRealMsg(msg)
	require.EqualValues(t, "帮助", msg)

	msg = "[@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12] [@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12] [@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12] 帮助 haha"
	msg = getRealMsg(msg)
	require.EqualValues(t, "帮助 haha", msg)

	msg = "[@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12] 帮助 [@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12] [@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12] 帮助 haha"
	msg = getRealMsg(msg)
	require.EqualValues(t, "帮助   帮助 haha", msg)
}

func TestGetAtMsg(t *testing.T) {
	msg := "帮助 [@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12]"
	data := getAtWeChatMsgs(msg)
	require.EqualValues(t, "数字货币机器人", data["wxid_xno0ahdy95zg12"])

	msg = "[@at,nickname=数字货币机器人2,wxid=wxid_xno0ahdy95zg122] 帮助 [@at,nickname=数字货币机器人,wxid=wxid_xno0ahdy95zg12]"
	data = getAtWeChatMsgs(msg)
	require.EqualValues(t, 2, len(data))
	require.EqualValues(t, "数字货币机器人", data["wxid_xno0ahdy95zg12"])
	require.EqualValues(t, "数字货币机器人2", data["wxid_xno0ahdy95zg122"])

	msg = "[@at,nickname=数字货币机器人1,wxid=wxid_xno0ahdy95zg121] " +
		"帮助 " +
		"[@at,nickname=数字货币机器人2,wxid=wxid_xno0ahdy95zg122] haha " +
		"[@at,nickname=数字货币机器人3,wxid=wxid_xno0ahdy95zg1233] 帮助 haha"
	data = getAtWeChatMsgs(msg)
	require.EqualValues(t, 3, len(data))
	require.EqualValues(t, "数字货币机器人1", data["wxid_xno0ahdy95zg121"])
	require.EqualValues(t, "数字货币机器人2", data["wxid_xno0ahdy95zg122"])
	require.EqualValues(t, "数字货币机器人3", data["wxid_xno0ahdy95zg1233"])
}
