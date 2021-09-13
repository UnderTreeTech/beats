package main

import (
	"flag"
	"fmt"
	"github.com/UnderTreeTech/waterdrop/pkg/log"
	"os"
	"os/signal"
	"syscall"

	"beats/internal/dao"
	"beats/internal/service"

	"github.com/UnderTreeTech/waterdrop/pkg/stats"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"
)

// run: go run main.go -conf=../configs/application.toml
func main() {
	flag.Parse()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	conf.Init()

	logCfg := &log.Config{}
	if err := conf.Unmarshal("log", logCfg); err != nil {
		panic(fmt.Sprintf("parse log config fail, err msg %s", err.Error()))
	}
	defer log.New(logCfg).Sync()

	dao := dao.New()
	s := service.New(dao)
	if _, err := stats.StartStats(); err != nil {
		panic(fmt.Sprintf("start stats fail, err msg is %s", err.Error()))
	}

	<-c

	s.Close()
}
