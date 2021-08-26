package statics

import (
	"github.com/jasonodonnell/go-opendota"
)

type RankScoreSort []opendota.PlayerRankings

func (s RankScoreSort) Len() int           { return len(s) }
func (s RankScoreSort) Less(i, j int) bool { return s[i].Score > s[j].Score }
func (s RankScoreSort) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type HeroGamesSort []opendota.PlayerHero

func (s HeroGamesSort) Len() int           { return len(s) }
func (s HeroGamesSort) Less(i, j int) bool { return s[i].Games > s[j].Games }
func (s HeroGamesSort) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type PeerGamesSort []opendota.PlayerPeers

func (s PeerGamesSort) Len() int           { return len(s) }
func (s PeerGamesSort) Less(i, j int) bool { return s[i].Games > s[j].Games }
func (s PeerGamesSort) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type Match struct {
	Detail *opendota.Match `json:",omitempty"`
}

type Player struct {
	Tire    string                     `json:",omitempty"`
	Info    *opendota.Player           `json:",omitempty"`
	Rank    *[]opendota.PlayerRankings `json:",omitempty"`
	Hero    *[]opendota.PlayerHero     `json:",omitempty"`
	Friends *[]opendota.PlayerPeers    `json:",omitempty"`
	Matches *[]opendota.PlayerMatch    `json:",omitempty"`
	Val     *[]opendota.PlayerTotals   `json:",omitempty"`
	WinLose *opendota.WinLoss          `json:",omitempty"`
}
