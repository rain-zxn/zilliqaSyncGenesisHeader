package main

import (
	"flag"
	sdk "github.com/polynetwork/poly-go-sdk"
	"time"
	"zilliqaSyncGenesisHeader/common"
	"zilliqaSyncGenesisHeader/log"
	"zilliqaSyncGenesisHeader/tool"
)

var (
	Config string //config file
	Method string //method list in cmdline
)

func init() {
	flag.StringVar(&Config, "cfg", "./config.json", "Config of poly-validator-tool")
	flag.StringVar(&Method, "m", "", "method to run")
	flag.Parse()
}

func main() {
	log.InitLog(log.InfoLog, "./Logs/", log.Stdout)
	defer time.Sleep(time.Second)

	err := common.DefConfig.Init(Config)
	if err != nil {
		log.Error("DefConfig.Init error:%s", err)
		return
	}

	polySdk := sdk.NewPolySdk()
	polySdk.NewRpcClient().SetAddress(common.DefConfig.JsonRpcAddress)
	hdr, err := polySdk.GetHeaderByHeight(0)
	if err != nil {
		log.Error("Failed to initialize poly chain id", "err", err)
		return
	}
	polySdk.SetChainId(hdr.ChainID)

	method := Method
	tool := tool.NewTool()
	tool.SetPolySdk(polySdk)
	tool.RegMethods()
	tool.Start(method)
}
