package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func request(url string) (*HTTPServerResponse, error) {
	response := &HTTPServerResponse{}

	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(bytes, response)

	if err != nil {
		return response, err
	}

	if response.Code > 0 {
		err = fmt.Errorf("%v", response.Message)
	}

	return response, err
}

func playerInfo(serverAddr, id string) (*HTTPServerResponse, error) {
	return request(fmt.Sprintf("%s/api/player?id=%s", serverAddr, id))
}

func matchInfo(serverAddr, id string) (*HTTPServerResponse, error) {
	return request(fmt.Sprintf("%s/api/match?id=%s", serverAddr, id))
}
