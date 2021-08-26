package ding

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/invxp/dota2-statics/internal/util/convert"
	"io/ioutil"
	"net/http"
	"time"
)

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

type BotFrom struct {
	Text Text `json:"text"`
}

type BotTo struct {
	MsgType  string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}

type ServerResponse struct {
	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

func sign(secret string) (string, string) {
	ts := convert.I64toA(time.Now().UnixNano() / int64(time.Millisecond))
	h := hmac.New(sha256.New, convert.StringToByte(secret))
	h.Write(convert.StringToByte(ts + "\n" + secret))
	return ts, base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func SendMarkdown(title, markdown, token string) error {
	b, _ := json.Marshal(&BotTo{MsgType: "markdown", Markdown: Markdown{Title: title, Text: markdown}})

	ts, s := sign("SEC68583bf7d6f538bb0caf0d305690ce32c82495064e803940ee60ea93a13db002")

	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&sign=%s&timestamp=%s", token, s, ts)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))

	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	b, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	dsr := &ServerResponse{}
	err = json.Unmarshal(b, dsr)
	if err != nil {
		return err
	}
	if dsr.ErrorCode != 0 {
		err = fmt.Errorf("%s", dsr.ErrorMsg)
	}
	return err
}
