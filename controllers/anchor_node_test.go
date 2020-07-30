package controllers

import (
	"fmt"
	"github.com/simplechain-org/go-simplechain/common"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestController_ListAnchorNode(t *testing.T) {
	var url string = "http://localhost:8092" + "/api/v1/anchor/node/list"
	request, err := http.NewRequest(http.MethodGet, url, nil)
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

func TestController_AddAnchorNode(t *testing.T) {
	fmt.Println(common.IsHexAddress("17529b05513e5595ceff7f4fb1e06512c271a5"))
}
