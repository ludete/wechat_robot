package app

import (
	"encoding/json"
	"net/http"
)

func responseData(w http.ResponseWriter, msgType int, msg string, robotID string, toWeChatID string) error {
	data := make(map[string]interface{})
	data[TypeKey] = msgType
	data[MsgKey] = msg
	data[RobotIDKey] = robotID
	data[ToWeChatIDKey] = toWeChatID
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	return err
}
