package controllers

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"math/big"
	"sipemanager/blockchain"
	"sipemanager/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/cross/core"
	"github.com/simplechain-org/go-simplechain/params"

	"sipemanager/dao"
)

// @Summary 补签列表
// @Tags RetroActiveList
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param current_page formData uint32 true "当前页，默认1"
// @Param page_size formData uint32 true "页的记录数，默认10"
// @Param status formData int false "补签状态"
// @Success 200 {object} JsonResult{data=[]dao.RetroActive}
// @Router /retro/list [post]
func (this *Controller) RetroActiveList(c *gin.Context) {
	type Param struct {
		CurrentPage uint32 `json:"current_page"`
		PageSize    uint32 `json:"page_size"`
		Status      int    `json:"status"`
	}
	var param Param
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, err)
		return
	}
	offset := (param.CurrentPage - 1) * param.PageSize
	var result []dao.RetroActive
	var err error
	if param.Status == 0 {
		fmt.Println(1)
		result, err = this.dao.ListRetroActive(offset, param.PageSize)
		if err != nil {
			this.echoError(c, err)
			return
		}
	} else {
		fmt.Println(2)
		result, err = this.dao.ListRetroActiveByStatus(param.Status, offset, param.PageSize)
		if err != nil {
			this.echoError(c, err)
			return
		}
	}
	fmt.Println(3)
	this.echoResult(c, result)
}

// @Summary 添加补签记录
// @Tags RetroActiveAdd
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param tx_hash formData string true "交易hash"
// @Param network_id formData uint64 true "交易所在链"
// @Success 200 {object} JsonResult{data=int}
// @Router /retro/add [post]
func (this *Controller) RetroActiveAdd(c *gin.Context) {
	var param dao.RetroActive
	if err := c.ShouldBindJSON(&param); err != nil {
		fmt.Println("1", err)
		this.echoError(c, err)
		return
	}
	if retro, err := this.dao.QueryRetroActive(param.TxHash); err == nil && retro != nil {
		fmt.Println("0", err)
		this.echoError(c, errors.New("already exit"))
		return
	}
	chain, err := this.dao.GetChainByNetWorkId(param.NetworkId)

	if err != nil {
		fmt.Println("2", err)
		this.echoError(c, err)
		return
	}

	contract, err := this.dao.GetContractInstanceById(chain.ContractInstanceId)
	if err != nil {
		fmt.Println("3", err)
		this.echoError(c, err)
		return
	}

	conAbi, err := this.dao.GetContractById(contract.ContractId)
	if err != nil {
		fmt.Println("111", err)
		this.echoError(c, err)
		return
	}
	user, err := this.GetUser(c)
	if err != nil {
		fmt.Println("4", err)
		this.echoError(c, err)
		return
	}

	sourceNode, err := this.dao.GetNodeByChainId(chain.ID)
	if err != nil {
		fmt.Println("11", err)
		this.echoError(c, err)
		return
	}
	api, err := this.getApiByNodeId(sourceNode.ID)
	if err != nil {
		fmt.Println("12", err)
		this.echoError(c, err)
		return
	}

	//api, err := this.getApi(user.ID,param.NetworkId)
	//if err != nil {
	//	fmt.Println("5",err)
	//	this.echoError(c, err)
	//	return
	//}
	//查询交易receipt
	receipt, err := api.TransactionReceipt(common.HexToHash(param.TxHash))
	if err != nil {
		fmt.Println("6", err)
		this.echoError(c, err)
		return
	}
	if len(receipt.Logs) > 0 {
		if receipt.Logs[0].Address == common.HexToAddress(contract.Address) {
			if len(receipt.Logs[0].Topics) >= 3 && receipt.Logs[0].Topics[0] == params.MakerTopic && len(receipt.Logs[0].Data) >= common.HashLength*6 {
				param.CtxId = receipt.Logs[0].Topics[1].Hex()
				param.Event = 1
			}
			if len(receipt.Logs[0].Topics) >= 3 && receipt.Logs[0].Topics[0] == params.TakerTopic && len(receipt.Logs[0].Data) >= common.HashLength*4 {
				param.CtxId = receipt.Logs[0].Topics[1].Hex()
				param.Event = 2
			}
		} else {
			fmt.Println("6", err)
			this.echoError(c, errors.New("address error"))
			return
		}
	} else {
		fmt.Println("7", err)
		this.echoError(c, errors.New("no logs"))
		return
	}
	wallets, err := this.dao.ListWalletByUserId(user.ID)
	if err != nil {
		fmt.Println("8", err)
		this.echoError(c, err)
		return
	}
	var wallet dao.Wallet
	if len(wallets) > 0 {
		wallet = wallets[0]
	} else {
		fmt.Println("9", err)
		this.echoError(c, errors.New("no wallets"))
		return
	}

	if param.Event == 1 {
		targetId, err := this.dao.GetTargetChainIdBySourceChainId(chain.ID)

		if err != nil {
			fmt.Println("12", err)
			this.echoError(c, err)
			return
		}
		chain, err := this.dao.GetChain(targetId)
		if err != nil {
			fmt.Println("14", err)

			this.echoError(c, err)
			return
		}
		stat, err := api.GetMakerTx(common.HexToHash(param.CtxId), common.HexToAddress(contract.Address), common.HexToAddress(wallet.Address), []byte(conAbi.Abi), big.NewInt(int64(chain.NetworkId)))

		fmt.Println(stat)
		if err != nil {
			fmt.Println("10", err)
			this.echoError(c, err)
			return
		}
		if stat {
			t, err := api.CtxGet(common.HexToHash(param.CtxId))
			if err != nil {
				fmt.Println("11", err)
				this.echoError(c, err)
				return
			}
			fmt.Println(t, param.CtxId)
			if t == nil || t.Status == core.CtxStatusPending {
				param.Status = 1
			} else {
				param.Status = 2
			}
		} else {
			param.Status = 2
		}
	} else {
		targetId, err := this.dao.GetTargetChainIdBySourceChainId(chain.ID)

		if err != nil {
			fmt.Println("12", err)
			this.echoError(c, err)
			return
		}
		targetChain, err := this.dao.GetChain(targetId)
		if err != nil {
			fmt.Println("121", err)
			this.echoError(c, err)
			return
		}
		objContract, err := this.dao.GetContractInstanceById(targetChain.ContractInstanceId)
		if err != nil {
			fmt.Println("13", err)
			this.echoError(c, err)
			return
		}
		//chain,err := this.dao.GetChain(targetId)
		//if err != nil {
		//	fmt.Println("14",err)
		//
		//	this.echoError(c, err)
		//	return
		//}

		sourceNode2, err := this.dao.GetNodeByChainId(targetId)
		if err != nil {
			fmt.Println("14", err)
			this.echoError(c, err)
			return
		}
		obApi, err := this.getApiByNodeId(sourceNode2.ID)
		if err != nil {
			fmt.Println("13", err)
			this.echoError(c, err)
			return
		}
		//obApi, err := this.getApi(user.ID,chain.NetworkId)
		//if err != nil {
		//	fmt.Println("5",err)
		//	this.echoError(c, err)
		//	return
		//}
		//todo obApi对称
		fmt.Println(objContract.Address, param.NetworkId)
		stat, err := obApi.GetMakerTx(common.HexToHash(param.CtxId), common.HexToAddress(objContract.Address), common.HexToAddress(wallet.Address), []byte(conAbi.Abi), big.NewInt(int64(param.NetworkId)))

		if err != nil {
			fmt.Println("16", err)
			this.echoError(c, err)
			return
		}
		if stat {
			param.Status = 1
		} else {
			param.Status = 2
		}
	}
	//验证流程
	id, err := this.dao.CreateRetroActive(&param)
	if err != nil {
		fmt.Println("17", err)
		this.echoError(c, err)
		return
	}
	this.echoResult(c, id)
}

// 每30分钟遍历更新一次补签状态
func (this *Controller) UpdateRetroActive() {
	cron := cron.New()
	cron.AddFunc("@every 30m", func() {
		fmt.Println("current UpdateRetroActive time is ", time.Now())
		var offset uint32
		var result []dao.RetroActive
		var err error
		wallet, err := this.dao.GetWallet(1)
		if err != nil {
			return
		}
		nodes, err := this.dao.GetInstancesJoinNode()
		if err != nil {
			fmt.Println(err)
			return
		}
		filterNodes := utils.RemoveRepByLoop(nodes)
		for ; ; offset += 20 {
			result, err = this.dao.ListRetroActive(offset, 20)
			if err != nil {
				fmt.Println(err)
				break
			}
			if len(result) == 0 {
				break
			}

			for _, v := range result {
				if v.Status == 1 {
					chain, err := this.dao.GetChainByNetWorkId(v.NetworkId)

					if err != nil {
						fmt.Println("2", err)
						break
					}
					if v.Event == 1 {
						var node *dao.InstanceNodes
						for _, v := range filterNodes {
							if v.NetworkId == v.NetworkId {
								node = &v
							}
						}

						var api *blockchain.Api
						if node != nil {
							api, err = GetRpcApi(*node)
							if err != nil {
								fmt.Println("5", err)
								break
							}
						}
						contract, err := this.dao.GetContractInstanceById(chain.ID)
						if err != nil {
							break
						}
						targetId, err := this.dao.GetTargetChainIdBySourceChainId(chain.ID)

						if err != nil {
							fmt.Println("12", err)
							break
						}
						chain, err := this.dao.GetChain(targetId)
						if err != nil {
							fmt.Println("14", err)
							break
						}
						conAbi, err := this.dao.GetContractById(contract.ContractId)
						if err != nil {
							fmt.Println("111", err)
							break
						}

						stat, err := api.GetMakerTx(common.HexToHash(v.CtxId), common.HexToAddress(contract.Address), common.HexToAddress(wallet.Address), []byte(conAbi.Abi), big.NewInt(int64(chain.NetworkId)))
						fmt.Println("777", stat, chain.NetworkId, v.CtxId)

						if err != nil {
							break
						}
						if stat {
							t, err := api.CtxGet(common.HexToHash(v.CtxId))
							if err != nil {
								break
							}
							fmt.Println("999", t, stat)
							if t == nil || t.Status == core.CtxStatusPending {
								v.Status = 1
							} else {
								v.Status = 2
							}
						} else {
							v.Status = 2
						}
					} else {
						targetId, err := this.dao.GetTargetChainIdBySourceChainId(chain.ID)

						if err != nil {
							fmt.Println("12", err)
							break
						}
						targetChain, err := this.dao.GetChain(targetId)
						if err != nil {
							fmt.Println("14", err)
							break
						}
						objContract, err := this.dao.GetContractInstanceById(targetChain.ContractInstanceId)
						if err != nil {
							fmt.Println("13", err)
							break
						}
						conAbi, err := this.dao.GetContractById(objContract.ContractId)
						if err != nil {
							fmt.Println("111", err)
							break
						}
						chain, err := this.dao.GetChain(targetId)

						if err != nil {
							fmt.Println("14", err)
							break
						}
						var node *dao.InstanceNodes
						for _, v := range filterNodes {
							if v.NetworkId == chain.NetworkId {
								node = &v
							}
						}

						var obApi *blockchain.Api
						if node != nil {
							obApi, err = GetRpcApi(*node)
							if err != nil {
								fmt.Println("5", err)
								break
							}
						}
						fmt.Println(objContract.Address, v.NetworkId)
						stat, err := obApi.GetMakerTx(common.HexToHash(v.CtxId), common.HexToAddress(objContract.Address), common.HexToAddress(wallet.Address), []byte(conAbi.Abi), big.NewInt(int64(v.NetworkId)))

						if err != nil {
							fmt.Println("16", err)
							break
						}
						fmt.Println("111", stat)
						if stat {
							v.Status = 1
						} else {
							v.Status = 2
						}
					}
					if v.Status == 2 {
						err := this.dao.UpdateRetroActiveStatus(v.ID, 2)
						if err != nil {
							fmt.Println(err)
							break
						}
					}
				}
			}
		}
	})
	cron.Start()
}
