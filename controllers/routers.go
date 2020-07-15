package controllers

import (
	"sipemanager/blockchain"
	"sipemanager/dao"
	"sipemanager/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SwaggerDoc(router *gin.Engine) {
	docs.SwaggerInfo.Title = "Sipe Manager API"
	docs.SwaggerInfo.Description = "区块链管理系统api文档"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "192.168.4.109:8092"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

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
	router.POST("/api/v1/wallet/update", validateLogin, c.UpdateWallet)

	router.POST("/api/v1/contract", validateLogin, c.AddContract)
	router.POST("/api/v1/contract/instance", validateLogin, c.DeployContract)
	router.POST("/api/v1/contract/register", validateLogin, c.RegisterChain)
	router.POST("/api/v1/contract/produce", validateLogin, c.Produce)
	router.POST("/api/v1/contract/consume", validateLogin, c.Consume)
	router.GET("/api/v1/contract/transaction", validateLogin, c.ListCrossTransaction)
	router.GET("/api/v1/contract/list", validateLogin, c.ListContract)
	router.GET("/api/v1/contract/chain", validateLogin, c.GetContractOnChain)
	router.GET("/api/v1/contract/instance/list", validateLogin, c.GetContractInstances)
	router.POST("/api/v1/contract/instance/import", validateLogin, c.GetContractInstances)

	router.POST("/api/v1/contract/register/once", validateLogin, c.RegisterChainTwoWay)

	router.POST("/api/v1/contract/instance/add", validateLogin, c.AddContractInstance)

	router.GET("/api/v1/chain/current", validateLogin, c.GetUserCurrentChain)
	router.GET("/api/v1/chain/info/:chain_id", validateLogin, c.GetChainInfo)
	router.POST("/api/v1/chain/create", validateLogin, c.CreateChain)

	router.GET("/api/v1/contract/register/list", validateLogin, c.ListRegisterChain)
	router.GET("/api/v1/contract/produce/list", validateLogin, c.ListMakerOrder)
	router.GET("/api/v1/contract/consume/list", validateLogin, c.ListTakerOrder)

	router.POST("/api/v1/contract/register/add", validateLogin, c.RegisterChainAdd)

	router.POST("/api/v1/retro/list", validateLogin, c.RetroActiveList)
	router.POST("/api/v1/retro/add", validateLogin, c.RetroActiveAdd)

	router.GET("/api/v1/chart/feeAndCount/list", c.FeeAndCount)
	router.GET("/api/v1/reward/list", validateLogin, c.ListSignReward)
	router.GET("/api/v1/reward/total", validateLogin, c.GetTotalReward)
	router.GET("/api/v1/reward/chain", validateLogin, c.GetChainReward)
	router.GET("/api/v1/anchor/work/count", validateLogin, c.GetAnchorWorkCount)
	router.POST("/api/v1/reward/add", validateLogin, c.AddSignReward)
	router.POST("/api/v1/reward/configure", validateLogin, c.ConfigureSignReward)
	router.GET("/api/v1/service/charge/list", validateLogin, c.ListServiceCharge)
	router.POST("/api/v1/service/charge/add", validateLogin, c.AddServiceCharge)
	router.GET("/api/v1//service/charge/fee", validateLogin, c.GetServiceChargeFee)
	router.POST("/api/v1/anchor/node/add", validateLogin, c.AddAnchorNode)
	router.POST("/api/v1/anchor/node/remove", validateLogin, c.RemoveAnchorNode)
	router.DELETE("/api/v1/wallet/remove", validateLogin, c.RemoveWallet)
	router.POST("/api/v1/punishment/add", validateLogin, c.AddPunishment)
	router.GET("/api/v1/punishment/list", validateLogin, c.ListPunishment)
	router.GET("/api/v1/anchor/node/list", validateLogin, c.ListAnchorNode)
	router.GET("/api/v1/anchor/node/obtain", validateLogin, c.GetAnchorNode)
	router.PUT("/api/v1/node", validateLogin, c.UpdateNode)
	router.DELETE("/api/v1/node/remove/:id", validateLogin, c.DeleteNode)
	router.PUT("/api/v1/chain/update", validateLogin, c.UpdateChain)
	router.GET("/api/v1/chain/list", validateLogin, c.ListChain)


}

type BlockChannel struct {
	ChainId     uint
	BlockNumber int64
	currentNode dao.InstanceNodes
}

func ListenEvent(object *dao.DataBaseAccessObject) {
	c := &Controller{userClient: make(map[uint]*blockchain.Api),
		dao: object,
	}
	//go c.ListenCrossEvent()
	go c.ListenBlock()
	go c.ListenAnchors()

}
