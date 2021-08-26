package server

import (
	"fmt"
	"github.com/invxp/dota2-statics/internal/util/ding"
	"net/http"
)

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

	response.Markdown.Text = s.processDingMessage(dingContent)
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

func (s *Server) processDingMessage(content string) string {
	go func() {
		tokens, err := s.dingTokens()
		if err != nil {
			s.log("get ding tokens error: %v", err)
			return
		}
		for _, token := range tokens {
			s.log("send ding msg: %s-%v", token, ding.SendMarkdown("谢猪猪", fmt.Sprintf("我要发消息(%s)", content), token))
		}
	}()
	return fmt.Sprintf("[色]爸爸们别着急,正在处理中(%s)", content)
}
