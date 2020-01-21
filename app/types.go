package app

import (
	"encoding/json"
	"net/http"
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
	RobotIDKey      = "robot_wxid"

	ResReceiveIDKey = "to_wxid"
	FriendIDKey     = "friend_wxid"

	GroupMsgSendKey = "final_from_wxid"
	GroupRoomKey    = "from_wxid"
	AtWeChatIDKey   = "at_wxid"
)

type privNews struct {
	typeKey int

	// private msg
	sendMsgWeChatID    string
	receiveMsgWeChatID string
	recvMsg            string
}

func (b *privNews) getNewsFromRequest(r *http.Request) {
	b.sendMsgWeChatID = r.PostForm.Get(SendMsgIDKey)
	b.receiveMsgWeChatID = r.PostForm.Get(ReceiveMsgIDKey)
	b.recvMsg = r.PostForm.Get(MsgKey)
	return
}

func (b *privNews) groupResMsg(msgType int, resMsg string) []byte {
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

type GroupMsg struct {
	// private msg
	groupRoomID     string
	sendMsgWeChatID string
	robotID         string
	revMsg          string
	atWeChatIDS     map[string]struct{}
}

func (g *GroupMsg) getGroupMsg(r *http.Request) {
	g.groupRoomID = r.PostForm.Get(GroupRoomKey)
	g.sendMsgWeChatID = r.PostForm.Get(GroupMsgSendKey)
	g.robotID = r.PostForm.Get(RobotIDKey)
	g.revMsg = r.PostForm.Get(MsgKey)
	g.getAtIDs(r)
}

func (g *GroupMsg) getAtIDs(r *http.Request) {
	g.atWeChatIDS = nil
}

func (g *GroupMsg) GroupMsg(typeKey int, msg string) []byte {
	data := make(map[string]interface{})
	data[TypeKey] = typeKey
	data[MsgKey] = msg
	if typeKey == PrivateChatType {
		data[ResReceiveIDKey] = g.sendMsgWeChatID
	} else {
		data[RobotIDKey] = g.robotID
		data[AtWeChatIDKey] = g.sendMsgWeChatID
		data[ResReceiveIDKey] = g.groupRoomID
	}
	bz, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	bz = []byte(toUnicode(string(bz)))
	return bz
}
