package app

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func handler(app *RobotApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		news := new(baseNews)
		if err := news.getNewsFromRequest(r); err != nil {
			return
		}
		log.Infof("typeKey : %d, sendMsgWeChatID : %s, receiveMsgWeChatID : %s, recvMsg : %s\n",
			news.typeKey, news.sendMsgWeChatID, news.receiveMsgWeChatID, news.recvMsg)

		switch news.typeKey {
		case PrivateChatType:
			handlerPrivateChatMsg(news, app, responseWeChat)
		case ReceiveTransferType:
			handlerReceiveTransfer(news, app, responseWeChat)
		case GroupChatType:
			handlerGroupChat(news, app)
		case AgreeGroupInvite:
			handlerGroupInvite(news, responseWeChat)
		case ReceiveAddFriendRequest:
			handlerFriendVerify(news, responseWeChat)
		}
	}
}

func handlerPrivateChatMsg(news *baseNews, app *RobotApp, fn func([]byte) error) {
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

func handlerReceiveTransfer(news *baseNews, app *RobotApp, fn func([]byte) error) {
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

func handlerGroupChat(news *baseNews, app *RobotApp) {
	//var resMsg []byte
}

func handlerGroupInvite(news *baseNews, fn func([]byte) error) {
	resMsg := news.groupResMsg(AgreeGroupInvite, news.recvMsg)
	if err := Retry(3, 3, func() error {
		return fn(resMsg)
	}); err != nil {
		log.Errorf("response group invite failed : %s\n", err.Error())
	}
}

func handlerFriendVerify(news *baseNews, fn func([]byte) error) {
	resMsg := news.groupResMsg(AgreeFriendVerify, news.recvMsg)
	if err := Retry(3, 3, func() error {
		return fn(resMsg)
	}); err != nil {
		log.Errorf("response friend verify failed : %s\n", err.Error())
	}
}
