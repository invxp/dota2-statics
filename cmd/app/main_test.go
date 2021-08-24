package main

import (
	"fmt"
	"github.com/invxp/dota2-statics/pkg/statics"
	"github.com/jasonodonnell/go-opendota"
	"net/http"
	"testing"
)

func TestD2(t *testing.T) {
	s := statics.New("KEY")
	s.MatchInfo("88888")
	client := opendota.NewClient(&http.Client{})

	fmt.Println(client.ProPlayerService.Players())

	p, r, e := client.PlayerService.Player(136700549)

	fmt.Println(p, r, e)
}
