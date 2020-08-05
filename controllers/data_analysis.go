package controllers

import (
	"github.com/simplechain-org/go-simplechain/common"
	"math/big"
	"sipemanager/blockchain"
	"strconv"
	"strings"

	"sipemanager/dao"
	"sipemanager/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AnchorsNodes struct {
	Id      uint
	Name    string
	Address string
}

func (this *Controller) AnalysisAnchors() {
	registers, err := this.dao.ListChainRegisterByStatus(1)
	if err != nil {
		logrus.Error(utils.ErrLogCode{LogType: "controller => data_analysis => AnalysisAnchors:", Code: 40002, Message: "Analysis Anchors Failed", Err: nil})
	}
	for _, register := range registers {
		sourceChain, err := this.dao.GetChain(register.SourceChainId)
		targetChain, err := this.dao.GetChain(register.TargetChainId)
		if err != nil {
			logrus.Error(utils.ErrLogCode{LogType: "controller => data_analysis => AnalysisAnchors:", Code: 40002, Message: "GetChain Anchors  Not Found", Err: nil})
		}
		anchorIds := strings.Split(register.AnchorAddresses, ",")
		for _, anchorId := range anchorIds {
			n, _ := strconv.Atoi(anchorId)
			anchor, err := this.dao.GetAnchorNode(uint(n))
			if err != nil {
				logrus.Error(utils.ErrLogCode{LogType: "controller => data_analysis => AnalysisAnchors:", Code: 40002, Message: "GetAnchorNode Anchors Not Found", Err: nil})
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
				logrus.Error(utils.ErrLogCode{LogType: "controller => data_analysis => AnalysisAnchors:", Code: 40001, Message: "Analysis Anchors Failed", Err: nil})
				continue
			}
		}
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
// @Param startTime query string true "hour:2020-07-10 12:00:00 day:2020-07-10 week:202025"
// @Param endTime query string true "hour:2020-07-12 12:00:00 day:2020-07-12 week:202029"
// @Param tokenKey query string true "1,2"
// @Param timeType query string true "hour,day,week"
// @Success 200 {object} JsonResult{data=dao.TokenListCount}
// @Router /chart/anchorCount/list [get]
func (this *Controller) AnchorCount(c *gin.Context) {
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	tokenKey := c.Query("tokenKey")
	timeType := c.Query("timeType")
	tokenList, err := this.dao.GetTxTokenList()
	token := tokenList[tokenKey]
	anchorIds := strings.Split(token.AnchorAddresses, ",")
	tokenCount := make(map[string][]dao.TokenListCount, 0)
	for _, id := range anchorIds {
		n, _ := strconv.Atoi(id)
		anchors, err := this.dao.TokenListAnchorCount(token, startTime, endTime, timeType, uint(n))
		if err != nil {
			this.echoError(c, err)
			return
		}
		tokenCount[id] = anchors
	}
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, tokenCount)
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

type FinishEventView struct {
	FinishEventList []FinishEvent
	Count           uint32
}

type FinishEvent struct {
	dao.CrossAnchors
	TokenName    string
	TokenListKey string
	AnchorId     uint
	AnchorName   string
}

// @Summary MakeFinish手续费记录
// @Tags Chart
// @Accept  json
// @Produce  json
// @Param startTime query string false "时间戳"
// @Param endTime query string false "时间戳"
// @Param anchorId query string false "锚定节点ID"
// @Param page query string true "页码"
// @Param limit query string true "页数"
// @Success 200 {object} JsonResult{data=FinishEventView}
// @Router /chart/finishList/list [get]
func (this *Controller) getFinishList(c *gin.Context) {
	pageParam := c.Query("page")
	limitParam := c.Query("limit")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	anchorId := c.Query("anchorId")
	page, err := strconv.Atoi(pageParam)
	limit, err := strconv.Atoi(limitParam)
	offset := (page - 1) * limit
	var finishList []dao.CrossAnchors
	var count uint32
	finishList, count, err = this.dao.QueryFinishList(uint32(offset), uint32(limit), startTime, endTime, anchorId)
	if err != nil {
		this.echoError(c, err)
		return
	}

	tokenList, err := this.dao.GetTxTokenList()
	finishEventArr := make([]FinishEvent, 0)
	var tokenKey string

	for _, item := range finishList {
		sourceId := strconv.Itoa(int(item.ChainId))
		targetId := strconv.Itoa(int(item.RemoteChainId))
		if _, exists := tokenList[sourceId+","+targetId]; exists {
			tokenKey = sourceId + "," + targetId
		} else {
			tokenKey = targetId + "," + sourceId
		}
		anchorId, anchorName, err := this.GetAnchorId(tokenList[tokenKey], item.AnchorAddress)
		if err != nil {
			this.echoError(c, err)
			return
		}
		finishEvent := FinishEvent{
			CrossAnchors: item,
			TokenName:    tokenList[tokenKey].Name,
			TokenListKey: tokenKey,
			AnchorId:     anchorId,
			AnchorName:   anchorName,
		}
		finishEventArr = append(finishEventArr, finishEvent)

	}
	result := FinishEventView{
		FinishEventList: finishEventArr,
		Count:           count,
	}

	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, result)
}

func (this *Controller) GetAnchorId(token dao.TokenListInterface, anchorAddress string) (uint, string, error) {
	anchorIds := strings.Split(token.AnchorAddresses, ",")
	for _, id := range anchorIds {
		n, _ := strconv.Atoi(id)
		anchor, anchorErr := this.dao.GetAnchorNode(uint(n))
		if anchorErr != nil {
			return 0, "", anchorErr
		}
		if anchor.Address == anchorAddress {
			return anchor.ID, anchor.Name, nil
		}
	}
	return 0, "nknown", nil
}

type AnchorNodeMonitor struct {
	SourceBalance *big.Int
	TargetBalance *big.Int
	SignTxCount   uint32
	AnchorId      uint
	AnchorName    string
	UnSignTxCount uint32
	OnLine        bool
}

// @Summary 锚定节点监控
// @Tags Chart
// @Accept  json
// @Produce  json
// @Param tokenKey query string true "1,2"
// @Success 200 {object} JsonResult{data=FinishEventView}
// @Router /chart/crossMonitor/list [get]
func (this *Controller) GetCrossMonitor(c *gin.Context) {
	tokenKey := c.Query("tokenKey")
	tokenList, err := this.dao.GetTxTokenList()
	token := tokenList[tokenKey]
	sourceNode, err := this.dao.GetNodeByChainId(token.ChainID)
	targetNode, err := this.dao.GetNodeByChainId(token.RemoteChainID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	source := &blockchain.Node{
		Address: sourceNode.Address,
		Port:    sourceNode.Port,
		ChainId: sourceNode.ChainId,
		IsHttps: sourceNode.IsHttps,
	}
	target := &blockchain.Node{
		Address: targetNode.Address,
		Port:    targetNode.Port,
		ChainId: targetNode.ChainId,
		IsHttps: targetNode.IsHttps,
	}
	sourceApi, err := blockchain.NewApi(source)
	targetApi, err := blockchain.NewApi(target)
	anchorIds := strings.Split(token.AnchorAddresses, ",")
	AnchorMon := make([]AnchorNodeMonitor, 0)

	MonCountMap, err := this.QueryMonitorBy(token)
	for _, anchorId := range anchorIds {
		anId, _ := strconv.Atoi(anchorId)
		anchor, err := this.dao.GetAnchorNode(uint(anId))
		souBal, err := sourceApi.LatestBalanceAt(common.HexToAddress(anchor.Address))
		tarBal, err := targetApi.LatestBalanceAt(common.HexToAddress(anchor.Address))
		if err != nil {
			this.echoError(c, err)
			return
		}
		var online bool
		if 200-MonCountMap[anchor.Address] == 200 {
			online = false
		} else {
			online = true
		}
		AM := AnchorNodeMonitor{
			SourceBalance: souBal,
			TargetBalance: tarBal,
			SignTxCount:   MonCountMap[anchor.Address],
			UnSignTxCount: 200 - MonCountMap[anchor.Address],
			AnchorId:      uint(anId),
			AnchorName:    anchor.Name,
			OnLine:        online,
		}
		AnchorMon = append(AnchorMon, AM)
	}

	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, AnchorMon)

}

type MonitorCount struct {
	Address string
	Count   uint32
}

func (this *Controller) QueryMonitorBy(token dao.TokenListInterface) (map[string]uint32, error) {
	anchorIds := strings.Split(token.AnchorAddresses, ",")
	MonCountMap := make(map[string]uint32, 0)
	for _, anchorId := range anchorIds {
		anId, _ := strconv.Atoi(anchorId)
		anchor, err := this.dao.GetAnchorNode(uint(anId))
		SApi, err := blockchain.NewDirectApi(anchor.SourceRpcUrl)
		TApi, err := blockchain.NewDirectApi(anchor.TargetRpcUrl)
		if SApi == nil || TApi == nil {
			return MonCountMap, nil
		}
		sourceMon, err := SApi.GetMonitor()
		targetMon, err := TApi.GetMonitor()

		for key, value := range sourceMon.Recently {
			if sourceMon.Tally[key] > 100 {
				for tKey, tValue := range targetMon.Recently {
					if sourceMon.Tally[tKey] > 100 {
						if key.Hex() == tKey.Hex() {
							MonCountMap[strings.ToLower(key.Hex())] = value + tValue
						}
					}
				}
			}
		}
		if err != nil {
			return MonCountMap, err
		}
	}
	return MonCountMap, nil
}
