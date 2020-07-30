package dao

import (
	"testing"
)

func TestDataBaseAccessObject_CreateWorkCount(t *testing.T) {
	//SourceChainId       uint `gorm:"source_chain_id" json:"source_chain_id"`
	//TargetChainId       uint `gorm:"target_chain_id" json:"target_chain_id"`
	//AnchorNodeId        uint `gorm:"anchor_node_id" json:"anchor_node_id"` //锚定节点编号
	//SignCount           uint `gorm:"sign_count" json:"sign_count"`
	//FinishCount         uint `gorm:"finish_count" json:"finish_count"`
	//PreviousSignCount   uint `gorm:"previous_sign_count" json:"previous_sign_count"`
	//PreviousFinishCount uint `gorm:"previous_finish_count" json:"previous_finish_count"`
	workCount := &WorkCount{
		SourceChainId:       1,
		TargetChainId:       2,
		AnchorNodeId:        1,
		SignCount:           10,
		FinishCount:         20,
		PreviousSignCount:   30,
		PreviousFinishCount: 60,
	}
	id, err := obj.CreateWorkCount(workCount)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("work count id:", id)
	workCount = &WorkCount{
		SourceChainId:       1,
		TargetChainId:       2,
		AnchorNodeId:        1,
		SignCount:           10,
		FinishCount:         20,
		PreviousSignCount:   40,
		PreviousFinishCount: 80,
	}
	id, err = obj.CreateWorkCount(workCount)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("work count id:", id)
	workCount = &WorkCount{
		SourceChainId:       1,
		TargetChainId:       2,
		AnchorNodeId:        1,
		SignCount:           10,
		FinishCount:         20,
		PreviousSignCount:   50,
		PreviousFinishCount: 100,
	}
	id, err = obj.CreateWorkCount(workCount)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("work count id:", id)
}
func TestDataBaseAccessObject_CreateWorkCount3(t *testing.T) {
	//SourceChainId       uint `gorm:"source_chain_id" json:"source_chain_id"`
	//TargetChainId       uint `gorm:"target_chain_id" json:"target_chain_id"`
	//AnchorNodeId        uint `gorm:"anchor_node_id" json:"anchor_node_id"` //锚定节点编号
	//SignCount           uint `gorm:"sign_count" json:"sign_count"`
	//FinishCount         uint `gorm:"finish_count" json:"finish_count"`
	//PreviousSignCount   uint `gorm:"previous_sign_count" json:"previous_sign_count"`
	//PreviousFinishCount uint `gorm:"previous_finish_count" json:"previous_finish_count"`
	workCount := &WorkCount{
		SourceChainId:       1,
		TargetChainId:       2,
		AnchorNodeId:        1,
		SignCount:           10,
		FinishCount:         20,
		PreviousSignCount:   60,
		PreviousFinishCount: 120,
	}
	id, err := obj.CreateWorkCount(workCount)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("work count id:", id)

	workCount = &WorkCount{
		SourceChainId:       1,
		TargetChainId:       2,
		AnchorNodeId:        1,
		SignCount:           10,
		FinishCount:         20,
		PreviousSignCount:   70,
		PreviousFinishCount: 140,
	}
	id, err = obj.CreateWorkCount(workCount)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("work count id:", id)
	workCount = &WorkCount{
		SourceChainId:       1,
		TargetChainId:       2,
		AnchorNodeId:        1,
		SignCount:           10,
		FinishCount:         20,
		PreviousSignCount:   80,
		PreviousFinishCount: 160,
	}
	id, err = obj.CreateWorkCount(workCount)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("work count id:", id)
}
func TestDataBaseAccessObject_CreateWorkCount2(t *testing.T) {
	//SourceChainId       uint `gorm:"source_chain_id" json:"source_chain_id"`
	//TargetChainId       uint `gorm:"target_chain_id" json:"target_chain_id"`
	//AnchorNodeId        uint `gorm:"anchor_node_id" json:"anchor_node_id"` //锚定节点编号
	//SignCount           uint `gorm:"sign_count" json:"sign_count"`
	//FinishCount         uint `gorm:"finish_count" json:"finish_count"`
	//PreviousSignCount   uint `gorm:"previous_sign_count" json:"previous_sign_count"`
	//PreviousFinishCount uint `gorm:"previous_finish_count" json:"previous_finish_count"`
	workCount := &WorkCount{
		SourceChainId:       1,
		TargetChainId:       2,
		AnchorNodeId:        1,
		SignCount:           10,
		FinishCount:         20,
		PreviousSignCount:   0,
		PreviousFinishCount: 0,
	}
	id, err := obj.CreateWorkCount(workCount)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("work count id:", id)

	workCount = &WorkCount{
		SourceChainId:       1,
		TargetChainId:       2,
		AnchorNodeId:        1,
		SignCount:           10,
		FinishCount:         20,
		PreviousSignCount:   10,
		PreviousFinishCount: 20,
	}
	id, err = obj.CreateWorkCount(workCount)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("work count id:", id)
	workCount = &WorkCount{
		SourceChainId:       1,
		TargetChainId:       2,
		AnchorNodeId:        1,
		SignCount:           10,
		FinishCount:         20,
		PreviousSignCount:   20,
		PreviousFinishCount: 40,
	}
	id, err = obj.CreateWorkCount(workCount)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("work count id:", id)
}
