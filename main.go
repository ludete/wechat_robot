package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ludete/wechat_robot/app"
	"github.com/ludete/wechat_robot/util"

	log "github.com/sirupsen/logrus"
)

var (
	cfgPath string
)

func init() {
	newFlag := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	newFlag.StringVar(&cfgPath, "c", "config.toml", "config path")
}

func main() {
	cfg, err := util.LoadConfig(cfgPath)
	if err != nil {
		fmt.Println("load config file failed : ", cfgPath)
		return
	}
	if err := util.InitLog(cfg); err != nil {
		fmt.Println("init util failed ")
		return
	}

	app, err := app.NewRobotApp(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.Start()
	util.WaitForSignal()
	log.Info("robot begin stop")
}
