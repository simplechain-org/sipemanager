package controllers

import (
	"github.com/gin-gonic/gin"
	"sipemanager/blockchain"
	"sipemanager/dao"
)

func Register(router *gin.Engine, object *dao.DataBaseAccessObject) {
	c := &Controller{userClient: make(map[uint]*blockchain.Api),
		dao: object,
	}
	validateLogin := ValidateTokenMiddleware()
	router.POST("/api/v1/user/register", c.Register)
	router.POST("/api/v1/user/login", c.Login)
	router.GET("/api/v1/check/health", c.CheckHealth)

	//auth
	router.GET("/api/v1/block/list", validateLogin, c.GetPageBlock)
	router.GET("/api/v1/block/transaction/:number", validateLogin, c.GetBlockTransaction)
	router.GET("/api/v1/transaction/:hash", validateLogin, c.GetTransactionReceipt)
	router.GET("/api/v1/node/list", validateLogin, c.GetNodes)
	router.POST("/api/v1/node", validateLogin, c.AddNode)
	router.POST("/api/v1/node/change", validateLogin, c.ChangeNode)
	router.GET("/api/v1/node/current", validateLogin, c.GetUserCurrentNode)

	router.DELETE("/api/v1/chain/:chain_id", c.RemoveChain)

	router.GET("/api/v1/wallet/list", validateLogin, c.ListWallet)
	router.POST("/api/v1/wallet", validateLogin, c.AddWallet)

	router.GET("/api/v1/chain/list", validateLogin, c.GetChains)

	router.POST("/api/v1/contract", validateLogin, c.AddContract)
	router.POST("/api/v1/contract/instance", validateLogin, c.DeployContract)
	router.POST("/api/v1/contract/register", validateLogin, c.RegisterChain)
	router.POST("/api/v1/contract/produce", validateLogin, c.Produce)
	router.POST("/api/v1/contract/consume", validateLogin, c.Consume)
	router.GET("/api/v1/contract/transaction", validateLogin, c.ListCrossTransaction)
	router.GET("/api/v1/contract/list", validateLogin, c.ListContract)
	router.GET("/api/v1/contract/chain", validateLogin, c.GetContractOnChain)
	router.GET("/api/v1/contract/instance/list", validateLogin, c.GetContractInstances)

	router.POST("/api/v1/contract/instance/add", validateLogin, c.AddContractInstance)

	router.POST("/api/v1/chain/address", validateLogin, c.UpdateChainContractAddress)
	router.GET("/api/v1/chain/current", validateLogin, c.GetUserCurrentChain)
	router.GET("/api/v1/chain/info/:chain_id", validateLogin, c.GetChainInfo)
	router.POST("/api/v1/chain/create", validateLogin, c.CreateChain)

	router.GET("/api/v1/contract/register/list", validateLogin, c.ListRegisterChain)
	router.GET("/api/v1/contract/produce/list", validateLogin, c.ListMakerOrder)
	router.GET("/api/v1/contract/consume/list", validateLogin, c.ListTakerOrder)

	router.POST("/api/v1/contract/register/add", validateLogin, c.RegisterChainAdd)

}
