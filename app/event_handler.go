package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func receiveRMB(app *RobotApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bz, err := ioutil.ReadAll(r.Body)
		if err != nil {
			responseData(w, PrivateChatType, "读取请求出错", "", "")
			return
		}

		data := make(map[string]interface{})
		err = json.Unmarshal(bz, data)
		if err != nil {
			responseData(w, PrivateChatType, "解析请求出错", "", "")
			return
		}

		robotID := data[RobotIDKey].(string)
		toChatID := data[ToWeChatIDKey].(string)
		msg := data[MsgKey].(string)

		//TODO. will store account state and amount to store

		responseData(w, ReceiveTransferType, msg, robotID, toChatID)
	}
}

func buyToken(app *RobotApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func withdrawToken(app *RobotApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func queryBalance(app *RobotApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func wechatReward(app *RobotApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func help(app *RobotApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func advert(app *RobotApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
