package tool

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// only testnet
func GetZilGenesisHeader(t *Tool) error {
	data, err := ioutil.ReadFile("./params/GetZilGenesisHeader_testnet.json")
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile failed %v", err)
	}
	getZilGenesisHeaderParam := new(GetZilGenesisHeaderParam)
	err = json.Unmarshal(data, getZilGenesisHeaderParam)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed %v", err)
	}
	type TxBlockAndDsComm struct {
		TxBlock *core.TxBlock
		DsBlock *core.DsBlock
		DsComm  []core.PairOfNode
	}

	zilSdk := provider.NewProvider(getZilGenesisHeaderParam.ZilliqaRPC)
	// ON TESTNET it gets the currentDScomm. The getMiner info returns the an empty dscommittee
	// for a previous DSBlock num
	initDsComm, err := zilSdk.GetCurrentDSComm()
	if err != nil {
		return fmt.Errorf("zilSdk.GetCurrentDSComm failed: %v", err)
	}
	// as its name suggest, the tx epoch is actually a future tx block
	// zilliqa side has this limitation to avoid some risk that no tx block got mined yet
	nextTxEpoch, err := strconv.ParseUint(initDsComm.CurrentTxEpoch, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseUintm failed: %v", err)
	}
	fmt.Printf("next tx epoch is %v, current tx block number is %s, ds block number is %s, number of ds guard is: %d\n", nextTxEpoch, initDsComm.CurrentTxEpoch, initDsComm.CurrentDSEpoch, initDsComm.NumOfDSGuard)

	for {
		latestTxBlock, err := zilSdk.GetLatestTxBlock()
		if err != nil {
			return fmt.Errorf("zilSdk.GetLatestTxBlock failed: %v", err)
		}
		fmt.Println("wait current tx block got generated")
		latestTxBlockNum, err := strconv.ParseUint(latestTxBlock.Header.BlockNum, 10, 64)
		if err != nil {
			return fmt.Errorf("strconv.ParseUint BlockNum failed: %v", err)
		}
		fmt.Printf("latest tx block num is: %d, next tx epoch num is: %d\n", latestTxBlockNum, nextTxEpoch)
		if latestTxBlockNum >= nextTxEpoch {
			break
		}
		time.Sleep(time.Second * 20)
	}

	var dsComm []core.PairOfNode
	for _, ds := range initDsComm.DSComm {
		dsComm = append(dsComm, core.PairOfNode{
			PubKey: ds,
		})
	}
	dsBlockT, err := zilSdk.GetDsBlockVerbose(initDsComm.CurrentDSEpoch)
	if err != nil {
		return fmt.Errorf("zilSdk.GetDsBlockVerbose get ds block %s failed: %v", initDsComm.CurrentDSEpoch, err)
	}
	dsBlock := core.NewDsBlockFromDsBlockT(dsBlockT)
	txBlockT, err := zilSdk.GetTxBlockVerbose(initDsComm.CurrentTxEpoch)
	if err != nil {
		return fmt.Errorf("zilSdk.GetTxBlockVerbose get tx block %s failed: %v", initDsComm.CurrentTxEpoch, err)
	}

	txBlock := core.NewTxBlockFromTxBlockT(txBlockT)

	txBlockAndDsComm := TxBlockAndDsComm{
		TxBlock: txBlock,
		DsBlock: dsBlock,
		DsComm:  dsComm,
	}

	raw, err := json.Marshal(txBlockAndDsComm)
	if err != nil {
		return fmt.Errorf("json.Marshal txBlockAndDsComm failed: %v", err)
	}
	tx, err := t.sdk.Native.Hs.NewSyncGenesisHeaderTransaction(getZilGenesisHeaderParam.ZilliqaChainId, raw)
	if err != nil {
		return fmt.Errorf("NewSyncGenesisHeaderTransaction failed: %v", err)
	}
	txString := hex.EncodeToString(tx.ToArray())
	if err != nil {
		return fmt.Errorf("hex.DecodeString sink error: %v", err)
	}
	file, err := os.OpenFile("sigDataIn.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("open file sigDataIn.txt err: %v", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(txString))
	if err != nil {
		return fmt.Errorf("write file sigDataIn.txt err: %v", err)
	}
	fmt.Println("success GetZilGenesisHeader, write sigDataIn.txt")
	return nil
}
