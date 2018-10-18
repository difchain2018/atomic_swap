// // Copyright (c) 2017 The Decred developers
// // Use of this source code is governed by an ISC
// // license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

type ethRPC struct {
	RPC       *RpcInfo
	rpcClient *ethrpc.EthRPC
}

func createClient() *ethRPC {
	matcher := new(ethRPC)
	Register(matcher.ID(), matcher)
	matcher.InitRPC(null)
	return matcher
}

func (m *ethRPC) ID() int {
	return 1002
}
func (m *ethRPC) Name() string {
	return "ETH"
}

func (m *ethRPC) MainFee() float64 {
	return 0.0001
}
func (m *ethRPC) MainAddress() string {
	return "2MyAYGXXdJBWrfU16f8HEy7Gmu5ZDEZVoXA"
}

func (m *ethRPC) MainMax() float64 {
	return 999999999999999999
}

func (m *ethRPC) InitRPC(rpcInfo *RpcInfo) error {
	m.RPC = rpcInfo
	url := "http://127.0.0.1:8080"

	client := ethrpc.New(url)

	version, err := client.Web3ClientVersion()
	if err != nil {
		return err
	}

	m.rpcClient = client
	return nil
}

func (m *ethRPC) Release() {
}

// func (m *ethRPC) GenerateAddress(account string) (string, string, string, error) {
// 	return keystore.NewAccountEx(account)
// }

func (m *ethRPC) UnlockWallet() error {
	return nil
}

func (m *ethRPC) SendTrans(account string, toaddress string, amount float64, coin string, tag string) (string, error) {

	balance, err := m.rpcClient.EthGetBalance(address, "latest")

	if err != nil {
		return "0.0", err
	}
	//todo: 取小数点8位，未解决，目前只取了6位
	//查询结果除以1000000000000000000得到以太坊余额、取小数点后8位
	//fbalance ,_:= strconv.ParseFloat( balance.String(),64)
	//bigbalance := big.NewFloat( fbalance )
	//resultF64 ,_:= bigbalance.Quo( bigbalance, ether).Float64()
	//resultStr := fmt.Sprintf("%f", resultF64)

	return balance.String(), nil
}

func (m *ethRPC) GetBalance(account string, coin string) (amount float64, err error) {

	balance, err := m.rpcClient.EthGetBalance(address, "latest")

	if err != nil {
		return "0.0", err
	}
	//todo: 取小数点8位，未解决，目前只取了6位
	//查询结果除以1000000000000000000得到以太坊余额、取小数点后8位
	//fbalance ,_:= strconv.ParseFloat( balance.String(),64)
	//bigbalance := big.NewFloat( fbalance )
	//resultF64 ,_:= bigbalance.Quo( bigbalance, ether).Float64()
	//resultStr := fmt.Sprintf("%f", resultF64)

	return 0, nil
}

func main() {

	c := createClient()

	fmt.Errorf("main: start %v", c)
}
