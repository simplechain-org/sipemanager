package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"

	"sipemanager/dao"
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
				fmt.Println("AnalysisAnchors", err)
				continue
			}
			txAnchor := dao.TxAnchors{
				AnchorAddress:   anchor.Address,
				SourceChainId:   register.SourceChainId,
				TargetChainId:   register.TargetChainId,
				AnchorId:        anchor.ID,
				ChainId:         register.SourceChainId,
				SourceNetworkId: sourceChain.NetworkId,
				TargetNetworkId: targetChain.NetworkId,
				ContractAddress: register.Address,
			}
			TxHourErr := this.dao.QueryTxByHours(txAnchor, "makerFinish")
			TxDayErr := this.dao.QueryTxByDays(txAnchor, "makerFinish")
			TxWeekErr := this.dao.QueryTxByWeeks(txAnchor, "makerFinish")
			if TxHourErr != nil || TxDayErr != nil || TxWeekErr != nil {
				fmt.Printf("-------23-----%+v\n", TxHourErr.Error())
			}
		}
	}

	if err != nil {

	}
}

//统计锚定节点验证数和手续费
// @Summary 统计锚定节点验证数和手续费
// @Tags Chart
// @Accept  json
// @Produce  json
// @Param startTime query string true "hour:2020-07-10 12:00:00 day:2020-07-10 week:202025"
// @Param endTime query string true "hour:2020-07-12 12:00:00 day:2020-07-12 week:202027"
// @Param chainId query int true "chainId"
// @Param timeType query string true "hour,day,week"
// @Success 200 {object} JsonResult{data=[]dao.TxAnchorsNode}
// @Router /chart/feeAndCount/list [get]
func (this *Controller) FeeAndCount(c *gin.Context) {
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	chainIdParam := c.Query("chainId")
	timeType := c.Query("timeType")
	chainId, err := strconv.Atoi(chainIdParam)
	anchors, err := this.dao.QueryAnchors(startTime, endTime, chainId, timeType)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, anchors)
}

//分叉监控
// @Summary 分叉监控
// @Tags Chart
// @Accept  json
// @Produce  json
// @Success 200 {object} JsonResult{data=[]dao.MaxUncle}
// @Router /chart/maxUncle/list [get]
func (this *Controller) MaxUncle(c *gin.Context) {
	anchors, err := this.dao.QueryMaxUncle()
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, anchors)
}

// @Summary 跨链交易对列表
// @Tags Chart
// @Accept  json
// @Produce  json
// @Success 200 {object} JsonResult{data=dao.TokenListInterface}
// @Router /chart/txTokenList/list [get]
func (this *Controller) TxTokenList(c *gin.Context) {
	tokenList, err := this.dao.GetTxTokenList()
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, tokenList)
}

// @Summary 签名监控
// @Tags Chart
// @Accept  json
// @Produce  json
// @Success 200 {object} JsonResult{data=dao.TokenListInterface}
// @Router /chart/anchorCount/list [get]
func (this *Controller) AnchorCount(c *gin.Context) {
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	tokenKey := c.Query("tokenKey")
	timeType := c.Query("timeType")
	tokenList, err := this.dao.GetTxTokenList()
	token := tokenList[tokenKey]
	anchors, err := this.dao.TokenListAnchorCount(token, startTime, endTime, timeType)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, anchors)
}

// @Summary 跨链交易数监控
// @Tags Chart
// @Accept  json
// @Produce  json
// @Success 200 {object} JsonResult{data=dao.TokenListInterface}
// @Router /chart/crossTxCount/list [get]
func (this *Controller) CrossTxCount(c *gin.Context) {
	tokenList, err := this.dao.GetTxTokenList()
	for key, value := range tokenList {
		count := this.dao.TokenListCount(value)
		value.Count = count
		tokenList[key] = value
	}

	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, tokenList)
}
