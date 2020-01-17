package util

import (
	"os"
	"os/signal"
	"syscall"

	toml "github.com/pelletier/go-toml"
)

func LoadConfig(file string) (*toml.Tree, error) {
	return toml.LoadFile(file)
}

func WaitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-c
}
