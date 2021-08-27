package server

import (
	"fmt"
	"strconv"
)

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

func (s *Server) processPlayer(nickname, history string) string {
	if history == "" {
		history = nickname
	}

	resp, err := s.loadPlayerInfo(history)
	if err != nil {
		resp, err = playerInfo("http://"+s.serverAddr, history)
		if err != nil {
			return fmt.Sprintf("获取玩家: %s 信息失败: %v", nickname, err)
		}
		_ = s.storePlayerInfo(history, resp)
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

	info := player.Info
	//Info
	wl := float32(0)
	if (player.WinLose.Win + player.WinLose.Lose) > 0 {
		wl = float32(player.WinLose.Win) / float32(player.WinLose.Win+player.WinLose.Lose) * 100.00
	}

	mdContent := fmt.Sprintf("玩家: **%s(%s)-%s(%s)**\n\n![avatar](%s)\n\n综合能力: **%d(整体胜率%.2f%%)**\n\n平均APM: **%d**\n\n平均GPM: **%d**\n\n平均KDA: **%.2f**\n\n", info.Profile.Personaname, player.Tire, nickname, history, info.Profile.AvatarFull, info.MmrEstimate.Estimate, wl, apm, gpm, kda)

	rank := *player.Rank
	//Rank
	mdContent += fmt.Sprintf("**绝活英雄**\n\n")
	n := 10
	if len(rank) < 10 {
		n = len(rank)
	}
	for i := 0; i < n; i++ {
		_, heroStat := s.statics.HeroIDToName(rank[i].HeroID)
		mdContent += fmt.Sprintf("![avatar](https://steamcdn-a.akamaihd.net/%s)", heroStat.Icon)
	}
	mdContent += "\n\n"

	hero := *player.Hero
	mdContent += fmt.Sprintf("**最爱英雄**\n\n")
	n = 10
	if len(hero) < 10 {
		n = len(hero)
	}
	for i := 0; i < n; i++ {
		num, _ := strconv.Atoi(hero[i].HeroID)
		_, heroStat := s.statics.HeroIDToName(num)
		mdContent += fmt.Sprintf("![avatar](https://steamcdn-a.akamaihd.net/%s)", heroStat.Icon)
	}
	mdContent += "\n\n"

	match := *player.Matches

	matchResult := ""
	singleTotal := float32(0)
	partyTotal := float32(0)
	singleWin := float32(0)
	partyWin := float32(0)
	totalWin := float32(0)
	n = 10
	if len(match) < 10 {
		n = len(match)
	}
	for i := 0; i < n; i++ {
		heroName, heroStat := s.statics.HeroIDToName(match[i].HeroID)
		win := (match[i].PlayerSlot <= 127 && match[i].RadiantWin) || (match[i].PlayerSlot > 127 && !match[i].RadiantWin)
		winStr := "负"

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
			winStr = "胜"
		}

		matchResult += fmt.Sprintf("\n\n%s-%s(%s-%d分钟)\n- KDA:%d/%d/%d\n", fmt.Sprintf("![avatar](https://steamcdn-a.akamaihd.net/%s)", heroStat.Icon), heroName, winStr, match[i].Duration/60, match[i].Kills, match[i].Deaths, match[i].Assists)
	}

	mdContent += fmt.Sprintf("最近%d场胜率: **单排/组排/总(%.2f%%/%.2f%%/%.2f%%)**", n, (singleWin/singleTotal)*100.00, (partyWin/partyTotal)*100.00, (totalWin)/float32(n)*100.00)
	mdContent += matchResult

	friends := *player.Friends

	mdContent += "\n\n最佳队友\n"

	n = 10
	if len(friends) < 10 {
		n = len(friends)
	}

	for i := 0; i < n; i++ {
		mdContent += fmt.Sprintf("\n\n![avatar](%s)%s(%d),胜率: %.2f%%\n", friends[i].Avatar, friends[i].Personaname, friends[i].AccountID, float32(friends[i].WithWin)/float32(friends[i].WithGames)*100.00)
	}

	//s.PublishMessages(mdContent)

	return mdContent
}

func (s *Server) processMatch(_ string) string {
	return fmt.Sprintf("查看比赛详情功能正在开发中...")
}

func (s *Server) processStatics(_ string) string {
	return fmt.Sprintf("查看统计详情功能正在开发中...")
}
