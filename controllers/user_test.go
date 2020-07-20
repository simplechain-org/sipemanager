package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

//用户登录
func TestController_Login(t *testing.T) {
	user := &User{
		Username: "damin",
		Password: "123456",
	}
	data, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
		return
	}
	params := bytes.NewBuffer(data)
	var url string = "http://192.168.3.14:8092" + "/api/v1/user/login"
	request, err := http.NewRequest(http.MethodPost, url, params)
	if err != nil {
		t.Error(err)
		return
	}
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

//用户注册
func TestController_Register(t *testing.T) {
	user := &User{
		Username: "damin",
		Password: "123456",
	}
	data, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
		return
	}
	params := bytes.NewBuffer(data)
	var url string = "http://192.168.3.14:8092" + "/api/v1/user/register"
	request, err := http.NewRequest(http.MethodPost, url, params)
	if err != nil {
		t.Error(err)
		return
	}
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
