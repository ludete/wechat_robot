package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMsg(t *testing.T) {
	app := RobotApp{advert: "广告"}
	fmt.Println(getHelpMsg(&app))
}

func TestRetry(t *testing.T) {
	runNum := 0
	Retry(3, 1, func() error {
		runNum++
		return fmt.Errorf("error")
	})
	require.EqualValues(t, 3, runNum)

	runNum = 0
	Retry(3, 1, func() error {
		runNum++
		return nil
	})
	require.EqualValues(t, 1, runNum)
}
