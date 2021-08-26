package server

import (
	"encoding/json"
	"fmt"
	"github.com/invxp/dota2-statics/internal/util/convert"
	"github.com/invxp/dota2-statics/internal/util/ding"
	"github.com/invxp/dota2-statics/internal/util/system"
	"github.com/invxp/dota2-statics/pkg/statics"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	HTTPClientError = 400
	HTTPServerError = 500
)

type HTTPServerResponse struct {
	Code       uint
	Message    string `json:",omitempty"`
	Hostname   string
	Time       string
	ClientIP   string
	MsgType    string          `json:"msgtype,omitempty"`
	Markdown   *ding.Markdown  `json:"markdown,omitempty"`
	MatchInfo  *statics.Match  `json:",omitempty"`
	PlayerInfo *statics.Player `json:",omitempty"`
}

func buildHTTPServerResponse() *HTTPServerResponse {
	return &HTTPServerResponse{Code: 0, Message: "", Hostname: system.Hostname(), Time: time.Now().Format("2006-01-02 15:04:05")}
}

func writeHTTPResponse(writer http.ResponseWriter, response interface{}) {
	_, _ = io.WriteString(writer, convert.MustMarshal(response))
}

func (s *Server) authHTTPClientRequest(request *http.Request, response *HTTPServerResponse) string {
	response.ClientIP = strings.Split(request.RemoteAddr, ":")[0]
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		s.log("client: %s request: %s, header: %v read body error: %v", response.ClientIP, request.URL.RequestURI(), request.Header, err)
		response.Code = HTTPClientError
		response.Message = fmt.Sprintf("client: %s request: %s, header: %v read body error: %v", response.ClientIP, request.URL.RequestURI(), request.Header, err)
		return ""
	}

	if s.tokenBucket != nil && s.tokenBucket.Take() == 0 {
		s.log("client: %s request: %s, header: %v rate limited", response.ClientIP, request.URL.RequestURI(), request.Header)
		response.Code = HTTPClientError
		response.Message = fmt.Sprintf("client: %s request: %s, header: %v rate limited", response.ClientIP, request.URL.RequestURI(), request.Header)
		return ""
	}

	dbf := &ding.BotFrom{}
	_ = json.Unmarshal(bytes, dbf)

	s.log("client: %s request: %s, header: %v success, %s", response.ClientIP, request.URL.RequestURI(), request.Header,  convert.ByteToString(bytes))

	return dbf.Text.Content
}

func (s *Server) handlePlayer(writer http.ResponseWriter, request *http.Request) {
	response := buildHTTPServerResponse()
	defer writeHTTPResponse(writer, response)
	s.authHTTPClientRequest(request, response)
	if response.Code > 0 {
		return
	}

	player, err := s.statics.Player(request.URL.Query().Get("id"))
	if err != nil {
		response.Code = HTTPServerError
		response.Message = fmt.Sprintf("%v", err)
		return
	}
	response.PlayerInfo = &player
}

func (s *Server) handleMatch(writer http.ResponseWriter, request *http.Request) {
	response := buildHTTPServerResponse()
	defer writeHTTPResponse(writer, response)
	s.authHTTPClientRequest(request, response)
	if response.Code > 0 {
		return
	}

	match, err := s.statics.Match(request.URL.Query().Get("id"))
	if err != nil {
		response.Code = HTTPServerError
		response.Message = fmt.Sprintf("%v", err)
		return
	}
	response.MatchInfo = &match
}

func (s *Server) handleBot(writer http.ResponseWriter, request *http.Request) {
	response := buildHTTPServerResponse()
	defer writeHTTPResponse(writer, response)
	dingContent := s.authHTTPClientRequest(request, response)
	if response.Code > 0 {
		return
	}

	token := request.Header.Get("Token")

	response.MsgType = "markdown"
	response.Markdown = &ding.Markdown{Title: "谢猪猪", Text: ""}

	if token != s.dingAuth {
		response.Markdown.Text = "没有访问权限[嘿嘿]"
		return
	}

	response.Markdown.Text = s.bot.ProcessDingMessage(dingContent)
}

func (s *Server) handleAddDingToken(writer http.ResponseWriter, request *http.Request) {
	response := buildHTTPServerResponse()
	defer writeHTTPResponse(writer, response)
	s.authHTTPClientRequest(request, response)
	if response.Code > 0 {
		return
	}

	err := s.storeDingTokens(request.URL.Query().Get("id"))
	if err != nil {
		response.Code = HTTPServerError
		response.Message = fmt.Sprintf("%v", err)
	}
}
