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
			handlerPrivateChatMsg(getPrivNews(r), app, responseWeChat)
		case ReceiveTransferType:
			handlerReceiveTransfer(getPrivNews(r), app, responseWeChat)
		case GroupChatType:
			handlerGroupChat(getGroupNews(r), app, responseWeChat)
		case AgreeGroupInvite:
			handlerGroupInvite(getPrivNews(r), app, responseWeChat)
		case ReceiveAddFriendRequest:
			handlerFriendVerify(getPrivNews(r), app, responseWeChat)
		}
	}
}

func handlerPrivateChatMsg(news *privNews, app *RobotApp, fn ResponseFunc) {
	var resMsg []byte
	if strings.HasPrefix(news.recvMsg, BUYTOKEN) {
		resMsg = buyTokens(app, news)
	} else if strings.HasPrefix(news.recvMsg, HELP) {
		resMsg = news.groupResMsg(PrivateChatType, getHelpMsg(app))
	} else {
		if price, err := app.exchange.QueryPrice(news.recvMsg); err == nil {
			resMsg = news.groupResMsg(PrivateChatType, price)
		}
	}
	if err := Retry(3, 3, func() error {
		return fn(app.resURL, resMsg)
	}); err != nil {
		log.Errorf("response private msg failed : %s\n", err.Error())
		return
	}
}

func handlerReceiveTransfer(news *privNews, app *RobotApp, fn ResponseFunc) {
	resMsg := news.groupResMsg(ResponseTransferType, news.recvMsg)
	if err := Retry(3, 3, func() error {
		return fn(app.resURL, resMsg)
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
func handlerGroupChat(news *GroupMsg, app *RobotApp, fn ResponseFunc) {
	var resMsg []byte
	if strings.HasPrefix(news.revMsg, HELP) {
		// 如果at 了机器人; 进行帮助信息的回复
		if _, ok := news.atWeChatIDS[news.robotID]; ok {
			resMsg = news.groupResMsg(ResGroupChatType, getHelpMsg(app))
		}
	} else if strings.HasPrefix(news.revMsg, TIPS) {
		resMsg = news.groupResMsg(ResGroupChatType, getHelpMsg(app))
		// 从发送信息的人的账户， 打赏 at的所有人，一定数量的金额
		datas := strings.Split(news.revMsg, TIPS)
		if len(datas) == 2 {
			if amount, err := strconv.Atoi(strings.TrimSpace(datas[1])); err == nil {
				if txid, err := tipDenomToPeoples(app, exchanges.SPICE, amount, news); err == nil {
					resMsg = news.groupResMsg(PrivateChatType, fmt.Sprintf("txid : %s", txid))
				}
			}
		} else {
			resMsg = news.groupResMsg(ResGroupChatType, "格式错误")
		}
	} else if _, ok := app.coins[strings.ToLower(news.revMsg)]; ok {
		price, err := app.exchange.QueryPrice(news.revMsg)
		if err == nil {
			resMsg = news.groupResMsg(ResGroupChatType, price)
		}
	}
	if len(resMsg) != 0 {
		log.Info("response wechat msg : ", string(resMsg))
		if err := Retry(3, 3, func() error {
			return fn(app.resURL, resMsg)
		}); err != nil {
			log.Errorf("回复群消息失败")
		}
	}
	return
}

func handlerGroupInvite(news *privNews, app *RobotApp, fn ResponseFunc) {
	resMsg := news.groupResMsg(AgreeGroupInvite, news.recvMsg)
	if err := Retry(3, 3, func() error {
		return fn(app.resURL, resMsg)
	}); err != nil {
		log.Errorf("response group invite failed : %s\n", err.Error())
	}
}

func handlerFriendVerify(news *privNews, app *RobotApp, fn ResponseFunc) {
	resMsg := news.groupResMsg(AgreeFriendVerify, news.recvMsg)
	if err := Retry(3, 3, func() error {
		return fn(app.resURL, resMsg)
	}); err != nil {
		log.Errorf("response friend verify failed : %s\n", err.Error())
	}
}
