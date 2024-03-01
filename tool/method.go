package tool

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/ontio/ontology-crypto/keypair"
	poly_go_sdk_utils "github.com/polynetwork/poly-go-sdk/utils"
	"io/ioutil"
	"os"
	"zilliqaSyncGenesisHeader/log"
)

type GetZilGenesisHeaderParam struct {
	ZilliqaRPC     string
	ZilliqaChainId uint64
}

// mainnet step1
func GetZilGenesisHeaderMain(t *Tool) error {
	data, err := ioutil.ReadFile("./params/GetZilGenesisHeader.json")
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile failed %v", err)
	}
	getZilGenesisHeaderParam := new(GetZilGenesisHeaderParam)
	err = json.Unmarshal(data, getZilGenesisHeaderParam)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed %v", err)
	}
	content, err := os.ReadFile("zb_genesis_0221.json")
	if err != nil {
		log.Fatal(err)
	}
	type TxBlockAndDsComm struct {
		TxBlock *core.TxBlock
		DsBlock *core.DsBlock
		DsComm  []core.PairOfNode
	}
	txBlockAndDsComm := new(TxBlockAndDsComm)
	err = json.Unmarshal(content, txBlockAndDsComm)
	if err != nil {
		log.Fatal(err)
	}
	raw, err := json.Marshal(txBlockAndDsComm)
	if err != nil {
		return fmt.Errorf("json.Marshal txBlockAndDsComm failed: %v", err)
	}
	tx, err := t.sdk.Native.Hs.NewSyncGenesisHeaderTransaction(getZilGenesisHeaderParam.ZilliqaChainId, raw)
	if err != nil {
		return fmt.Errorf("NewSyncGenesisHeaderTransaction failed: %v", err)
	}
	fmt.Println("tx.ChainID", tx.ChainID)
	txString := hex.EncodeToString(tx.ToArray())
	if err != nil {
		return fmt.Errorf("hex.DecodeString sink error: %v", err)
	}
	file, err := os.OpenFile("sigDataIn_main_0221.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("open file sigDataIn.txt err: %v", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(txString))
	if err != nil {
		return fmt.Errorf("write file sigDataIn_main.txt err: %v", err)
	}
	fmt.Println("success GetZilGenesisHeader, write sigDataIn_main.txt")
	return nil
}

type SignatureDataParam struct {
	Path      string
	SigDataIn string
	Pubkeys   []string
	SigM      uint16
}

// mainnet step2
func SignatureData(t *Tool) error {
	data, err := ioutil.ReadFile("./params/SignatureData.json")
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile failed %v", err)
	}
	signatureDataParam := new(SignatureDataParam)
	err = json.Unmarshal(data, signatureDataParam)
	if err != nil {
		return fmt.Errorf("json.Unmarshal signatureDataParam failed %v", err)
	}
	pubKeys := make([]keypair.PublicKey, 0)
	for _, s := range signatureDataParam.Pubkeys {
		sBytes, err := hex.DecodeString(s)
		if err != nil {
			return fmt.Errorf("hex.DecodeString Pubkeys error:%s", err)
		}
		pk, err := keypair.DeserializePublicKey(sBytes)
		if err != nil {
			return fmt.Errorf("keypair.DeserializePublicKey error:%s", err)
		}
		pubKeys = append(pubKeys, pk)
	}
	tx, err := poly_go_sdk_utils.TransactionFromHexString(signatureDataParam.SigDataIn)
	if err != nil {
		return fmt.Errorf("poly_go_sdk_utils.TransactionFromHexString failed %v", err)
	}
	user, err := getAccountByPassword(t, signatureDataParam.Path)
	if err != nil {
		return fmt.Errorf("getAccountByPassword failed %v", err)
	}
	err = t.sdk.MultiSignToTransaction(tx, signatureDataParam.SigM, pubKeys, user)
	if err != nil {
		return fmt.Errorf("SignMToTransaction failed, err: %s", err)
	}
	txString := hex.EncodeToString(tx.ToArray())
	if err != nil {
		return fmt.Errorf("hex.DecodeString sink error: %v", err)
	}
	file, err := os.OpenFile("sigDataOut.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("open file sigDataOut.txt err: %v", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(txString))
	if err != nil {
		return fmt.Errorf("write file sigDataOut.txt err: %v", err)
	}
	fmt.Println("success SignatureData, write sigDataOut.txt")
	return nil
}

type SyncZilGenesisHeaderParam struct {
	Tx string
}

// mainnet step3
func SyncZilGenesisHeader(t *Tool) error {
	data, err := ioutil.ReadFile("./params/SyncZilGenesisHeader.json")
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile failed %v", err)
	}
	syncZilGenesisHeaderParam := new(SyncZilGenesisHeaderParam)
	err = json.Unmarshal(data, syncZilGenesisHeaderParam)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed %v", err)
	}
	tx, err := poly_go_sdk_utils.TransactionFromHexString(syncZilGenesisHeaderParam.Tx)
	if err != nil {
		return fmt.Errorf("poly_go_sdk_utils.TransactionFromHexString failed %v", err)
	}
	fmt.Println("tx.ChainID", tx.ChainID)
	txhash, err := t.sdk.SendTransaction(tx)
	if err != nil {
		return fmt.Errorf("t.sdk.SendTransaction failed %v", err)
	}
	fmt.Println("success send tx hash:", txhash.ToHexString())
	return nil
}
