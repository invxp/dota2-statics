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
