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
		log.Infof("typeKey : %s, sendMsgWeChatID : %s, receiveMsgWeChatID : %s, recvMsg : %s\n",
			news.typeKey, news.sendMsgWeChatID, news.receiveMsgWeChatID, news.recvMsg)

		switch news.typeKey {
		case PrivateChatType:
			handlerPrivateChatMsg(news)
		case ReceiveTransferType:
			handlerReceiveTransfer(news)
		case GroupChatType:
			handlerGroupChat(news)
		case AgreeGroupInvite:
			handlerGroupInvite(news)
		case ReceiveAddFriendRequest:
			handlerFriendVerify(news)
		}
	}
}

func handlerPrivateChatMsg(news *baseNews) {
	var (
		resMsg []byte
		err    error
	)
	//var resMsg *url.Values
	if news.recvMsg == HELP || news.recvMsg == "help" {
		resMsg, err = news.groupResMsg(PrivateChatType, "help me !")
	} else if strings.HasPrefix(news.recvMsg, BUYTOKEN) {
		resMsg, err = news.groupResMsg(PrivateChatType, "buy token is not implement !")
	} else {
		resMsg, err = news.groupResMsg(PrivateChatType, "pls input help me !")
	}
	if err != nil {
		log.Errorf("group response message error : %s\n", err.Error())
		return
	}
	responseWeChat(resMsg)
	//responseURLVal(resMsg)
}

func handlerReceiveTransfer(news *baseNews) {
	resMsg, err := news.groupResMsg(ResponseTransferType, news.recvMsg)
	if err != nil {
		log.Errorf("group response message error : %s\n", err.Error())
		return
	}
	responseWeChat(resMsg)
}

func handlerGroupChat(news *baseNews) {

}

func handlerGroupInvite(news *baseNews) {

}

func handlerFriendVerify(news *baseNews) {
	resMsg, err := news.groupResMsg(AgreeFriendVerify, news.recvMsg)
	if err != nil {
		log.Errorf("group response message error : %s\n", err.Error())
		return
	}
	responseWeChat(resMsg)
}
