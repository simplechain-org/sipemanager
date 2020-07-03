package controllers

import (
	"math/big"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

//签名奖励日志列表
func (this *Controller) ListSignReward(c *gin.Context) {
	//分页显示
	//发放时间
	//锚定节点名称
	//奖励池总额度
	//签名量占比
	//奖励值
	//交易哈希

}

//剩余奖池总额

//单笔签名奖励

//本期总签名数

//签名工作量占比

type AddSignRewardParam struct {
	AnchorNodeId uint     `json:"anchor_node_id" form:"anchor_node_id"`
	NodeId       uint     `json:"node_id" form:"node_id"`
	WalletId     uint     `json:"wallet_id" form:"wallet_id"`
	Password     string   `json:"password" form:"password"`
	Reward       *big.Int `json:"reward" form:"reward"` //奖励值
	Coin         string   `json:"coin" form:"coin"`     //奖励币种
}

func (this *Controller) AddSignReward(c *gin.Context) {
	var param AddSignRewardParam
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
		this.echoError(c, err)
		return
	}
	privateKey, err := blockchain.GetPrivateKey(wallet.Content, param.Password)
	if err != nil {
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
	hash, err := source.AccumulateRewards(config, callerConfig, param.Reward)
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
	signRewardLog := &dao.SignRewardLog{
		AnchorNodeId:    param.AnchorNodeId,
		AnchorNodeName:  anchorNode.Name,
		TransactionHash: hash,
		Reward:          param.Reward,
		Coin:            param.Coin,
		Sender:          address.String(),
	}
	id, err := this.dao.CreateSignRewardLog(signRewardLog)
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
			//时间到，做一下检测
			receipt, err := source.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateSignRewardLogStatus(id, uint(receipt.Status))
				if err != nil {
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

type ConfigureSignRewardParam struct {
	AnchorNodeId uint     `json:"anchor_node_id" form:"anchor_node_id"`
	NodeId       uint     `json:"node_id" form:"node_id"`
	WalletId     uint     `json:"wallet_id" form:"wallet_id"`
	Password     string   `json:"password" form:"password"`
	Reward       *big.Int `json:"reward" form:"reward"` //奖励值
	Coin         string   `json:"coin" form:"coin"`     //奖励币种
}

//配置签名奖励
func (this *Controller) ConfigureSignReward(c *gin.Context) {

}
