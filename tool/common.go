package tool

import (
	"fmt"
	"time"

	sdk "github.com/polynetwork/poly-go-sdk"
	"github.com/polynetwork/poly/common/password"
)

func getAccountByPassword(t *Tool, path string) (*sdk.Account, error) {
	wallet, err := t.sdk.OpenWallet(path)
	if err != nil {
		return nil, fmt.Errorf("open wallet error:%s", err)
	}
	pwd, err := password.GetPassword()
	if err != nil {
		return nil, fmt.Errorf("getPassword error:%s", err)
	}
	user, err := wallet.GetDefaultAccount(pwd)
	if err != nil {
		return nil, fmt.Errorf("getDefaultAccount error:%s", err)
	}
	return user, nil
}

func waitForBlock(t *Tool) error {
	_, err := t.sdk.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		return fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return nil
}

func ConcatKey(args ...[]byte) []byte {
	temp := []byte{}
	for _, arg := range args {
		temp = append(temp, arg...)
	}
	return temp
}
