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
			handlerPrivateChatMsg(news, app)
		case ReceiveTransferType:
			handlerReceiveTransfer(news)
		case GroupChatType:
			handlerGroupChat(news, app)
		case AgreeGroupInvite:
			handlerGroupInvite(news)
		case ReceiveAddFriendRequest:
			handlerFriendVerify(news)
		}
	}
}

func handlerPrivateChatMsg(news *baseNews, app *RobotApp) {
	var resMsg []byte
	resMsg = news.groupResMsg(PrivateChatType, getHelpMsg(app))
	if strings.HasPrefix(news.recvMsg, BUYTOKEN) {
		resMsg = news.groupResMsg(PrivateChatType, "buy token is not implement !")
	} else if strings.HasPrefix(news.recvMsg, QUERY) {
		if infos := strings.Split(news.recvMsg, QUERY); len(infos) == 2 {
			if price, err := app.exchange.QueryPrice(infos[1]); err == nil {
				resMsg = []byte(getPriceMsg(infos[1], price))
			}
		}
	}
	responseWeChat(resMsg)
	//var resMsg *url.Values
	//responseURLVal(resMsg)
}

func handlerReceiveTransfer(news *baseNews) {
	resMsg := news.groupResMsg(ResponseTransferType, news.recvMsg)
	responseWeChat(resMsg)
}

func handlerGroupChat(news *baseNews, app *RobotApp) {
	//var resMsg []byte
}

func handlerGroupInvite(news *baseNews) {
	resMsg := news.groupResMsg(AgreeGroupInvite, news.recvMsg)
	responseWeChat(resMsg)
}

func handlerFriendVerify(news *baseNews) {
	resMsg := news.groupResMsg(AgreeFriendVerify, news.recvMsg)
	responseWeChat(resMsg)
}
