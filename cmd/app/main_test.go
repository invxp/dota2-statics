package main

import (
	"bytes"
	"encoding/json"
	"github.com/invxp/dota2-statics/internal/util/ding"
	"github.com/invxp/dota2-statics/internal/util/io"
	"github.com/invxp/dota2-statics/internal/util/log"
	"github.com/invxp/dota2-statics/internal/util/redis"
	"github.com/invxp/dota2-statics/pkg/server"
	"github.com/invxp/dota2-statics/pkg/statics"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func fakeBot(t *testing.T, message string) {
	info := &server.HTTPServerResponse{}
	d := &ding.BotFrom{}
	d.Text.Content = message
	b, _ := json.Marshal(d)

	resp, err := http.Post("http://localhost:7777/bot", "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(b, info)
	if err != nil {
		t.Fatal(err)
	}

	err = ding.SendMarkdown("谢猪猪", info.Markdown.Text, "3d79de450f8a4bc155d9e531507ee112ad2bdb4ea4fc1c832ae2ee5bba93749f")

	if err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	currentPath, currentExecutable := io.CurrentExecutablePath()
	logger := log.New(currentPath, "test", currentExecutable+".log", 66666, 66666)

	server.Start("localhost:7777", logger,
		redis.SimpleClient("localhost:6379", "", 3, logger),
		100,
		statics.New("E09635A9F555CE8F0B0CCEECE8E40434", logger),
		"")

	time.Sleep(time.Second)

	m.Run()
}

func complexStr(n1, n2 string) string {
	var r rune
	var ret string

	if len(n1) < 1 || len(n2) < 1 {
		return ""
	}

	if n1 == "1" {
		return n2
	}

	if n2 == "1" {
		return n1
	}

	if n1 == "0" || n2 == "0" {
		return "0"
	}

	// 找负数
	comp := (n1[0] == '-') == (n2[0] == '-')

	if n1[0] == '-' {
		n1 = n1[1:]
	}
	if n2[0] == '-' {
		n2 = n2[1:]
	}

	tmp := map[int]rune{}
	i := 0

	// 逐位相乘，错位相加
	for _, a := range n1 {
		j := 0
		for _, b := range n2 {
			tmp[i+j] += (a - '0') * (b - '0')
			j++
		}
		i++
	}

	// 倒着连接字符串
	for i := len(tmp) - 1; i >= 0; i-- {
		tmp[i] += r
		r = tmp[i] / 10
		ret = string([]rune{tmp[i]%10 + '0'}) + ret
	}

	if r != 0 {
		ret = string([]rune{r + '0'}) + ret
	}

	if comp {
		return ret
	}
	return "-" + ret
}

func TestBot(t *testing.T) {
	/*
		require.Equal(t, strconv.Itoa(5*6), complexStr("5", "6"))
		require.Equal(t, strconv.Itoa(-15*6), complexStr("-15", "6"))
		require.Equal(t, strconv.Itoa(-15*-3), complexStr("-15", "-3"))
		require.Equal(t, strconv.Itoa(1*1), complexStr("1", "1"))
		require.Equal(t, strconv.Itoa(1*2), complexStr("1", "2"))
		require.Equal(t, strconv.Itoa(1*0), complexStr("1", "0"))
		require.Equal(t, strconv.Itoa(0*1), complexStr("0", "1"))
		require.Equal(t, strconv.Itoa(666*777), complexStr("666", "777"))
	*/

	fakeBot(t, "解绑 爸爸")
	fakeBot(t, "解绑 阿猫")
	fakeBot(t, "绑定 136700549 阿猫")
	fakeBot(t, "绑定 6666 阿猫")
	//fakeBot(t, "玩家 阿猫")
	fakeBot(t, "玩家 136700549")
	fakeBot(t, "队友 136700549")

	fakeBot(t, "玩家 4445")
	fakeBot(t, "比赛 4445")
	fakeBot(t, "比赛 3559037317")

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
