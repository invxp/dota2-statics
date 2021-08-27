package server

import (
	"encoding/json"
	"fmt"
	"github.com/invxp/dota2-statics/internal/util/convert"
)

const (
	redisHashDingTokens = "D_Ding_Tokens"
	redisHashBinds      = "D_Bot_Binds"
	redisKeyPlayerInfo  = "D_Player_Info"
)

func (s *Server) storePlayerInfo(id string, value *HTTPServerResponse) error {
	if s.redis == nil {
		return fmt.Errorf("store player info: %s redis was not opened", id)
	}
	b, _ := json.Marshal(value)
	_, err := s.redis.SetEx(redisKeyPlayerInfo+"_"+id, b, 60*60)
	return err
}

func (s *Server) loadPlayerInfo(id string) (*HTTPServerResponse, error) {
	if s.redis == nil {
		return nil, fmt.Errorf("get player info: %s redis was not opened", id)
	}

	ret, err := s.redis.Get(redisKeyPlayerInfo + "_" + id)
	if err != nil {
		return nil, err
	}

	info := &HTTPServerResponse{}
	err = json.Unmarshal(convert.StringToByte(ret), info)
	return info, err
}

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

func (s *Server) binds(nickname string) (string, error) {
	if s.redis == nil {
		return "", fmt.Errorf("load binds redis was not opened")
	}
	return s.redis.HGet(redisHashBinds, nickname)
}

func (s *Server) setBinds(id, nickname string) error {
	if s.redis == nil {
		return fmt.Errorf("save binds redis was not opened")
	}
	_, err := s.redis.HSet(redisHashBinds, nickname, id)
	return err
}

func (s *Server) delBinds(nickname string) error {
	if s.redis == nil {
		return fmt.Errorf("del binds redis was not opened")
	}
	_, err := s.redis.HDel(redisHashBinds, nickname)
	return err
}
