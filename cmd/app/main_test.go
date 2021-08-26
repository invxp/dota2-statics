package main

import (
	"encoding/json"
	"github.com/invxp/dota2-statics/internal/util/redis"
	"github.com/invxp/dota2-statics/pkg/bot"
	"github.com/invxp/dota2-statics/pkg/server"
	"github.com/invxp/dota2-statics/pkg/statics"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	server.Start("localhost:7777", nil,
		redis.SimpleClient("localhost:6379", "", 3, nil),
		100,
		statics.New("E09635A9F555CE8F0B0CCEECE8E40434", nil),
		"")

	time.Sleep(time.Second)

	m.Run()
}

func TestBot(t *testing.T) {
	b := bot.New(redis.SimpleClient("localhost:6379", "", 3, nil), nil)
	t.Log(b.ProcessDingMessage("  绑定 1267736 阿布"))
	t.Log(b.ProcessDingMessage("  哈哈 1267736 阿布"))
}

func TestPlayer(t *testing.T) {
	player := &server.HTTPServerResponse{}

	resp, err := http.Get("http://localhost:7777/api/player?id=136700549")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(bytes, player)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMatch(t *testing.T) {
	match := &server.HTTPServerResponse{}

	resp, err := http.Get("http://localhost:7777/api/match?id=3559037317")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(bytes, match)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStoreToken(t *testing.T) {
	token := &server.HTTPServerResponse{}
	resp, err := http.Get("http://localhost:7777/ctrl/addDingToken?id=3d79de450f8a4bc155d9e531507ee112ad2bdb4ea4fc1c832ae2ee5bba93749f")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(bytes, token)
	if err != nil {
		t.Fatal(err)
	}
}
