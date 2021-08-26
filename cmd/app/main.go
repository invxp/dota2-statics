package main

import (
	"flag"
	"github.com/invxp/dota2-statics/internal/util/config"
	"github.com/invxp/dota2-statics/internal/util/deamon"
	"github.com/invxp/dota2-statics/internal/util/io"
	"github.com/invxp/dota2-statics/internal/util/log"
	"github.com/invxp/dota2-statics/internal/util/redis"
	"github.com/invxp/dota2-statics/pkg/cron"
	"github.com/invxp/dota2-statics/pkg/server"
	"github.com/invxp/dota2-statics/pkg/statics"
	"os"
	"os/signal"
	"syscall"
)

const (
	flagConfig    = "c"
	flagDaemon    = "d"
	flagDaemonize = "daemon"
)

var (
	configFile   = flag.String(flagConfig, "service.toml", "set a config file")
	enableDaemon = flag.Bool(flagDaemon, false, "run program in daemonize")
	_            = flag.Int(flagDaemonize, 0, "daemonize pid flag")
)

const (
	appVersion = "0.0.1-alpha"
)

func main() {
	flag.Parse()

	currentPath, currentExecutable := io.CurrentExecutablePath()

	daemon.Daemonize(*enableDaemon, flagDaemon, flagDaemonize, currentPath)

	conf := config.Load(currentPath, *configFile)

	logFile := log.New(currentPath, conf.Log.Path, currentExecutable+".log", conf.Log.MaxAge, conf.Log.MaxRotationSize)

	logFile.Printf("application started, version: %s, key: %s", appVersion, conf.Server.APIKey)

	server.Start(conf.Server.Address, logFile,
		redis.SimpleClient(conf.Redis.Address, conf.Redis.Password, conf.Redis.Database, logFile),
		conf.Server.Rate,
		statics.New(conf.Server.APIKey, logFile),
		conf.Server.DingAuth)

	c := cron.New(logFile)
	c.Add("0 */10 * * * ?", func() {})
	c.Start()
	defer c.Stop()

	quitSig := make(chan os.Signal)
	signal.Notify(quitSig, syscall.SIGINT, syscall.SIGTERM)
	<-quitSig
	logFile.Printf("application exit")
}
