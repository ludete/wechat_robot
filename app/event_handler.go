package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

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
			handlerGroupInvite(getPrivNews(r), responseWeChat)
		case ReceiveAddFriendRequest:
			handlerFriendVerify(getPrivNews(r), responseWeChat)
		}
	}
}

func handlerPrivateChatMsg(news *privNews, app *RobotApp, fn func([]byte) error) {
	var resMsg []byte
	resMsg = news.groupResMsg(PrivateChatType, getHelpMsg(app))
	if strings.HasPrefix(news.recvMsg, BUYTOKEN) {
		resMsg = buyTokens(app, news)
	} else if strings.HasPrefix(news.recvMsg, QUERY) {
		resMsg = news.groupResMsg(PrivateChatType, queryTokenPrice(app, news.recvMsg))
	}
	if err := Retry(3, 3, func() error {
		return fn(resMsg)
	}); err != nil {
		log.Errorf("response private msg failed : %s\n", err.Error())
		return
	}
}

func handlerReceiveTransfer(news *privNews, app *RobotApp, fn func([]byte) error) {
	resMsg := news.groupResMsg(ResponseTransferType, news.recvMsg)
	if err := Retry(3, 3, func() error {
		return fn(resMsg)
	}); err != nil {
		log.Errorf("response receive transfer failed : %s\n", err.Error())
		return
	}
	if err := app.db.ReceiveRMB(news.sendMsgWeChatID, news.typeKey); err != nil {
		log.Errorf("store amount RMB value in db failed : %s\n", err.Error())
	}
}

func handlerGroupChat(news *GroupMsg, app *RobotApp, fn func([]byte) error) {
	var resMsg []byte
	if strings.HasPrefix(news.revMsg, HELP) {
		// 如果at 了机器人; 进行帮助信息的回复
		if _, ok := news.atWeChatIDS[news.robotID]; ok {
			resMsg = news.GroupMsg(ResGroupChatType, getHelpMsg(app))
		}
	} else if strings.HasPrefix(news.revMsg, TIPS) {
		resMsg = news.GroupMsg(PrivateChatType, "打赏失败")
		// 从发送信息的人的账户， 打赏 at的所有人，一定数量的金额
		datas := strings.Split(news.revMsg, " ")
		amountStr := datas[1]
		denom := datas[2]
		amount, err := strconv.Atoi(amountStr)
		if err == nil {
			if txid, err := tipDenomToPeoples(app, denom, amount, news); err == nil {
				resMsg = news.GroupMsg(PrivateChatType, fmt.Sprintf("txid : %s", txid))
			}
		}
	}

	if err := Retry(3, 3, func() error {
		return fn(resMsg)
	}); err != nil {
		log.Errorf("回复群消息失败")
	}
	return
}

func handlerGroupInvite(news *privNews, fn func([]byte) error) {
	resMsg := news.groupResMsg(AgreeGroupInvite, news.recvMsg)
	if err := Retry(3, 3, func() error {
		return fn(resMsg)
	}); err != nil {
		log.Errorf("response group invite failed : %s\n", err.Error())
	}
}

func handlerFriendVerify(news *privNews, fn func([]byte) error) {
	resMsg := news.groupResMsg(AgreeFriendVerify, news.recvMsg)
	if err := Retry(3, 3, func() error {
		return fn(resMsg)
	}); err != nil {
		log.Errorf("response friend verify failed : %s\n", err.Error())
	}
}
