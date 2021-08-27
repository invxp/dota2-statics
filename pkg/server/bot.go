package server

import (
	"github.com/invxp/dota2-statics/internal/util/ding"
	"strings"
)

func (s *Server) ProcessDingMessage(content string) string {
	return s.parseFunctions(strings.Trim(content, " \t\r\n"))
}

func (s *Server) parseFunctions(content string) string {
	for k, f := range s.functions {
		if strings.HasPrefix(content, k) {
			return f(strings.Split(content, " "))
		}
	}
	return notFound()
}

func (s *Server) PublishMessages(content string) {
	tokens, err := s.dingTokens()
	if err != nil {
		s.log("publish ding message error: %v", err)
	}
	for _, token := range tokens {
		go func(token string) {
			s.log("publish ding message: %v, %s", ding.SendMarkdown("谢猪猪", content, token), token)
		}(token)
	}
}

func (s *Server) findBinds(nickname string) string {
	accountID, err := s.binds(nickname)
	s.log("find binds, k: %s, v: %s, %v", nickname, accountID, err)
	return accountID
}
