package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ludete/wechat_robot/app"
	"github.com/ludete/wechat_robot/util"
	toml "github.com/pelletier/go-toml"

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
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		fmt.Println("load config file failed : ", cfgPath)
		return
	}
	if err := util.InitLog(cfg); err != nil {
		fmt.Println("init util failed ")
		return
	}
	app := app.NewRobotApp(cfg)
	app.Start()
	waitForSignal()
	log.Info("robot begin stop")
}

func loadConfig(file string) (*toml.Tree, error) {
	return toml.LoadFile(file)
}

func waitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-c
}
