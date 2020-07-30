package controllers

import (
	"fmt"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/robfig/cron/v3"
	"github.com/simplechain-org/go-simplechain/common"
)

func (this *Controller) ListenWorkCount() {
	c := cron.New()
	spec := "* 8 * * *" // 每天上午8点执行一下
	c.AddFunc(spec, func() {
		//获取所有有效的的锚定节点
		anchorNodes, err := this.dao.ListAnchorNode()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, anchorNode := range anchorNodes {
			//记得双向记录
			sourceContract, err := this.dao.GetContractByChainId(anchorNode.SourceChainId)
			if err != nil {
				fmt.Println(err)
				continue
			}
			targetContract, err := this.dao.GetContractByChainId(anchorNode.TargetChainId)
			if err != nil {
				fmt.Println(err)
				continue
			}
			sourceChain, err := this.dao.GetChain(anchorNode.SourceChainId)
			if err != nil {
				fmt.Println(err)
				continue
			}
			targetChain, err := this.dao.GetChain(anchorNode.TargetChainId)
			if err != nil {
				fmt.Println(err)
				continue
			}
			sourceNode, err := this.dao.GetNodeByChainId(anchorNode.SourceChainId)
			if err != nil {
				fmt.Println(err)
				continue
			}
			targetNode, err := this.dao.GetNodeByChainId(anchorNode.TargetChainId)
			if err != nil {
				fmt.Println(err)
				continue
			}
			sourceApi, err := this.getApiByNodeId(sourceNode.ID)
			if err != nil {
				fmt.Println(err)
				continue
			}
			targetApi, err := this.getApiByNodeId(targetNode.ID)
			if err != nil {
				fmt.Println(err)
				continue
			}

			config := &blockchain.AnchorNodeRewardConfig{
				AbiData:         []byte(sourceContract.Abi),
				ContractAddress: common.HexToAddress(sourceContract.Address),
				TargetNetworkId: targetChain.NetworkId,
				AnchorAddress:   common.HexToAddress(anchorNode.Address),
			}
			callerConfig := &blockchain.CallerConfig{
				NetworkId: sourceChain.NetworkId,
			}
			//正向
			signCount, finishCount, err := sourceApi.GetAnchorWorkCount(config, callerConfig)
			if err != nil {
				fmt.Println(err)
				continue
			}
			var previousSignCount uint64
			var previousFinishCount uint64
			workCount, err := this.dao.GetlatestWorkCount(anchorNode.SourceChainId, anchorNode.TargetChainId, anchorNode.ID)
			if err != nil {
				//没有记录
				fmt.Println(err)
			} else {
				previousSignCount = workCount.PreviousSignCount
				previousFinishCount = workCount.PreviousFinishCount
			}
			sign := signCount.Uint64() - previousSignCount

			finish := finishCount.Uint64() - previousFinishCount
			sourceWorkCount := &dao.WorkCount{
				SourceChainId:       anchorNode.SourceChainId,
				TargetChainId:       anchorNode.TargetChainId,
				AnchorNodeId:        anchorNode.ID,
				SignCount:           sign,
				FinishCount:         finish,
				PreviousSignCount:   signCount.Uint64(),
				PreviousFinishCount: finishCount.Uint64(),
			}
			id, err := this.dao.CreateWorkCount(sourceWorkCount)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("CreateWorkCount id:", id)
			//反向
			config = &blockchain.AnchorNodeRewardConfig{
				AbiData:         []byte(targetContract.Abi),
				ContractAddress: common.HexToAddress(targetContract.Address),
				TargetNetworkId: sourceChain.NetworkId,
				AnchorAddress:   common.HexToAddress(anchorNode.Address),
			}
			callerConfig = &blockchain.CallerConfig{
				NetworkId: targetChain.NetworkId,
			}
			signCount, finishCount, err = targetApi.GetAnchorWorkCount(config, callerConfig)
			if err != nil {
				fmt.Println(err)
				continue
			}
			workCount, err = this.dao.GetlatestWorkCount(anchorNode.TargetChainId, anchorNode.SourceChainId, anchorNode.ID)
			if err != nil {
				//没有记录
				fmt.Println(err)
				previousSignCount = 0
				previousFinishCount = 0
			} else {
				previousSignCount = workCount.PreviousSignCount
				previousFinishCount = workCount.PreviousFinishCount
			}
			sign = signCount.Uint64() - previousSignCount

			finish = finishCount.Uint64() - previousFinishCount
			targetWorkCount := &dao.WorkCount{
				SourceChainId:       anchorNode.TargetChainId,
				TargetChainId:       anchorNode.SourceChainId,
				AnchorNodeId:        anchorNode.ID,
				SignCount:           sign,
				FinishCount:         finish,
				PreviousSignCount:   signCount.Uint64(),
				PreviousFinishCount: finishCount.Uint64(),
			}
			id, err = this.dao.CreateWorkCount(targetWorkCount)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("CreateWorkCount id:", id)
		}
	})
	c.Start()
}
