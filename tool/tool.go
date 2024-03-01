package tool

import (
	"fmt"
	"zilliqaSyncGenesisHeader/log"

	sdk "github.com/polynetwork/poly-go-sdk"
)

type Method func(t *Tool) error

type Tool struct {
	methodMap map[string]Method
	sdk       *sdk.PolySdk
}

func NewTool() *Tool {
	return &Tool{
		methodMap: make(map[string]Method, 0),
	}
}

func (this *Tool) RegMethod(name string, method Method) {
	this.methodMap[name] = method
}

func (this *Tool) Start(name string) {
	method, err := this.getMethodByName(name)
	if err != nil {
		log.Errorf("getMethodByName", err)
	}
	err = method(this)
	if err != nil {
		log.Errorf("method %s failed: %s", name, err)
	}
}

func (this *Tool) SetPolySdk(sdk *sdk.PolySdk) {
	this.sdk = sdk
}

func (this *Tool) getMethodByName(name string) (Method, error) {
	method, ok := this.methodMap[name]
	if !ok {
		return nil, fmt.Errorf("method name not exist")
	}
	return method, nil
}
