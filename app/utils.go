package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ludete/wechat_robot/wallets"

	"github.com/syndtr/goleveldb/leveldb"

	log "github.com/sirupsen/logrus"
)

func getHelpMsg(app *RobotApp) string {
	helpMsg := `
		--机器人沟通指南--

		打赏 - 给某人打赏(仅群聊有效)
				[语法：=打赏 @某人 666]

		查询 - 获取币种信息； 
				[语法: =spice]

		帮助 - 获取机器人的帮助信息
				[语法：=帮助]

		广告 - [语法：=广告]
				`
	if app != nil {
		helpMsg += app.advert
	}
	return helpMsg
}

func Retry(num int, sleep int, fn func() error) error {
	if err := fn(); err != nil {
		if num--; num > 0 {
			return Retry(num, sleep, fn)
		}
		return err
	}
	return nil
}

func responseWeChat(url string, msg []byte) error {
	res, err := http.Post(url, "application/json; Charset=UTF-8", bytes.NewBuffer(msg))
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

//入群 - 机器人邀请进群
//[语法：进群]

//买币 - 依据当前交易所的价格，购买指定币种(仅私聊有效)；进行买币前：必须先给机器人转账(不可发红包)；
//[语法：买币 bch]

func getOrCreateWallet(app *RobotApp, weChatID string) (string, error) {
	walletID, err := app.db.GetUserWalletKeyID(weChatID)
	if err == nil {
		return walletID, nil
	}
	if err = Retry(3, 3, func() error {
		walletID, err = app.wallet.CreateUserWallet()
		return err
	}); err != nil {
		return "", err
	}
	err = app.db.PutUserWalletKeyID(weChatID, walletID)
	return walletID, err
}

func buyTokens(app *RobotApp, news *privNews) []byte {
	denom := getCoinDenomFromMsg(news.recvMsg)
	if !checkBuyCoins(denom) {
		log.Errorf("不支持购买改币种 : %s\n", news.recvMsg)
		return news.groupResMsg(PrivateChatType, SupportTokens())
	}
	amountRMB, err := app.db.GetUserStoreRMB(news.sendMsgWeChatID)
	if err != nil {
		return news.groupResMsg(PrivateChatType, "未查到用户存储的资金")
	}
	walletID, err := getOrCreateWallet(app, news.sendMsgWeChatID)
	if err != nil {
		return news.groupResMsg(PrivateChatType, "购买失败")
	}
	var price string
	if err := Retry(3, 2, func() error {
		price, err = app.exchange.QueryPrice(denom)
		return err
	}); err != nil {
		log.Errorf("交易所查询币种价格失败; %s\n", err.Error())
		return news.groupResMsg(PrivateChatType, "交易所查询币种价格失败")
	}

	buyTokenAmount := calRMBToTokenAmount(amountRMB, price)
	if checkIsTooSmallToken(buyTokenAmount) {
		log.Errorf("购买的币种 : %s 数量 %d 太少\n", denom, buyTokenAmount)
		return news.groupResMsg(PrivateChatType, "购买的数量太少")
	}

	toAddr, err := app.db.GetUserDenomAddrInWallet(news.sendMsgWeChatID, "", denom)
	if err == leveldb.ErrNotFound {
		_, err = app.wallet.SendMoney(walletID, []wallets.TransferNews{
			{
				Address: toAddr,
				Denom:   denom,
				Amount:  int64(buyTokenAmount),
			},
		})
	}

	if err := Retry(3, 3, func() error {
		//return app.wallet.SendMoney(news.receiveMsgWeChatID, news.sendMsgWeChatID, denom, buyTokenAmount)
		return nil
	}); err != nil {
		log.Errorf("send %s token from %s to %s amount %d failed in wallet\n",
			denom, news.receiveMsgWeChatID, news.sendMsgWeChatID, buyTokenAmount)
		return news.groupResMsg(PrivateChatType, "购买失败")
	}
	app.db.ClearUserStoreRMB(news.sendMsgWeChatID)
	app.db.BuyTokenRecord(news.sendMsgWeChatID, denom, buyTokenAmount)
	return news.groupResMsg(PrivateChatType, fmt.Sprintf("购买%s成功，数量%d\n", denom, buyTokenAmount))
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

func calRMBToTokenAmount(rmb int, price string) int {
	return rmb / 1
}

func checkIsTooSmallToken(tokenAmount int) bool {
	if tokenAmount < 1 {
		return true
	}
	return false
}

func getPriceMsg(denom string, price string) string {
	return denom + " 价格：" + price
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

func getPrivNews(r *http.Request) *privNews {
	log.Info(r.PostForm)
	news := new(privNews)
	news.getNewsFromRequest(r)
	return news
}

func getGroupNews(r *http.Request) *GroupMsg {
	log.Info(r.PostForm)
	news := new(GroupMsg)
	news.getGroupMsg(r)
	return news
}

func getKeysFromRequest(r *http.Request) (int, error) {
	//bz, _ := ioutil.ReadAll(r.Body)
	//log.Info(string(bz))
	err := r.ParseForm()
	if err != nil {
		return -1, err
	}
	typeStr := r.PostForm.Get(TypeKey)
	typeKey, err := strconv.Atoi(typeStr)
	if err != nil {
		return -1, err
	}
	return typeKey, nil
}

func tipDenomToPeoples(app *RobotApp, denom string, amount int, news *GroupMsg) (string, error) {
	sendNews, err := getOrCreateWalletIDAndAddr(app, news.sendMsgWeChatID, denom)
	if err != nil {
		return "", err
	}
	sendWallet := sendNews[0]

	recvAddrs := make([]string, 0, len(news.atWeChatIDS))
	for id := range news.atWeChatIDS {
		atWalletAndAddr, err := getOrCreateWalletIDAndAddr(app, id, denom)
		if err != nil {
			return "", err
		}
		recvAddrs = append(recvAddrs, atWalletAndAddr[1])
		if len(recvAddrs) == len(news.atWeChatIDS) {
			time.Sleep(3 * time.Millisecond)
		}
	}

	transfers := make([]wallets.TransferNews, len(recvAddrs))
	for i, addr := range recvAddrs {
		transfers[i] = wallets.TransferNews{
			Address: addr,
			Amount:  int64(amount),
			Denom:   denom,
		}
	}

	txID, err := app.wallet.SendMoney(sendWallet, transfers)
	if err != nil {
		return "", err
	}
	return txID, nil
}

func getOrCreateWalletIDAndAddr(app *RobotApp, weChatID string, denom string) ([]string, error) {
	walletID, err := getOrCreateWalletForUser(app, weChatID)
	if err != nil {
		return nil, err
	}
	addr, err := getOrCreateDenomAddr(app, walletID, weChatID, denom)
	if err != nil {
		return nil, err
	}
	return []string{
		walletID,
		addr,
	}, nil
}

func getOrCreateWalletForUser(app *RobotApp, weChatID string) (string, error) {
	sendWalletID, err := app.db.GetUserWalletKeyID(weChatID)
	if err != nil {
		if err = Retry(3, 3, func() error {
			sendWalletID, err = app.wallet.CreateUserWallet()
			return err
		}); err != nil {
			log.Errorf("create user wallet failed, err : %s", err.Error())
			return "", err
		}

		if err := app.db.PutUserWalletKeyID(weChatID, sendWalletID); err != nil {
			log.Errorf("store walletID to db failed; err : %s", err.Error())
			return "", err
		}
	}
	return sendWalletID, nil
}

func getOrCreateDenomAddr(app *RobotApp, walletID, weChatID, denom string) (string, error) {
	fmt.Printf("begin  walletID : %s\n", walletID)
	addr, err := app.db.GetUserDenomAddrInWallet(weChatID, walletID, denom)
	if err != nil {
		if err = Retry(3, 3, func() error {
			fmt.Printf("middle walletID : %s\n", walletID)
			addr, _, err = app.wallet.GetAmountOfDenoms(walletID, denom)
			return err
		}); err != nil {
			log.Errorf("create %s addr failed in wallet, err : %s", denom, err.Error())
			return "", err
		}

		if err = app.db.PutUserDenomAddrInWallet(weChatID, walletID, denom, addr); err != nil {
			log.Errorf("store %s denom addr in db failed, err : %s", denom, err.Error())
			return "", err
		}
	}
	return addr, nil
}
