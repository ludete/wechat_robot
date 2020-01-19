package app

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"

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

func getHelpMsg(app *RobotApp) string {
	helpMsg := `		--机器人沟通指南--
		查询 - 获取币种信息； 
			语法: 币种 bch
		买币 - 依据当前交易所的价格，购买指定币种(仅私聊有效)；
			进行买币前：必须先给机器人转账(不可发红包)；
			语法：买币 bch
		打赏 - 给某人打赏(仅群聊有效)
			语法：打赏 1cet @某人
		入群 - 机器人邀请进群
			语法：进群
		帮助 - 获取机器人的帮助信息
			语法：帮助
				`
	if app != nil {
		helpMsg += app.advert
	}
	return toUnicode(helpMsg)
}

func getPriceMsg(denom string, price int) string {
	return toUnicode(denom + " 价格：" + strconv.Itoa(price))
}

func toUnicode(str string) string {
	runes := []rune(str)
	res := ""
	for _, r := range runes {
		if r < rune(128) {
			res += string(r)
		} else {
			res += "\\u" + strconv.FormatInt(int64(r), 16)
		}
	}
	return res
}
