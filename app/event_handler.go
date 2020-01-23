package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ludete/wechat_robot/exchanges"

	log "github.com/sirupsen/logrus"
)

type ResponseFunc func(string, []byte) error

func handler(app *RobotApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key, err := getKeysFromRequest(r)
		if err != nil {
			log.Errorf(err.Error())
			return
		}

		switch key {
		case PrivateChatType:
			handlerPrivateChatMsg(w, getPrivNews(r), app, responseWeChat)
		case ReceiveTransferType:
			handlerReceiveTransfer(w, getPrivNews(r), app, responseWeChat)
		case GroupChatType:
			handlerGroupChat(w, getGroupNews(r), app, responseWeChat)
		case AgreeGroupInvite:
			handlerGroupInvite(w, getPrivNews(r), app, responseWeChat)
		case ReceiveAddFriendRequest:
			handlerFriendVerify(w, getPrivNews(r), app, responseWeChat)
		}
	}
}

func getBalanceAndAddr(app *RobotApp, news AssemblyMsg) []byte {
	walletID, err := getOrCreateWalletForUser(app, news.getSendWeChatID())
	if err != nil {
		return news.groupResMsg(PrivateChatType, "未查询到余额信息")
	}
	addr, amount, err := app.wallet.GetAmountOfDenoms(walletID, exchanges.SPICE)
	if err != nil {
		return news.groupResMsg(PrivateChatType, "未查询到余额信息")
	}
	return news.groupResMsg(PrivateChatType, fmt.Sprintf("余额 : %d;\n 地址：%s\n", amount, addr))
}

func handlerPrivateChatMsg(w http.ResponseWriter, news *privNews, app *RobotApp, fn ResponseFunc) {
	var resMsg []byte
	if strings.HasPrefix(news.recvMsg, BUYTOKEN) {
		resMsg = buyTokens(app, news)
	} else if strings.HasPrefix(news.recvMsg, HELP) {
		resMsg = news.groupResMsg(PrivateChatType, getHelpMsg(app))
	} else if _, ok := app.coins[strings.ToLower(news.recvMsg)]; ok {
		if price, err := app.exchange.QueryPrice(news.recvMsg); err == nil {
			resMsg = news.groupResMsg(PrivateChatType, price)
		}
	} else if strings.HasPrefix(news.recvMsg, ADVERT) {
		resMsg = getAdvert(app, news)
	} else if strings.HasPrefix(news.recvMsg, BALANCE) {
		resMsg = getBalanceAndAddr(app, news)
	}
	if len(resMsg) > 0 {
		if err := Retry(3, 3, func() error {
			_, errT := w.Write(resMsg)
			return errT
		}); err != nil {
			log.Errorf("response private msg failed : %s\n", err.Error())
			return
		}
	}
}

func handlerReceiveTransfer(w http.ResponseWriter, news *privNews, app *RobotApp, fn ResponseFunc) {
	resMsg := news.groupResMsg(ResponseTransferType, news.recvMsg)
	if err := Retry(3, 3, func() error {
		_, err := w.Write(resMsg)
		return err
	}); err != nil {
		log.Errorf("response receive transfer failed : %s\n", err.Error())
		return
	}
	if err := app.db.ReceiveRMB(news.sendMsgWeChatID, news.typeKey); err != nil {
		log.Errorf("store amount RMB value in db failed : %s\n", err.Error())
	}
}

// 1. 帮助
// 2. 打赏
func handlerGroupChat(w http.ResponseWriter, news *GroupMsg, app *RobotApp, fn ResponseFunc) {
	var resMsg []byte
	if strings.HasPrefix(news.revMsg, HELP) {
		if _, ok := news.atWeChatIDS[news.robotID]; ok {
			resMsg = news.groupResMsg(ResGroupChatType, getHelpMsg(app))
		}
	} else if strings.HasPrefix(news.revMsg, TIPS) {
		resMsg = tipToken(app, news)
	} else if _, ok := app.coins[strings.ToLower(news.revMsg)]; ok {
		resMsg = queryPrice(app, news)
	} else if strings.HasPrefix(news.revMsg, ADVERT) {
		resMsg = getAdvert(app, news)
	} else if strings.HasPrefix(news.revMsg, BALANCE) {
		if _, ok := news.atWeChatIDS[news.robotID]; ok {
			resMsg = getBalanceAndAddr(app, news)
		}
	}
	if len(resMsg) > 0 {
		retryReq(func() error {
			length, err := w.Write(resMsg)
			if length != len(resMsg) {
				log.Errorf("write response to client failed; the length is not match, expect : %d, actual : %d", len(resMsg), length)
			}
			return err
		})
	}
	return
}

func getAdvert(app *RobotApp, news AssemblyMsg) []byte {
	num, err := strconv.Atoi(strings.TrimPrefix(news.getMsg(), ADVERT))
	if err != nil {
		return nil
	}
	if num > len(app.advert) || num <= 0 {
		return nil
	}
	return news.groupResMsg(ResGroupChatType, app.advert[num-1])
}

func tipToken(app *RobotApp, news *GroupMsg) []byte {
	var resMsg []byte
	log.Info(news)
	resMsg = news.groupResMsg(ResGroupChatType, getHelpMsg(app))
	// 从发送信息的人的账户， 打赏 at的所有人，一定数量的金额
	msg := strings.TrimSpace(strings.Trim(news.revMsg, TIPS))
	if amount, err := strconv.Atoi(msg); err == nil {
		if txid, err := tipDenomToPeoples(app, exchanges.SPICE, amount, news); err == nil {
			resMsg = news.groupResMsg(PrivateChatType, fmt.Sprintf("txid : %s", txid))
		} else {
			resMsg = news.groupResMsg(PrivateChatType, err.Error())
		}
	} else {
		resMsg = news.groupResMsg(ResGroupChatType, "格式错误")
	}
	return resMsg
}

func queryPrice(app *RobotApp, news *GroupMsg) []byte {
	price, err := app.exchange.QueryPrice(news.revMsg)
	if err == nil {
		return news.groupResMsg(ResGroupChatType, price)
	}
	return nil
}

func retryReq(fn func() error) {
	if err := Retry(3, 3, fn); err != nil {
		log.Errorf("回复群消息失败")
	}
}

func handlerGroupInvite(w http.ResponseWriter, news *privNews, app *RobotApp, fn ResponseFunc) {
	resMsg := news.groupResMsg(AgreeGroupInvite, news.recvMsg)
	if err := Retry(3, 3, func() error {
		_, err := w.Write(resMsg)
		return err
	}); err != nil {
		log.Errorf("response group invite failed : %s\n", err.Error())
	}
}

func handlerFriendVerify(w http.ResponseWriter, news *privNews, app *RobotApp, fn ResponseFunc) {
	resMsg := news.groupResMsg(AgreeFriendVerify, news.recvMsg)
	if err := Retry(3, 3, func() error {
		_, err := w.Write(resMsg)
		return err
	}); err != nil {
		log.Errorf("response friend verify failed : %s\n", err.Error())
	}
}
