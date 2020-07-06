package controllers

import (
	"fmt"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"
	"github.com/gin-gonic/gin"
)

type AddAnchorNodeParam struct {
	SourceChainId uint   `json:"source_chain_id" form:"source_chain_id"`
	TargetChainId uint   `json:"target_chain_id" form:"target_chain_id"`
	SourceNodeId  uint   `json:"source_node_id" form:"source_node_id"`
	TargetNodeId  uint   `json:"target_node_id" form:"target_node_id"`
	AnchorAddress string `json:"anchor_address" form:"anchor_address"`
	AnchorName    string `json:"anchor_name" form:"anchor_name"`
	WalletId      uint   `json:"wallet_id" form:"wallet_id"`
	Password      string `json:"password" form:"password"`
}

//锚定节点管理
//新增锚定节点
func (this *Controller) AddAnchorNode(c *gin.Context) {
	var param AddAnchorNodeParam
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, err)
		return
	}
	//调用合约增加锚定节点，要注意是双链
	//添加到数据库
	sourceContract, err := this.dao.GetContractByChainId(param.SourceChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		fmt.Println("GetWallet:", err.Error())
		this.echoError(c, err)
		return
	}
	privateKey, err := blockchain.GetPrivateKey(wallet.Content, param.Password)
	if err != nil {
		fmt.Println("GetPrivateKey:", err.Error())
		this.echoError(c, err)
		return
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	source, err := this.getApiByNodeId(param.SourceNodeId)
	if err != nil {
		fmt.Println("GetApiByNodeId:", err.Error())
		this.echoError(c, err)
		return
	}
	target, err := this.getApiByNodeId(param.TargetNodeId)
	if err != nil {
		fmt.Println("GetApiByNodeId:", err.Error())
		this.echoError(c, err)
		return
	}

	configSource := &blockchain.AnchorNodeConfig{
		AbiData:         []byte(sourceContract.Abi),
		ContractAddress: common.HexToAddress(sourceContract.Address),
		TargetNetworkId: target.GetNetworkId(),
		AnchorAddresses: []common.Address{common.HexToAddress(param.AnchorAddress)},
	}
	callerConfigSource := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  source.GetNetworkId(),
	}

	sourceHash, err := source.AddAnchors(configSource, callerConfigSource)
	if err != nil {
		this.echoError(c, err)
		return
	}

	targetContract, err := this.dao.GetContractByChainId(param.TargetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}

	configTarget := &blockchain.AnchorNodeConfig{
		AbiData:         []byte(targetContract.Abi),
		ContractAddress: common.HexToAddress(targetContract.Address),
		TargetNetworkId: source.GetNetworkId(),
		AnchorAddresses: []common.Address{common.HexToAddress(param.AnchorAddress)},
	}
	callerConfigTarget := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  target.GetNetworkId(),
	}
	targetHash, err := target.AddAnchors(configTarget, callerConfigTarget)
	if err != nil {
		this.echoError(c, err)
		return
	}
	//todo 万一其中一条链添加失败，这个时候该怎么处理

	anchorNode := &dao.AnchorNode{
		Name:          param.AnchorName,
		Address:       param.AnchorAddress,
		SourceChainId: param.SourceChainId,
		TargetChainId: param.TargetChainId,
		SourceHash:    sourceHash,
		TargetHash:    targetHash,
	}
	id, err := this.dao.CreateAnchorNode(anchorNode)
	if err != nil {
		this.echoError(c, err)
		return
	}

	go func(api *blockchain.Api, id uint, hash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
			receipt, err := api.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateSourceStatus(id, uint(receipt.Status))
				if err != nil {
					fmt.Println(err)
					continue
				}
				break
			}
			if i >= maxCount {
				break
			}
			i++
		}
	}(source, id, sourceHash)

	go func(api *blockchain.Api, id uint, hash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
			receipt, err := api.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateTargetStatus(id, uint(receipt.Status))
				if err != nil {
					fmt.Println(err)
					continue
				}
				break
			}
			if i >= maxCount {
				break
			}
			i++
		}
	}(target, id, targetHash)
	this.echoSuccess(c, "success")
}

type RemoveAnchorNodeParam struct {
	AnchorNodeId uint   `json:"anchor_node_id" form:"anchor_node_id"`
	SourceNodeId uint   `json:"source_node_id" form:"source_node_id"`
	TargetNodeId uint   `json:"target_node_id" form:"target_node_id"`
	WalletId     uint   `json:"wallet_id" form:"wallet_id"`
	Password     string `json:"password" form:"password"`
}

//锚定节点管理
//删除锚定节点
func (this *Controller) RemoveAnchorNode(c *gin.Context) {
	//调用合约增加锚定节点，要注意是双链
	//从数据库中删除
	var param RemoveAnchorNodeParam
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, err)
		return
	}
	anchorNode, err := this.dao.GetAnchorNode(param.AnchorNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	sourceContract, err := this.dao.GetContractByChainId(anchorNode.SourceChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	targetContract, err := this.dao.GetContractByChainId(anchorNode.TargetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	source, err := this.getApiByNodeId(param.SourceNodeId)
	if err != nil {
		fmt.Println("GetApiByNodeId:", err.Error())
		this.echoError(c, err)
		return
	}
	target, err := this.getApiByNodeId(param.TargetNodeId)
	if err != nil {
		fmt.Println("GetApiByNodeId:", err.Error())
		this.echoError(c, err)
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		fmt.Println("GetWallet:", err.Error())
		this.echoError(c, err)
		return
	}
	privateKey, err := blockchain.GetPrivateKey(wallet.Content, param.Password)
	if err != nil {
		fmt.Println("GetPrivateKey:", err.Error())
		this.echoError(c, err)
		return
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	targetConfig := &blockchain.AnchorNodeConfig{
		AbiData:         []byte(targetContract.Abi),
		ContractAddress: common.HexToAddress(targetContract.Address),
		TargetNetworkId: source.GetNetworkId(),
		AnchorAddresses: []common.Address{common.HexToAddress(anchorNode.Address)},
	}
	targetCallerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  target.GetNetworkId(),
	}
	targetHash, err := target.RemoveAnchors(targetConfig, targetCallerConfig)
	if err != nil {
		this.echoError(c, err)
		return
	}
	sourceConfig := &blockchain.AnchorNodeConfig{
		AbiData:         []byte(sourceContract.Abi),
		ContractAddress: common.HexToAddress(sourceContract.Address),
		TargetNetworkId: target.GetNetworkId(),
		AnchorAddresses: []common.Address{common.HexToAddress(anchorNode.Address)},
	}
	sourceCallerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  source.GetNetworkId(),
	}
	sourceHash, err := target.RemoveAnchors(sourceConfig, sourceCallerConfig)
	if err != nil {
		this.echoError(c, err)
		return
	}
	go func(source *blockchain.Api, target *blockchain.Api, id uint, sourceHash string, targetHash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
			receipt, err := source.TransactionReceipt(common.HexToHash(sourceHash))
			if err == nil && receipt != nil {
				if receipt.Status == 1 {
					receipt, err := target.TransactionReceipt(common.HexToHash(targetHash))
					if err == nil && receipt != nil {
						if receipt.Status == 1 {
							//两条链都已经完成合约调用成功，那么就移除删除数据库中锚定节点的数据
							err = this.dao.RemoveAnchorNode(id)
							if err != nil {
								fmt.Println("RemoveAnchorNode error:", err)
								continue
							}
							break
						}
					}
				}
			}
			if i >= maxCount {
				break
			}
			i++
		}
	}(source, target, param.AnchorNodeId, sourceHash, targetHash)
}

//锚定节点列表
//makefinish应该是有两方面的费用，A链和B链
func (this *Controller) ListAnchorNode(c *gin.Context) {
	//创建时间
	//锚定节点名称
	//归属链A
	//归属链B
	//makefinish手续费
	//已报销手续费
	//累计有效签名  //
	//累计发放奖励  //
	//质押金额
	//身份状态

	//分页显示，每页10条记录

}
