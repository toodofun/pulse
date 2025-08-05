package main

import (
	"github.com/sirupsen/logrus"
	"pulse/internal/config"
	"pulse/internal/server"
)

func main() {
	svc, err := server.New(config.New("config.yaml"))
	if err != nil {
		panic(err)
	}

	logrus.Infof("running with config: %+v", config.Current())

	if err = svc.Run(); err != nil {
		panic(err)
	}
}
