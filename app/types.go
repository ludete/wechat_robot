package app

import (
	"encoding/json"
	"net/http"
	"strconv"
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
)

type baseNews struct {
	recvMsg            string
	typeKey            int
	sendMsgWeChatID    string
	receiveMsgWeChatID string
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
	b.sendMsgWeChatID = r.PostForm.Get(SendMsgIDKey)
	b.receiveMsgWeChatID = r.PostForm.Get(ReceiveMsgIDKey)
	b.recvMsg = r.PostForm.Get(MsgKey)
	return nil
}

func (b *baseNews) groupResMsg(msgType int, resMsg string) ([]byte, error) {
	data := make(map[string]interface{})
	data[TypeKey] = msgType
	data[MsgKey] = resMsg
	data[ReceiveMsgIDKey] = b.receiveMsgWeChatID
	if msgType == 301 {
		data[FriendIDKey] = b.sendMsgWeChatID
	} else {
		data[ResReceiveIDKey] = b.sendMsgWeChatID
	}
	bz, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	//s, err := simplifiedchinese.GBK.NewEncoder().String(msg)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	return bz, nil
}
