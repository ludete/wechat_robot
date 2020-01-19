package app

import (
	"bytes"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//func groupURLVal(msgType int, msg string, robotID, toAccountID string) *url.Values {
//	v := &url.Values{}
//	v.Set(TypeKey, strconv.Itoa(msgType))
//	v.Set(MsgKey, msg)
//	v.Set(RobotIDKey, robotID)
//	if msgType == 301 {
//		v.Set("friend_wxid", toAccountID)
//	} else {
//		v.Set(ToWeChatIDKey, toAccountID)
//	}
//	return v
//}

func responseWeChat(msg []byte) {
	//	res, err := http.PostForm("http://192.168.1.2:8073/send", *values)
	res, err := http.Post("http://192.168.1.2:8073/send",
		"application/x-www-form-urlencoded; Charset=UTF-8", bytes.NewBuffer(msg))
	if err != nil {
		log.Error("send post request failed ...")
		return
	}
	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("read body failed, error : %s\n", err)
		return
	}
	log.Infof("receive response : %s\n", bz)
}
