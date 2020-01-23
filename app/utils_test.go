package app

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestTrim(t *testing.T) {
	msg := "=帮助"
	data := strings.Trim(msg, msg)
	fmt.Println(data)
}
