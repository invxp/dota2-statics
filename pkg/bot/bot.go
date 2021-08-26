package bot

import (
	"fmt"
	"github.com/invxp/dota2-statics/internal/util/ding"
	"github.com/invxp/dota2-statics/internal/util/log"
	"github.com/invxp/dota2-statics/internal/util/redis"
	"strings"
)

type Bot struct {
	redis  *redis.Redis
	logger *log.Log
}

func New(redis *redis.Redis, log *log.Log) *Bot {
	return &Bot{redis, log}
}

func (b *Bot) ProcessDingMessage(content string) string {
	return b.parseFunctions(strings.TrimLeft(content, " "))
}

func (b *Bot) parseFunctions(content string) string {
	for k, f := range function {
		if strings.HasPrefix(content, k) {
			return f(b.redis, b.logger, strings.Split(content, " "))
		}
	}
	return notFound()
}

func (b *Bot) publishMessages(content string) {
	tokens, err := b.dingTokens()
	if err != nil {
		b.log("publish ding message error: %v", err)
	}
	for _, token := range tokens {
		go func(token string) {
			b.log("publish ding message: %v, %s", ding.SendMarkdown("谢猪猪", content, token), token)
		}(token)
	}
}

func (b *Bot) log(format string, v ...interface{}) {
	if b.logger == nil {
		fmt.Printf(format, v...)
	} else {
		b.logger.Printf(format, v...)
	}
}