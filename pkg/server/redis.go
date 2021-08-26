package server

import "fmt"

const (
	redisHashDingTokens = "D_Ding_Tokens"
)

func (s *Server) storeDingTokens(token string) error {
	if s.redis == nil {
		return fmt.Errorf("store ding token: %s redis was not opened", token)
	}
	_, err := s.redis.HSet(redisHashDingTokens, token, token)
	return err
}

func (s *Server) dingTokens() ([]string, error) {
	var keys []string
	if s.redis == nil {
		return keys, fmt.Errorf("load ding token redis was not opened")
	}
	return s.redis.HKeys(redisHashDingTokens)
}
