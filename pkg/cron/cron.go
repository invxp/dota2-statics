package cron

import (
	"fmt"
	"github.com/invxp/dota2-statics/internal/util/log"
	"github.com/robfig/cron"
)

type Cron struct {
	c *cron.Cron
	l *log.Log
}

func New(log *log.Log) *Cron {
	return &Cron{cron.New(), log}
}

func (c *Cron) Add(spec string, cmd func()) {
	if err := c.c.AddFunc(spec, cmd); err != nil {
		c.log("add cron error: %v", err)
	}
}

func (c *Cron) Start() {
	c.c.Start()
}

func (c *Cron) Stop() {
	c.c.Stop()
}

func (c *Cron) log(format string, v ...interface{}) {
	if c.l == nil {
		fmt.Printf(format, v...)
	}else{
		c.l.Printf(format, v...)
	}
}