package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"

	log "github.com/sirupsen/logrus"
)

func Retry(num int, sleep int, fn func() error) error {
	if err := fn(); err != nil {
		if num--; num > 0 {
			return Retry(num, sleep, fn)
		}
		return err
	}
	return nil
}

func responseWeChat(msg []byte) error {
	//	res, err := http.PostForm("http://192.168.1.2:8073/send", *values)
	res, err := http.Post("http://192.168.1.2:8073/send",
		"application/json; Charset=UTF-8", bytes.NewBuffer(msg))
	if err != nil {
		log.Error("send post request failed ...")
		return err
	}
	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("read body failed, error : %s\n", err)
		return nil
	}
	log.Infof("receive response : %s\n", bz)
	return nil
}

func queryTokenPrice(app *RobotApp, msg string) string {
	if denom := getCoinDenomFromMsg(msg); denom != "" {
		if price, err := app.exchange.QueryPrice(denom); err == nil {
			return getPriceMsg(denom, price)
		}
	}
	return toUnicode("查询价格失败")
}

func getHelpMsg(app *RobotApp) string {
	helpMsg := `		--机器人沟通指南--
		查询 - 获取币种信息； 
				[语法: 币种 bch]

		买币 - 依据当前交易所的价格，购买指定币种(仅私聊有效)；进行买币前：必须先给机器人转账(不可发红包)；
				[语法：买币 bch]

		打赏 - 给某人打赏(仅群聊有效)
				[语法：打赏 1cet @某人]

		入群 - 机器人邀请进群
				[语法：进群]

		帮助 - 获取机器人的帮助信息
				[语法：帮助]
				`
	if app != nil {
		helpMsg += app.advert
	}

	//return url.QueryEscape(helpMsg)
	//return toUnicode(helpMsg)
	return helpMsg
}

func buyTokens(app *RobotApp, news *baseNews) []byte {
	var (
		denom = getCoinDenomFromMsg(news.recvMsg)
		price int
	)

	amountRMB, err := app.db.GetUserStoreRMB(news.sendMsgWeChatID)
	if err != nil {
		return news.groupResMsg(PrivateChatType, "未查到用户存储的资金")
	}
	if !checkBuyCoins(denom) {
		log.Errorf("不支持购买改币种 : %s\n", news.recvMsg)
		return news.groupResMsg(PrivateChatType, SupportTokens())
	}

	if err := Retry(3, 2, func() error {
		price, err = app.exchange.QueryPrice(denom)
		return err
	}); err != nil {
		log.Errorf("交易所查询币种价格失败; %s\n", err.Error())
		return news.groupResMsg(PrivateChatType, "交易所查询币种价格失败")
	}

	tokenAmount := calRMBToTokenAmount(amountRMB, price)
	if checkIsTooSmallToken(tokenAmount) {
		log.Errorf("购买的币种 : %s 数量 %d 太少\n", denom, tokenAmount)
		return news.groupResMsg(PrivateChatType, "购买的数量太少")
	}

	walletID := app.db.GetUserWalletKeyID(news.receiveMsgWeChatID)
	toAddr, err := app.db.GetUserDeomAddr(news.sendMsgWeChatID, denom)
	if err == leveldb.ErrNotFound {
		err = app.wallet.SendMoney(walletID, toAddr, denom, tokenAmount)
	}

	if err := Retry(3, 3, func() error {
		return app.wallet.SendMoney(news.receiveMsgWeChatID, news.sendMsgWeChatID, denom, tokenAmount)
	}); err != nil {
		log.Errorf("send %s token from %s to %s amount %d failed in wallet\n",
			denom, news.receiveMsgWeChatID, news.sendMsgWeChatID, tokenAmount)
		return news.groupResMsg(PrivateChatType, "购买失败")
	}
	app.db.ClearUserStoreRMB(news.sendMsgWeChatID)
	app.db.BuyTokenRecord(news.sendMsgWeChatID, denom, tokenAmount)
	return news.groupResMsg(PrivateChatType, fmt.Sprintf("购买%s成功，数量%d\n", denom, tokenAmount))
}

func SupportTokens() string {
	support := "支持购买的币种 : " + BCH + "、" + CET
	return support
}

func getCoinDenomFromMsg(msg string) string {
	if infos := strings.Split(msg, ""); len(infos) == 2 {
		return infos[1]
	}
	return ""
}

func calRMBToTokenAmount(rmb int, price int) int {
	return rmb / price
}

func checkIsTooSmallToken(tokenAmount int) bool {
	if tokenAmount < 1 {
		return true
	}
	return false
}

func getPriceMsg(denom string, price int) string {
	return toUnicode(denom + " 价格：" + strconv.Itoa(price))
}

func toUnicode(str string) string {
	runes := []rune(str)
	res := ""
	for _, r := range runes {
		if r < rune(128) {
			res += string(r)
		} else {
			res += "\\u" + strconv.FormatInt(int64(r), 16)
		}
	}
	return res
}

func checkBuyCoins(denom string) bool {
	switch denom {
	case BCH, CET:
		return true
	}
	return false
}
