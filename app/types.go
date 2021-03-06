package app

import (
	"encoding/json"
	"net/http"
	"strings"

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
	RobotIDKey      = "robot_wxid"

	ResReceiveIDKey = "to_wxid"
	FriendIDKey     = "friend_wxid"

	GroupMsgSendKey  = "final_from_wxid"
	GroupFromName    = "from_name"
	GroupRoomKey     = "from_wxid"
	AtWeChatIDKey    = "at_wxid"
	AtWeChatNickName = "at_name"
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
	sendMsgNickName string
	robotID         string
	revMsg          string
	atWeChatIDS     map[string]string
}

func (g *GroupMsg) getGroupMsg(r *http.Request) {
	g.groupRoomID = r.PostForm.Get(GroupRoomKey)
	g.sendMsgWeChatID = r.PostForm.Get(GroupMsgSendKey)
	g.sendMsgNickName = r.PostForm.Get(GroupFromName)
	g.robotID = r.PostForm.Get(RobotIDKey)
	g.revMsg = getRealMsg(r.PostForm.Get(MsgKey))
	logrus.Infof("%s\n", r.PostForm.Get(MsgKey))
	g.atWeChatIDS = getAtWeChatMsgs(r.PostForm.Get(MsgKey))
}

func (g *GroupMsg) groupResMsg(typeKey int, msg string) []byte {
	data := make(map[string]interface{})
	data[TypeKey] = typeKey
	data[MsgKey] = msg
	data[RobotIDKey] = g.robotID
	if typeKey == PrivateChatType {
		data[ResReceiveIDKey] = g.sendMsgWeChatID
	} else {
		data[AtWeChatIDKey] = g.sendMsgWeChatID
		data[ResReceiveIDKey] = g.groupRoomID
		data[AtWeChatNickName] = g.sendMsgNickName
	}
	bz, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	bz = []byte(toUnicode(string(bz)))
	return bz
}

func getRealMsg(msg string) string {
	if !strings.Contains(msg, "[@at,") {
		index := strings.Index(msg, "[@at,")
		if index < 0 {
			return strings.TrimSpace(msg)
		}
		return strings.TrimSpace(msg[:index])
	}
	msg = trimAtWeChatMsg(strings.TrimSpace(msg))
	return getRealMsg(msg)
}

func trimAtWeChatMsg(msg string) string {
	begin := strings.Index(msg, "[@at,")
	end := strings.Index(msg, "]")
	return strings.TrimSpace(msg[:begin] + msg[end+1:])
}

func getAtWeChatMsgs(msg string) map[string]string {
	data := make(map[string]string)
	for {
		atMsg := getAtWeChatMsg(msg)
		if len(atMsg) == 0 {
			return data
		}
		data[getAtID(atMsg)] = getNickName(atMsg)
		index := strings.Index(msg, atMsg)
		msg = msg[index+len(atMsg):]
	}
}

func getAtWeChatMsg(msg string) string {
	begin := strings.Index(msg, "[@at,nickname=")
	end := strings.Index(msg, "]")
	if begin < 0 || end < 0 {
		return ""
	}
	return msg[begin : end+1]
}

func getNickName(msg string) string {
	begin := strings.Index(msg, "nickname") + 9
	end := strings.Index(msg, ",wxid")
	return msg[begin:end]
}

func getAtID(msg string) string {
	begin := strings.Index(msg, "wxid") + 5
	end := strings.Index(msg, "]")
	return msg[begin:end]
}
