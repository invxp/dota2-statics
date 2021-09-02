package server

import (
	"fmt"
)

var (
	functionInfo = []string{
		"绑定玩家账号: '绑定 游戏数字ID 昵称'",
		"解绑玩家账号: '解绑 昵称'",
		"查询玩家信息: '玩家 游戏数字ID 或 昵称'",
		//"查询队友信息: '队友 游戏数字ID 或 昵称'",
		"查询比赛信息: '比赛 比赛数字ID 或 昵称'",
	}
)

func showFunctions() string {
	ret := ""
	for _, val := range functionInfo {
		ret += "* " + val + "\n"
	}
	return ret
}

func notFound() string {
	return fmt.Sprintf("[流泪]命令不正确或暂时此功能,现有功能如下\n\n%s", showFunctions())
}

func (s *Server) bind(contents []string) string {
	if len(contents) != 3 {
		return notFound()
	}

	return s.processBind(contents[1], contents[2], s.findBinds(contents[2]))
}

func (s *Server) unbind(contents []string) string {
	if len(contents) != 2 {
		return notFound()
	}

	return s.processUnbind(contents[1], s.findBinds(contents[1]))
}

func (s *Server) player(contents []string) string {
	if len(contents) != 2 {
		return notFound()
	}

	return s.processPlayer(contents[1], s.findBinds(contents[1]))
}

func (s *Server) match(contents []string) string {
	if len(contents) != 2 {
		return notFound()
	}

	return s.processMatch(contents[1])
}

/*
func (s *Server) friend(contents []string) string {
	if len(contents) != 2 {
		return notFound()
	}
	return s.processFriend(contents[1], s.findBinds(contents[1]))
}
*/
