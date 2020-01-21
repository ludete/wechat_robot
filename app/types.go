package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	READTIMEOUT  = 3
	WRITETIMEOUT = 3
)

const (
	MsgKey          = "msg"
	TypeKey         = "type"
	SendMsgIDKey    = "from_wxid"
	ReceiveMsgIDKey = "robot_wxid"

	ResSendMsgIDKey = "robot_wxid"
	ResReceiveIDKey = "to_wxid"
	FriendIDKey     = "friend_wxid"
	AtWeChatIDKey   = "at_wxid"
)

type baseNews struct {
	typeKey int

	// private msg
	sendMsgWeChatID    string
	receiveMsgWeChatID string

	// group msg
	atWeChatID string
	groupID    string

	recvMsg string

	money int
}

func (b *baseNews) getNewsFromRequest(r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	typeStr := r.PostForm.Get(TypeKey)
	if b.typeKey, err = strconv.Atoi(typeStr); err != nil {
		return err
	}
	logrus.Info(r.PostForm)
	b.sendMsgWeChatID = r.PostForm.Get(SendMsgIDKey)
	b.receiveMsgWeChatID = r.PostForm.Get(ReceiveMsgIDKey)
	b.recvMsg = r.PostForm.Get(MsgKey)
	b.atWeChatID = r.PostForm.Get(AtWeChatIDKey)
	return nil
}

func (b *baseNews) groupResMsg(msgType int, resMsg string) []byte {
	data := make(map[string]interface{})
	data[TypeKey] = msgType
	data[MsgKey] = resMsg
	data[ReceiveMsgIDKey] = b.receiveMsgWeChatID
	if msgType == ResponseTransferType {
		data[FriendIDKey] = b.sendMsgWeChatID
	} else {
		data[ResReceiveIDKey] = b.sendMsgWeChatID
	}
	bz, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	bz = []byte(toUnicode(string(bz)))
	return bz
}

type stop struct {
	error
}

func NoRetryError(err error) stop {
	return stop{err}
}
