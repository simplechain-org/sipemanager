package controllers

import (
	"fmt"
	"math/big"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"
	"github.com/gin-gonic/gin"
)

//已报销手续费列表
func (this *Controller) ListServiceCharge(c *gin.Context) {
	//报销时间
	//锚定节点名称
	//交易哈希
	//报销手续费
	//分页显示
}

type AddServiceChargeParam struct {
	AnchorNodeId uint     `json:"anchor_node_id" form:"anchor_node_id"`
	NodeId       uint     `json:"node_id" form:"node_id"`
	WalletId     uint     `json:"wallet_id" form:"wallet_id"`
	Password     string   `json:"password" form:"password"`
	Fee          *big.Int `json:"fee" form:"fee"`   //报销手续费
	Coin         string   `json:"coin" form:"coin"` //报销的币种
}

//新增手续费报销
func (this *Controller) AddServiceCharge(c *gin.Context) {
	var param AddServiceChargeParam
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, err)
		return
	}
	anchorNode, err := this.dao.GetAnchorNode(param.AnchorNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	node, err := this.dao.GetNodeById(param.NodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	source, err := this.getApiByNodeId(param.NodeId)
	if err != nil {
		fmt.Println("GetApiByNodeId:", err.Error())
		this.echoError(c, err)
		return
	}
	//链的合约
	contract, err := this.dao.GetContractByChainId(node.ChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var targetChainId uint
	if anchorNode.SourceChainId == node.ChainId {
		targetChainId = anchorNode.TargetChainId
	} else if anchorNode.TargetChainId == node.ChainId {
		targetChainId = anchorNode.SourceChainId
	}
	chain, err := this.dao.GetChain(targetChainId)
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

	config := &blockchain.AnchorNodeRewardConfig{
		AbiData:         []byte(contract.Abi),
		ContractAddress: common.HexToAddress(contract.Address),
		TargetNetworkId: source.GetNetworkId(),
		AnchorAddress:   common.HexToAddress(anchorNode.Address),
	}
	callerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  chain.NetworkId,
	}
	hash, err := source.AccumulateRewards(config, callerConfig, param.Fee)
	if err != nil {
		this.echoError(c, err)
		return
	}
	//AnchorNodeId    uint     //锚定节点编号
	//AnchorNodeName  uint     //锚定节点名称，冗余方便查询
	//TransactionHash string   //交易哈希
	//Fee             *big.Int //报销手续费
	//Coin            string   //报销的币种
	//Sender          string   //出账账户地址
	//BlockNumber     uint     //区块高度
	//Status          uint     //状态

	//todo 区块高度
	serviceChargeLog := &dao.ServiceChargeLog{
		AnchorNodeId:    param.AnchorNodeId,
		AnchorNodeName:  anchorNode.Name,
		TransactionHash: hash,
		Fee:             param.Fee,
		Coin:            param.Coin,
		Sender:          address.String(),
	}
	id, err := this.dao.CreateServiceChargeLog(serviceChargeLog)
	if err != nil {
		this.echoError(c, err)
		return
	}
	go func(source *blockchain.Api, id uint, hash string) {
		//todo 应考虑回滚，让过了多少个块以后确认为不可推翻的（已成永久）
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
			receipt, err := source.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateServiceChargeLogSourceStatus(id, uint(receipt.Status))
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
	}(source, id, hash)
	this.echoSuccess(c, "success")

}

//统计累计消耗手续费,从报销表

//统计累计已报销手续费
//计算本期应报销手续费

//这里应注意到：一条链只能报销一种token(币)
