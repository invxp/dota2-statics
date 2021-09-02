package server

import (
	"fmt"
	"github.com/jasonodonnell/go-opendota"
	"sort"
	"strconv"
)

type HeroWLSort []opendota.PlayerHero

func (s HeroWLSort) Len() int { return len(s) }
func (s HeroWLSort) Less(i, j int) bool {
	return float32(s[i].Win)/float32(s[i].Games)*100.00 > float32(s[j].Win)/float32(s[j].Games)*100.00
}
func (s HeroWLSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s *Server) processBind(id, nickname, history string) string {
	if history != "" {
		return fmt.Sprintf("账号已绑定")
	}
	if err := s.setBinds(id, nickname); err == nil {
		return fmt.Sprintf("绑定: %s-%s 成功", id, nickname)
	} else {
		return fmt.Sprintf("绑定: %s-%s 异常: %v", id, nickname, err)
	}
}

func (s *Server) processUnbind(nickname, history string) string {
	if history == "" {
		return fmt.Sprintf("账号: %s 未绑定", nickname)
	}
	if err := s.delBinds(nickname); err == nil {
		return fmt.Sprintf("解绑: %s-%s 成功", nickname, history)
	} else {
		return fmt.Sprintf("解绑: %s-%s 异常: %v", nickname, history, err)
	}
}

func (s *Server) processPlayer(nickname, accountID string) string {
	if accountID == "" {
		accountID = nickname
	}

	resp, err := s.loadPlayerInfo(accountID)
	if err != nil {
		resp, err = playerInfo("http://"+s.serverAddr, accountID)
		if err != nil {
			return fmt.Sprintf("获取玩家: %s 信息失败: %v", nickname, err)
		}
		_ = s.storePlayerInfo(accountID, resp)
	}

	player := resp.PlayerInfo

	val := *player.Val
	apm := 0
	kda := float32(0)
	gpm := 0
	for _, v := range val {
		if v.Field == "actions_per_min" {
			if v.N == 0 {
				apm = v.Sum
			} else {
				apm = v.Sum / v.N
			}
		}
		if v.Field == "kda" {
			if v.N == 0 {
				kda = float32(v.Sum)
			} else {
				kda = float32(v.Sum) / float32(v.N)
			}
		}

		if v.Field == "gold_per_min" {
			if v.N == 0 {
				gpm = v.Sum
			} else {
				gpm = v.Sum / v.N
			}
		}
	}

	match := *player.Matches

	matchResult := ""
	singleTotal := float32(0)
	partyTotal := float32(0)
	singleWin := float32(0)
	partyWin := float32(0)
	totalWin := float32(0)

	n := 10
	if len(match) < 10 {
		n = len(match)
	}
	var winLose []string
	for i := 0; i < n; i++ {
		_, heroStat := s.statics.HeroIDToName(match[i].HeroID)
		win := (match[i].PlayerSlot <= 127 && match[i].RadiantWin) || (match[i].PlayerSlot > 127 && !match[i].RadiantWin)
		wLose := "<font color=#FF4848 size=3 style=\"font-weight:bold\">负</font>"

		if match[i].PartySize > 1 {
			partyTotal++
		} else {
			singleTotal++
		}

		if win {
			totalWin++
			if match[i].PartySize > 1 {
				partyWin++
			} else {
				singleWin++
			}
			wLose = "<font color=#28DF99 size=3 style=\"font-weight:bold\">胜</font>"
		}

		matchResult += fmt.Sprintf("%s", fmt.Sprintf("![avatar](https://steamcdn-a.akamaihd.net/%s)", heroStat.Icon))
		winLose = append(winLose, wLose+"&nbsp;&nbsp;&nbsp;&nbsp;")
	}

	if singleTotal == 0 {
		singleTotal = 1
	}
	if partyTotal == 0 {
		partyTotal = 1
	}

	info := player.Info
	wl := float32(0)
	if (player.WinLose.Win + player.WinLose.Lose) > 0 {
		wl = float32(player.WinLose.Win) / float32(player.WinLose.Win+player.WinLose.Lose) * 100.00
	}

	mdContent := fmt.Sprintf("![avatar](%s)\n\n<font color=#FF4848 size=2 style=\"font-weight:bold\">玩家信息: %s(%s)-%s\n\n综合能力: %d(整体胜率%.2f%%)\n\nAPM: %d / GPM: %d / KDA: %.2f\n\n%s\n\n%s\n\n%v\n\n</font>",
		info.Profile.AvatarFull,
		info.Profile.Personaname,
		player.Tire,
		accountID,
		info.MmrEstimate.Estimate,
		wl,
		apm, gpm, kda,
		fmt.Sprintf("<font color=#64C9CF size=2 style=\"font-weight:bold\">近%d场胜率: 单排/组排/总(%.2f%%/%.2f%%/%.2f%%)</font>",
			n,
			(singleWin/singleTotal)*100.00,
			(partyWin/partyTotal)*100.00,
			(totalWin)/float32(n)*100.00),
		matchResult,
		winLose)

	rank := *player.Rank

	mdContent += fmt.Sprintf("<font color=#64C9CF size=2 style=\"font-weight:bold\">榜上英雄</font>\n\n")
	n = 10
	if len(rank) < 10 {
		n = len(rank)
	}
	for i := 0; i < n; i++ {
		_, heroStat := s.statics.HeroIDToName(rank[i].HeroID)
		mdContent += fmt.Sprintf("![avatar](https://steamcdn-a.akamaihd.net/%s)", heroStat.Icon)
	}
	mdContent += "\n\n"

	hero := *player.Hero
	mdContent += fmt.Sprintf("<font color=#64C9CF size=2 style=\"font-weight:bold\">最爱英雄</font>\n\n")
	n = 10
	if len(hero) < 10 {
		n = len(hero)
	}
	j := n / 2
	if n >= 10 && j < 10 {
		j = 10
	}
	if n < 10 {
		j = n
	}
	for i := 0; i < n; i++ {
		num, _ := strconv.Atoi(hero[i].HeroID)
		_, heroStat := s.statics.HeroIDToName(num)
		mdContent += fmt.Sprintf("![avatar](https://steamcdn-a.akamaihd.net/%s)", heroStat.Icon)
	}

	var heroes []opendota.PlayerHero
	for i := 0; i < j; i++ {
		heroes = append(heroes, hero[i])
	}

	sort.Sort(HeroWLSort(heroes))
	mdContent += fmt.Sprintf("\n\n<font color=#64C9CF size=2 style=\"font-weight:bold\">绝活英雄</font>\n\n")
	for i := 0; i < len(heroes); i++ {
		num, _ := strconv.Atoi(heroes[i].HeroID)
		_, heroStat := s.statics.HeroIDToName(num)
		mdContent += fmt.Sprintf("![avatar](https://steamcdn-a.akamaihd.net/%s)", heroStat.Icon)
	}

	mdContent += "\n\n"

	friends := *resp.PlayerInfo.Friends

	mdContent += "\n\n<font color=#64C9CF size=2 style=\"font-weight:bold\">最佳队友(总场次/胜率)</font>\n"

	n = 10
	if len(friends) < 10 {
		n = len(friends)
	}

	for i := 0; i < n; i++ {
		mdContent += fmt.Sprintf("\n\n![avatar](%s)<font color=#0F52BA size=2 style=\"font-weight:bold\">%s(%d)%d/%.2f%%\n</font>", friends[i].Avatar, friends[i].Personaname, friends[i].AccountID, friends[i].WithGames, float32(friends[i].WithWin)/float32(friends[i].WithGames)*100.00)
	}

	return mdContent
}

func (s *Server) processMatch(_ string) string {
	return fmt.Sprintf("查看比赛详情功能正在开发中...")
}
