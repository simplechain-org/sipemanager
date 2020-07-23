package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
	"time"
)

func TestDataBaseAccessObject_CreateRewardConfig(t *testing.T) {
	d, _ := time.ParseDuration("-24h")
	mo:=gorm.Model{
		CreatedAt: time.Now().Add(d),
		UpdatedAt: time.Now().Add(d),
	}
	//SourceChainId   uint   `json:"source_chain_id"`
	//TargetChainId   uint   `json:"target_chain_id"`
	//RegulationCycle uint   `json:"regulation_cycle"` //调控周期
	//SignReward      string `json:"sign_reward"`      //单笔签名奖励
	rewardConfig := &RewardConfig{
		Model:           mo,
		SourceChainId:   1,
		TargetChainId:   2,
		RegulationCycle: 90,
		SignReward:      "1090909345",
	}
	id, err := obj.CreateRewardConfig(rewardConfig)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("reward config id:", id)
}
func TestDataBaseAccessObject_GetRewardConfig(t *testing.T) {
	rewardConfigView,err:=obj.GetRewardConfig(1)
	if err!=nil{
		t.Fatal(err)
		return
	}
	fmt.Printf("rewardConfigView=%+v",rewardConfigView)
}
func TestDataBaseAccessObject_CreateRewardConfig2(t *testing.T) {
	//d, _ := time.ParseDuration("-24h")
	//mo:=gorm.Model{
	//	CreatedAt: time.Now().Add(d),
	//	UpdatedAt: time.Now().Add(d),
	//}
	//SourceChainId   uint   `json:"source_chain_id"`
	//TargetChainId   uint   `json:"target_chain_id"`
	//RegulationCycle uint   `json:"regulation_cycle"` //调控周期
	//SignReward      string `json:"sign_reward"`      //单笔签名奖励
	rewardConfig := &RewardConfig{
		//Model:           mo,
		SourceChainId:   1,
		TargetChainId:   2,
		RegulationCycle: 100,
		SignReward:      "1090909346",
	}
	id, err := obj.CreateRewardConfig(rewardConfig)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("reward config id:", id)
}
func TestDataBaseAccessObject_CreateRewardConfig3(t *testing.T) {
	d, _ := time.ParseDuration("-5h")
	mo:=gorm.Model{
		CreatedAt: time.Now().Add(d),
		UpdatedAt: time.Now().Add(d),
	}
	//SourceChainId   uint   `json:"source_chain_id"`
	//TargetChainId   uint   `json:"target_chain_id"`
	//RegulationCycle uint   `json:"regulation_cycle"` //调控周期
	//SignReward      string `json:"sign_reward"`      //单笔签名奖励
	rewardConfig := &RewardConfig{
		Model:           mo,
		SourceChainId:   1,
		TargetChainId:   2,
		RegulationCycle: 120,
		SignReward:      "1090909348",
	}
	id, err := obj.CreateRewardConfig(rewardConfig)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("reward config id:", id)
}

func TestDataBaseAccessObject_GetLatestRewardConfig(t *testing.T) {
	result,err:=obj.GetLatestRewardConfig(1,2)
	if err!=nil{
		t.Fatal(err)
		return
	}
	t.Log(result)
}

func TestDataBaseAccessObject_CreateRewardConfig4(t *testing.T) {
	d, _ := time.ParseDuration("-24h")
	mo:=gorm.Model{
		CreatedAt: time.Now().Add(d),
		UpdatedAt: time.Now().Add(d),
	}
	//SourceChainId   uint   `json:"source_chain_id"`
	//TargetChainId   uint   `json:"target_chain_id"`
	//RegulationCycle uint   `json:"regulation_cycle"` //调控周期
	//SignReward      string `json:"sign_reward"`      //单笔签名奖励
	rewardConfig := &RewardConfig{
		Model:           mo,
		SourceChainId:   2,
		TargetChainId:   1,
		RegulationCycle: 100,
		SignReward:      "2090909345",
	}
	id, err := obj.CreateRewardConfig(rewardConfig)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("reward config id:", id)
	d, _ = time.ParseDuration("-5h")
	mo=gorm.Model{
		CreatedAt: time.Now().Add(d),
		UpdatedAt: time.Now().Add(d),
	}
	rewardConfig = &RewardConfig{
		Model:           mo,
		SourceChainId:   2,
		TargetChainId:   1,
		RegulationCycle: 120,
		SignReward:      "2090909348",
	}
	id, err = obj.CreateRewardConfig(rewardConfig)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("reward config id:", id)

	rewardConfig = &RewardConfig{
		//Model:           mo,
		SourceChainId:   2,
		TargetChainId:   1,
		RegulationCycle: 150,
		SignReward:      "2090909346",
	}
	id, err = obj.CreateRewardConfig(rewardConfig)
	if err != nil {
		t.Fatal(err)
		return
	}
}
func TestDataBaseAccessObject_GetRewardConfigPage(t *testing.T) {
	result,err:=obj.GetRewardConfigPage(0,10)
	if err!=nil{
		t.Fatal(err)
		return
	}
	for _,o:=range result{
		t.Log(o)
	}
	count,err:=obj.GetRewardConfigCount()
	if err!=nil{
		t.Fatal(err)
		return
	}
	t.Log("count:",count)
}
