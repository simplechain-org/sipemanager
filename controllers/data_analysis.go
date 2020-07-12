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

func (this *Controller) AnalysisAnchors() {
	registers, err := this.dao.ListChainRegisterByStatus(1)
	for _, register := range registers {
		sourceChain, err := this.dao.GetChain(register.SourceChainId)
		targetChain, err := this.dao.GetChain(register.TargetChainId)
		if err != nil {

		}
		anchorIds := strings.Split(register.AnchorAddresses, ",")
		for _, anchorId := range anchorIds {
			n, _ := strconv.Atoi(anchorId)
			anchor, err := this.dao.GetAnchorNode(uint(n))
			if err != nil {

			}
			txAnchor := dao.TxAnchors{
				AnchorAddress:   anchor.Address,
				SourceChainId:   register.SourceChainId,
				TargetChainId:   register.TargetChainId,
				AnchorId:        anchor.ID,
				ChainId:         register.SourceChainId,
				SourceNetworkId: sourceChain.NetworkId,
				TargetNetworkId: targetChain.NetworkId,
			}
			TxErr := this.dao.QueryTxByHours(txAnchor, "makerFinish")
			if TxErr != nil {
				fmt.Printf("fdfdfd-", TxErr.Error())
			}
		}
	}

	if err != nil {

	}
}

func (this *Controller) MakerFinishAnalysis() {

}

func (this *Controller) CountAnalysis() {

}
