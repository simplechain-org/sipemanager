package controllers

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/simplechain-org/go-simplechain/cross/core"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/params"
	"github.com/gin-gonic/gin"

	"sipemanager/dao"
)

func (this *Controller) RetroActiveList(c *gin.Context) {
	type Param struct {
		Page  uint32 `json:"page"`
		Limit uint32 `json:"limit"`
		Status int `json:"status"`
	}
	var param Param
	if err := c.ShouldBindJSON(&param); err != nil {
		this.echoError(c, err)
		return
	}
	offset := (param.Page - 1) * param.Limit
	var result []dao.RetroActive
	var err error
	if param.Status == 0 {
		result, err = this.dao.ListRetroActive(offset,param.Limit)
		if err != nil {
			this.echoError(c, err)
			return
		}
	} else {
		result, err = this.dao.ListRetroActiveByStatus(param.Status,offset,param.Limit)
		if err != nil {
			this.echoError(c, err)
			return
		}
	}

	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}

	wallets, err := this.dao.ListWalletByUserId(user.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var wallet dao.Wallet
	if len(wallets) > 0 {
		wallet = wallets[0]
	} else {
		this.echoError(c, errors.New("no wallets"))
		return
	}

	for k,v := range result {
		if v.Status == 1 {
			chainId,err :=  this.dao.GetChainIdByNetworkId(v.NetworkId)
			if err != nil {
				fmt.Println("2",err)
				this.echoError(c, err)
				return
			}
			if v.Event == 1 {
				api, err := this.getApi(user.ID,v.NetworkId)
				if err != nil {
					fmt.Println("5",err)
					this.echoError(c, err)
					return
				}
				contract, err := this.dao.GetContractByChainId(chainId)
				if err != nil {
					this.echoError(c, err)
					return
				}
				targetId,err := this.dao.GetTargetChainIdBySourceChainId(chainId)
				if err != nil {
					fmt.Println("12",err)
					this.echoError(c, err)
					return
				}
				chain,err := this.dao.GetChain(targetId)
				if err != nil {
					fmt.Println("14",err)

					this.echoError(c, err)
					return
				}

				stat,err := api.GetMakerTx(common.HexToHash(v.CtxId),common.HexToAddress(contract.Address),common.HexToAddress(wallet.Address),[]byte(contract.Abi),big.NewInt(int64(chain.NetworkId)))
				fmt.Println("777",stat,chain.NetworkId,v.CtxId)
				if err != nil {
					this.echoError(c, err)
					return
				}
				if stat {
					t,err := api.CtxGet(common.HexToHash(v.CtxId))
					if err != nil {
						this.echoError(c, err)
						return
					}
					fmt.Println("999",t,stat)
					if t == nil || t.Status == core.CtxStatusPending {
						v.Status = 1
					} else {
						v.Status = 2
					}
				} else {
					v.Status = 2
				}
			} else {
				targetId,err := this.dao.GetTargetChainIdBySourceChainId(chainId)
				if err != nil {
					fmt.Println("12",err)
					this.echoError(c, err)
					return
				}
				objContract, err := this.dao.GetContractByChainId(targetId)
				if err != nil {
					fmt.Println("13",err)
					this.echoError(c, err)
					return
				}
				chain,err := this.dao.GetChain(targetId)
				if err != nil {
					fmt.Println("14",err)

					this.echoError(c, err)
					return
				}
				obApi, err := this.getApi(user.ID,chain.NetworkId)
				if err != nil {
					fmt.Println("5",err)
					this.echoError(c, err)
					return
				}
				//todo obApi对称
				fmt.Println(objContract.Address,v.NetworkId)
				stat,err := obApi.GetMakerTx(common.HexToHash(v.CtxId),common.HexToAddress(objContract.Address),common.HexToAddress(wallet.Address),[]byte(objContract.Abi),big.NewInt(int64(v.NetworkId)))
				if err != nil {
					fmt.Println("16",err)
					this.echoError(c, err)
					return
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
					this.echoError(c, err)
					return
				}
				result[k].Status = 2
			}
		}
	}
	this.echoResult(c, result)
}

func (this *Controller) RetroActiveAdd(c *gin.Context) {
	var param dao.RetroActive
	if err := c.ShouldBindJSON(&param); err != nil {
		fmt.Println("1",err)
		this.echoError(c, err)
		return
	}
	if retro,err := this.dao.QueryRetroActive(param.TxHash);err == nil && retro != nil {
		fmt.Println("0",err)
		this.echoError(c, errors.New("already exit"))
		return
	}
	chainId,err :=  this.dao.GetChainIdByNetworkId(param.NetworkId)
	if err != nil {
		fmt.Println("2",err)
		this.echoError(c, err)
		return
	}

	contract, err := this.dao.GetContractByChainId(chainId)
	if err != nil {
		fmt.Println("3",err)
		this.echoError(c, err)
		return
	}

	user, err := this.GetUser(c)
	if err != nil {
		fmt.Println("4",err)
		this.echoError(c, err)
		return
	}
	api, err := this.getApi(user.ID,param.NetworkId)
	if err != nil {
		fmt.Println("5",err)
		this.echoError(c, err)
		return
	}
	//查询交易receipt
	receipt, err := api.TransactionReceipt(common.HexToHash(param.TxHash))
	if err != nil {
		fmt.Println("6",err)
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
			fmt.Println("6",err)
			this.echoError(c, errors.New("address error"))
			return
		}
	} else {
		fmt.Println("7",err)
		this.echoError(c, errors.New("no logs"))
		return
	}
	wallets, err := this.dao.ListWalletByUserId(user.ID)
	if err != nil {
		fmt.Println("8",err)
		this.echoError(c, err)
		return
	}
	var wallet dao.Wallet
	if len(wallets) > 0 {
		wallet = wallets[0]
	} else {
		fmt.Println("9",err)
		this.echoError(c, errors.New("no wallets"))
		return
	}

	if param.Event == 1 {
		targetId,err := this.dao.GetTargetChainIdBySourceChainId(chainId)
		if err != nil {
			fmt.Println("12",err)
			this.echoError(c, err)
			return
		}
		chain,err := this.dao.GetChain(targetId)
		if err != nil {
			fmt.Println("14",err)

			this.echoError(c, err)
			return
		}
		stat,err := api.GetMakerTx(common.HexToHash(param.CtxId),common.HexToAddress(contract.Address),common.HexToAddress(wallet.Address),[]byte(contract.Abi),big.NewInt(int64(chain.NetworkId)))
		fmt.Println(stat)
		if err != nil {
			fmt.Println("10",err)
			this.echoError(c, err)
			return
		}
		if stat {
			t,err := api.CtxGet(common.HexToHash(param.CtxId))
			if err != nil {
				fmt.Println("11",err)
				this.echoError(c, err)
				return
			}
			fmt.Println(t,param.CtxId)
			if t == nil || t.Status == core.CtxStatusPending {
				param.Status = 1
			} else {
				param.Status = 2
			}
		} else {
			param.Status = 2
		}
	} else {
		targetId,err := this.dao.GetTargetChainIdBySourceChainId(chainId)
		if err != nil {
			fmt.Println("12",err)
			this.echoError(c, err)
			return
		}
		objContract, err := this.dao.GetContractByChainId(targetId)
		if err != nil {
			fmt.Println("13",err)
			this.echoError(c, err)
			return
		}
		chain,err := this.dao.GetChain(targetId)
		if err != nil {
			fmt.Println("14",err)

			this.echoError(c, err)
			return
		}
		obApi, err := this.getApi(user.ID,chain.NetworkId)
		if err != nil {
			fmt.Println("5",err)
			this.echoError(c, err)
			return
		}
		//todo obApi对称
		fmt.Println(objContract.Address,param.NetworkId)
		stat,err := obApi.GetMakerTx(common.HexToHash(param.CtxId),common.HexToAddress(objContract.Address),common.HexToAddress(wallet.Address),[]byte(objContract.Abi),big.NewInt(int64(param.NetworkId)))
		if err != nil {
			fmt.Println("16",err)
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
		fmt.Println("17",err)
		this.echoError(c, err)
		return
	}
	this.echoResult(c, id)
}