package bot

import (
	"fmt"
	"github.com/invxp/dota2-statics/internal/util/redis"
)

const (
	redisHashDingTokens = "D_Ding_Tokens"
	redisHashBinds = "D_Bot_Binds"
)

func (b *Bot) dingTokens() ([]string, error) {
	var keys []string
	if b.redis == nil {
		return keys, fmt.Errorf("load ding token redis was not opened")
	}
	return b.redis.HKeys(redisHashDingTokens)
}

func binds(redis *redis.Redis, nickname string) (string, error) {
	if redis == nil {
		return "", fmt.Errorf("load binds redis was not opened")
	}
	return redis.HGet(redisHashBinds, nickname)
}

func setBinds(redis *redis.Redis, id, nickname string) error {
	if redis == nil {
		return fmt.Errorf("save binds redis was not opened")
	}
	_, err := redis.HSet(redisHashBinds, nickname, id)
	return err
}