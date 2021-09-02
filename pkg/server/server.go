package server

import (
	"fmt"
	"github.com/invxp/dota2-statics/internal/util/log"
	"github.com/invxp/dota2-statics/internal/util/redis"
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
	serverAddr  string
	functions   map[string]func([]string) string
}

func Start(serverAddr string, log *log.Log, db *redis.Redis, rate uint32, statics *statics.D2Statics, dingAuth string) *Server {
	srv := &Server{log, db, tokenbucket.NewTokenBucket(time.Second, rate), statics, dingAuth, serverAddr, make(map[string]func([]string) string)}

	srv.functions["绑定"] = srv.bind
	srv.functions["解绑"] = srv.unbind
	srv.functions["玩家"] = srv.player
	srv.functions["比赛"] = srv.match
	//srv.functions["队友"] = srv.friend

	http.HandleFunc(Player, srv.handlePlayer)
	http.HandleFunc(Match, srv.handleMatch)
	http.HandleFunc(Bot, srv.handleBot)
	http.HandleFunc(AddDingToken, srv.handleAddDingToken)

	go func() {
		log.Printf("http server stop: %v", http.ListenAndServe(srv.serverAddr, nil))
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
