package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sipemanager/dao"
	"testing"
)

func TestController_CreateChain(t *testing.T) {
	chain := &dao.Chain{
		Name:      "simplechain",
		NetworkId: 1,
		CoinName:  "SIPC",
		Symbol:    "SIPC",
	}
	data, err := json.Marshal(chain)
	if err != nil {
		t.Error(err)
		return
	}
	params := bytes.NewBuffer(data)
	var url string = "http://127.0.0.1:8092" + "/api/v1/chain/create"
	request, err := http.NewRequest(http.MethodPost, url, params)
	if err != nil {
		t.Error(err)
		return
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error(err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(string(body))
	}
}
