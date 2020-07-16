package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestController_AddWallet(t *testing.T) {

	walletParam := &WalletParam{
		Name:     "daminyang",
		Content:  `{"address":"1de1d9a800ce214a807b6cfc819935144b51c4d7","crypto":{"cipher":"aes-128-ctr","ciphertext":"fbcdc0ed90ae8630abdf7e22b1574655a966180c60d1dd37f3edba217d5b6083","cipherparams":{"iv":"822d5e5000f035439443bc6ede90f0fd"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"eb31b6e42067f9c395194b6006ef3faa1b62324facef8c65c8bcd683ed07b450"},"mac":"ec06a8f5cd9834f7aebeb04dde76bf1559931ff3fad3c9f89f9c5d6661921e6b"},"id":"0f845b23-083c-42e3-bcf1-53111c384f17","version":3}`,
		Password: "123456",
	}
	data, err := json.Marshal(walletParam)
	if err != nil {
		t.Error(err)
		return
	}
	params := bytes.NewBuffer(data)
	var url string = "http://127.0.0.1:8092" + "/api/v1/wallet"
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
func TestController_AddWallet2(t *testing.T) {
	walletParam := &WalletParam{
		Name:     "daminyang",
		Content:  `2d0b168d0c8a4b1d9f86d061a5bb4b18a390eb66fbc7c26c55aa4fdb11a2e981`,
		Password: "123456",
	}
	data, err := json.Marshal(walletParam)
	if err != nil {
		t.Error(err)
		return
	}
	params := bytes.NewBuffer(data)
	var url string = "http://127.0.0.1:8092" + "/api/v1/wallet"
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
func TestController_AddWallet3(t *testing.T) {
	walletParam := &WalletParam{
		Name:     "daminyang",
		Content:  `royal food modify legal duck venture size trigger hospital brave typical kid season define pattern session arctic upon mention nature route volume thunder pulse`,
		Password: "123456",
	}
	data, err := json.Marshal(walletParam)
	if err != nil {
		t.Error(err)
		return
	}
	params := bytes.NewBuffer(data)
	var url string = "http://127.0.0.1:8092" + "/api/v1/wallet"
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

//没有加载keystore文件，所以content为空
func TestController_ListWallet(t *testing.T) {
	var url string = "http://127.0.0.1:8092" + "/api/v1/wallet/list"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Error(err)
		return
	}
	request.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error(err)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(body))
}

func TestController_UpdateWallet(t *testing.T) {
	paramObj := &UpdateWalletParam{
		WalletId:    8,
		NewPassword: "12345678",
		OldPassword: "123456",
	}
	data, err := json.Marshal(paramObj)
	if err != nil {
		t.Error(err)
		return
	}
	params := bytes.NewBuffer(data)
	var url string = "http://127.0.0.1:8092" + "/api/v1/wallet/update"
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
func TestController_UpdateWallet2(t *testing.T) {
	paramObj := &UpdateWalletParam{
		WalletId:    7,
		NewPassword: "12345678",
		OldPassword: "12345645",
	}
	data, err := json.Marshal(paramObj)
	if err != nil {
		t.Error(err)
		return
	}
	params := bytes.NewBuffer(data)
	var url string = "http://127.0.0.1:8092" + "/api/v1/wallet/update"
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
