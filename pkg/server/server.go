package server

import (
	"fmt"
	"github.com/invxp/dota2-statics/internal/util/log"
	"github.com/invxp/dota2-statics/internal/util/redis"
	"github.com/invxp/dota2-statics/pkg/bot"
	"github.com/invxp/dota2-statics/pkg/statics"
	"github.com/invxp/tokenbucket"
	"net/http"
	"time"
)

const (
	Player       = "/api/player"
	Match        = "/api/match"
	Bot          = "/bot"
	AddDingToken = "/ctrl/addDingToken"
)

type Server struct {
	logger      *log.Log
	redis       *redis.Redis
	tokenBucket *tokenbucket.TokenBucket
	statics     *statics.D2Statics
	dingAuth    string
	bot         *bot.Bot
}

func Start(serverAddr string, log *log.Log, db *redis.Redis, rate uint32, statics *statics.D2Statics, dingAuth string) *Server {
	srv := &Server{log, db, tokenbucket.NewTokenBucket(time.Second, rate), statics, dingAuth, bot.New(db, log)}

	http.HandleFunc(Player, srv.handlePlayer)
	http.HandleFunc(Match, srv.handleMatch)
	http.HandleFunc(Bot, srv.handleBot)
	http.HandleFunc(AddDingToken, srv.handleAddDingToken)

	go func() {
		log.Printf("http server stop: %v", http.ListenAndServe(serverAddr, nil))
	}()

	return srv
}

func (s *Server) log(format string, v ...interface{}) {
	if s.logger == nil {
		fmt.Printf(format, v...)
	} else {
		s.logger.Printf(format, v...)
	}
}