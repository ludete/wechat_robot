package app

import (
	"sync"

	"github.com/ludete/wechat_robot/storage"
)

type RobotApp struct {
	dbMutex sync.Mutex
	db      storage.DB
}

func NewRobotApp(dbConfig string) *RobotApp {

	return &RobotApp{
		dbMutex: sync.Mutex{},
		db:      storage.NewDB(dbConfig),
	}
}
