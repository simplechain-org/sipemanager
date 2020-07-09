package controllers

import (
	"fmt"
	"sipemanager/dao"
	"strconv"
	"strings"
)

type AnchorsNodes struct {
	Id      uint
	Name    string
	Address string
}

func (this *Controller) findAnchors(node dao.InstanceNodes) {
	anchorsId := strings.Split("1,3,2", ",")
	for _, item := range anchorsId {
		n, _ := strconv.Atoi(item)
		anchor, err := this.dao.GetAnchorNode(uint(n))
		if err != nil {

		}
		//txs, txErr := this.dao.GetTxByAnchors(node.ChainId, anchor.Address, node.CrossAddress)
		//if txErr != nil {
		//
		//}
		//
		//txAnchot := dao.TxAnchors{
		//	From:          anchor.Address,
		//	To:            node.Address,
		//	SourceChainId: anchor.SourceChainId,
		//	TargetChainId: anchor.TargetChainId,
		//}
		//this.dao.CreateTxAnchors(txAnchot)
		fmt.Printf("fdfd %+v\n", anchor)
	}

}
