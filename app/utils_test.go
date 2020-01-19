package app

import (
	"fmt"
	"testing"
)

func TestGetMsg(t *testing.T) {
	app := RobotApp{advert: "广告"}
	fmt.Println(getHelpMsg(&app))
}
