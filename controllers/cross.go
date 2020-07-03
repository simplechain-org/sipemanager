package controllers

import (
	"sipemanager/dao"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

func (this *Controller) EventLog(logs []types.Log, abiParsed abi.ABI, node dao.InstanceNodes) {
	makerTx := abiParsed.Events["MakerTx"].ID.Hex()
	takerTx := abiParsed.Events["TakerTx"].ID.Hex()
	makerFinish := abiParsed.Events["MakerFinish"].ID.Hex()
	for _, event := range logs {
		var item dao.CrossEvents
		switch event.Topics[0].Hex() {
		case makerTx:
			var args CrossMakerTx
			err := abiParsed.Unpack(&args, "MakerTx", event.Data)
			if err != nil {
				logrus.Error("makerTx:", err.Error())
			}
			item = dao.CrossEvents{
				BlockNumber:     event.BlockNumber,
				TxId:            event.Topics[1].Hex(),
				Event:           "MakerTx",
				From:            "0x" + event.Topics[2].Hex()[26:],
				NetworkId:       node.NetworkId,
				RemoteNetworkId: args.RemoteChainId.Int64(),
				Value:           args.Value.String(),
				DestValue:       args.DestValue.String(),
				TransactionHash: event.TxHash.Hex(),
				CrossAddress:    node.CrossAddress,
				ChainId:         node.ChainId,
			}
			this.dao.MakerEventUpsert(item)
		case takerTx:
			var args CrossTakerTx
			err := abiParsed.Unpack(&args, "TakerTx", event.Data)
			if err != nil {
				logrus.Error("takerTx:", err.Error())
			}
			item = dao.CrossEvents{
				BlockNumber:     event.BlockNumber,
				TxId:            event.Topics[1].Hex(),
				Event:           "TakerTx",
				From:            strings.ToLower(args.From.Hex()),
				To:              "0x" + event.Topics[2].Hex()[26:],
				NetworkId:       node.NetworkId,
				RemoteNetworkId: args.RemoteChainId.Int64(),
				Value:           args.Value.String(),
				DestValue:       args.DestValue.String(),
				TransactionHash: event.TxHash.Hex(),
				CrossAddress:    node.CrossAddress,
				ChainId:         node.ChainId,
			}
			this.dao.TakerEventUpsert(item)
		case makerFinish:
			item = dao.CrossEvents{
				BlockNumber:     event.BlockNumber,
				TxId:            event.Topics[1].Hex(),
				Event:           "MakerFinish",
				To:              "0x" + event.Topics[2].Hex()[26:],
				NetworkId:       node.NetworkId,
				TransactionHash: event.TxHash.Hex(),
				CrossAddress:    node.CrossAddress,
				ChainId:         node.ChainId,
			}
			this.dao.MakerFinishEventUpsert(item)
		}
	}
}
