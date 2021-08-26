package bot

import (
	"fmt"
	"github.com/invxp/dota2-statics/internal/util/log"
	"github.com/invxp/dota2-statics/internal/util/redis"
)

var (
	function = map[string]func(*redis.Redis, *log.Log, []string) string {
		"绑定": bind,
		"玩家": player,
		"比赛": match,
	}

	functionInfo = []string {
		"绑定人物账号: '绑定 游戏数字ID 昵称'",
		"查询玩家信息: '玩家 游戏数字ID 或 昵称'",
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

func bind(redis *redis.Redis, log *log.Log, contents []string) string {
	if len(contents) != 3 {
		return notFound()
	}

	if err := setBinds(redis, contents[1], contents[2]); err == nil {
		return fmt.Sprintf("绑定: %s-%s 成功", contents[1], contents[2])
	}else{
		return fmt.Sprintf("绑定: %s-%s 异常: %v", contents[1], contents[2], err)
	}
}

func player(redis *redis.Redis, log *log.Log, contents []string) string {
	if len(contents) != 2 {
		return notFound()
	}

	if id, err := binds(redis, contents[1]); err == nil {
		if id != "" {
			//TODO
			return fmt.Sprintf("获取: %s 成功: %s", contents[1], id)
		}
		//TODO
		id = contents[1]
		return fmt.Sprintf("获取: %s 失败", contents[1])

	}else{
		//TODO
		id = contents[1]
		return fmt.Sprintf("获取: %s 异常: %v", contents[1], err)
	}
}

func match(redis *redis.Redis, log *log.Log, contents []string) string {
	if len(contents) != 2 {
		return notFound()
	}
	if id, err := binds(redis, contents[1]); err == nil {
		if id != "" {
			//TODO
			return fmt.Sprintf("获取: %s 成功: %s", contents[1], id)
		}
		//TODO
		id = contents[1]
		return fmt.Sprintf("获取: %s 失败", contents[1])

	}else{
		//TODO
		id = contents[1]
		return fmt.Sprintf("获取: %s 异常: %v", contents[1], err)
	}
}
