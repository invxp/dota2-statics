package server

import (
	"github.com/fvbock/endless"
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
}

func Start(serverAddr string, log *log.Log, db *redis.Redis, rate uint32, statics *statics.D2Statics, dingAuth string) *Server {
	srv := &Server{log, db, tokenbucket.NewTokenBucket(time.Second, rate), statics, dingAuth}

	http.HandleFunc(Player, srv.handlePlayer)
	http.HandleFunc(Match, srv.handleMatch)
	http.HandleFunc(Bot, srv.handleBot)
	http.HandleFunc(AddDingToken, srv.handleAddDingToken)

	go func() {
		httpServer := endless.NewServer(serverAddr, nil)
		log.Printf("http server stop: %v", httpServer.ListenAndServe())
	}()

	return srv
}
