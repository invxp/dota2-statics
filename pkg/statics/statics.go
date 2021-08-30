package statics

import (
	"fmt"
	"github.com/Katsusan/go-dota2"
	"github.com/invxp/dota2-statics/internal/util/convert"
	"github.com/invxp/dota2-statics/internal/util/log"
	"github.com/jasonodonnell/go-opendota"
	"net/http"
	"sort"
	"sync"
)

const (
	SteamID64 = 76561197960265728
)

type D2Statics struct {
	odClient *opendota.Client
	stClient *dota2.Dota2api
	logger   *log.Log
	heroName map[int]string
	heroStat map[int]opendota.HeroStat
}

func New(apiKey string, log *log.Log) *D2Statics {
	d2s := &D2Statics{opendota.NewClient(&http.Client{}), dota2.NewApi(&http.Client{}), log, parseHeroes(), make(map[int]opendota.HeroStat)}
	d2s.stClient.SetApiKey(apiKey)
	d2s.parseHeroStat()
	return d2s
}

func (d *D2Statics) Match(matchID string) (Match, error) {
	match := Match{}
	result, _, err := d.odClient.MatchService.Match(convert.AtoI64(matchID))
	if err != nil {
		return match, err
	}

	if result.MatchID == 0 {
		return match, fmt.Errorf("match id: %s not found", matchID)
	}

	match.Detail = &result
	return match, nil
}

func (d *D2Statics) Player(accountID string) (Player, error) {
	var errors error
	player := Player{}
	_, _, _, d2ID := d.convertSteamIds(convert.AtoI64(accountID))
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		info, _, err := d.odClient.PlayerService.Player(d2ID)
		if err != nil {
			errors = err
			return
		}
		if info.Profile.AccountID == 0 {
			errors = fmt.Errorf("玩家基础信息未找到, ID: %s", accountID)
			return
		}
		player.Info = &info
		player.Tire = rank(info.RankTier)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rk, _, _ := d.odClient.PlayerService.Rankings(d2ID)
		if len(rk) == 0 {
			errors = fmt.Errorf("玩家排名信息未找到, ID: %s", accountID)
			return
		}
		sort.Sort(RankScoreSort(rk))
		player.Rank = &rk
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		he, _, err := d.odClient.PlayerService.Heroes(d2ID, nil)
		if err != nil {
			errors = err
			return
		}
		if len(he) == 0 {
			errors = fmt.Errorf("玩家英雄信息未找到, ID: %s", accountID)
			return
		}
		sort.Sort(HeroGamesSort(he))
		player.Hero = &he
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p, _, err := d.odClient.PlayerService.Peers(d2ID, nil)
		if err != nil {
			errors = err
			return
		}
		if len(p) == 0 {
			errors = fmt.Errorf("玩家英雄信息未找到, ID: %s", accountID)
			return
		}

		sort.Sort(PeerGamesSort(p))
		player.Friends = &p
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		m, _, err := d.odClient.PlayerService.RecentMatches(d2ID)
		if err != nil {
			errors = err
			return
		}
		if len(m) == 0 {
			errors = fmt.Errorf("玩家比赛信息未找到, ID: %s", accountID)
			return
		}
		player.Matches = &m
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, _, _ := d.odClient.PlayerService.Totals(d2ID, nil)
		if len(val) == 0 {
			errors = fmt.Errorf("玩家统计信息未找到, ID: %s", accountID)
			return
		}
		player.Val = &val
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		wl, _, err := d.odClient.PlayerService.WinLoss(d2ID, nil)
		if err != nil {
			errors = err
			return
		}
		player.WinLose = &wl
	}()

	wg.Wait()
	return player, errors
}

/*
func (d *D2Statics) PlayerSummary(steamID string) (dota2.PlayerSummary, error) {
	stID, _, _, _ := d.convertSteamIds(convert.AtoI64(steamID))
	summary, err := d.stClient.GetPlayerSummaries(stID)
	if err != nil {
		return dota2.PlayerSummary{}, err
	}
	if len(summary.PlayerSummary) == 0 {
		return dota2.PlayerSummary{}, fmt.Errorf("player id: %s summary not found", steamID)
	}
	return summary.PlayerSummary[0], nil
}
*/

func (d *D2Statics) log(format string, v ...interface{}) {
	if d.logger == nil {
		fmt.Printf(format, v...)
	} else {
		d.logger.Printf(format, v...)
	}
}

func (d *D2Statics) convertSteamIds(id int64) (steamID string, D2ID string, nSteamID int64, nD2ID int64) {
	nD2ID, nSteamID = id, id

	if nD2ID > SteamID64 {
		nD2ID -= SteamID64
	}

	if nSteamID < SteamID64 {
		nSteamID += SteamID64
	}

	return convert.I64toA(nSteamID), convert.I64toA(nD2ID), nSteamID, nD2ID
}

func (d *D2Statics) HeroIDToName(id int) (string, opendota.HeroStat) {
	return d.heroName[id], d.heroStat[id]
}
