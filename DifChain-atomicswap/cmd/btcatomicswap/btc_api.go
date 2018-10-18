// // Copyright (c) 2017 The Decred developers
// // Use of this source code is governed by an ISC
// // license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"net"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
)

var (
	chainParams = &chaincfg.MainNetParams
)

var (
	flagset = flag.NewFlagSet("", flag.ExitOnError)
	//connectFlag = flagset.String("s", "localhost", "host[:port] of Bitcoin Core wallet RPC server")
	connectFlag = flagset.String("s", "10.220.10.104:19001", "host[:port] of Bitcoin Core wallet RPC server")
	rpcuserFlag = flagset.String("rpcuser", "admin1", "username for wallet RPC authentication")
	rpcpassFlag = flagset.String("rpcpass", "123", "password for wallet RPC authentication")
	//testnetFlag = flagset.Bool("testnet", false, "use testnet network")
	testnetFlag = flagset.Bool("testnet", true, "use testnet network")

	// connectFlag = flagset.String("s", "192.168.5.229:8332", "host[:port] of Bitcoin Core wallet RPC server")
	// rpcuserFlag = flagset.String("rpcuser", "bitcoin", "username for wallet RPC authentication")
	// rpcpassFlag = flagset.String("rpcpass", "local321", "password for wallet RPC authentication")
	// //testnetFlag = flagset.Bool("testnet", false, "use testnet network")
	// testnetFlag = flagset.Bool("testnet", false, "use testnet network")
)

func normalizeAddress(addr string, defaultPort string) (hostport string, err error) {
	host, port, origErr := net.SplitHostPort(addr)
	if origErr == nil {
		return net.JoinHostPort(host, port), nil
	}
	addr = net.JoinHostPort(addr, defaultPort)
	_, _, err = net.SplitHostPort(addr)
	if err != nil {
		return "", origErr
	}
	return addr, nil
}

func walletPort(params *chaincfg.Params) string {
	switch params {
	case &chaincfg.MainNetParams:
		return "8332"
	case &chaincfg.TestNet3Params:
		return "18332"
	default:
		return ""
	}
}

func GetConnect() (c *rpc.Client, err error) {
	connect, err := normalizeAddress(*connectFlag, walletPort(chainParams))
	if err != nil {
		return nil, err
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         *rpcuserFlag,
		Pass:         *rpcpassFlag,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	return rpc.New(connConfig, nil)
}

func ReleaseConnect(c *rpc.Client) {
	if c == nil {
		return
	}

	c.Shutdown()
	c.WaitForShutdown()
}

func SetAccount(address string, account string, c *rpc.Client) (err error) {

	addr, err := btcutil.DecodeAddress(address, chainParams)
	if err != nil {
		return err
	}

	rpcClient := c
	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return err
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	err = rpcClient.SetAccount(addr, account)
	return err
}

func CreateAddress(account string, c *rpc.Client) (address string, privatekey string, err error) {

	rpcClient := c
	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return "", "", err
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	addr, err := rpcClient.GetNewAddress(account)

	if err != nil {
		return "", "", err
	}

	privkey, err := rpcClient.DumpPrivKey(addr)
	if err != nil {
		return "", "", err
	}

	return addr.EncodeAddress(), privkey.String(), err
}

func GetBalenceByAccount(account string, c *rpc.Client) (amount float64, err error) {

	rpcClient := c
	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return 0, err
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	balence, err := rpcClient.GetBalance(account)

	if err != nil {
		return 0, err
	}

	return balence.ToBTC(), err
}

func ListAccount(c *rpc.Client) (acc map[string]btcutil.Amount, err error) {

	rpcClient := c
	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return nil, err
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	return rpcClient.ListAccounts()
}

func GetAccountByAddress(address string, c *rpc.Client) (acc string, err error) {
	addr, err := btcutil.DecodeAddress(address, chainParams)
	if err != nil {
		return "", err
	}

	rpcClient := c
	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return "", err
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	return rpcClient.GetAccount(addr)
}

func GetAddressByAccount(account string, c *rpc.Client) (address string, err error) {

	rpcClient := c
	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return "", err
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	addr, err := rpcClient.GetAddressesByAccount(account)
	if err != nil {
		return "", err
	}

	fmt.Printf("GetAddressByAccount acc=%v addrs=%v", account, addr)
	return addr[0].EncodeAddress(), err
}

func ListTransaction(account string, c *rpc.Client) {
	rpcClient := c
	var err error = nil

	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	arr, err := rpcClient.ListTransactions(account)

	if err == nil {
		fmt.Printf("ListTransaction %v", arr)
	}
}

func UnlockWallet(walletpwd string, timeout_second int64, c *rpc.Client) (err error) {
	rpcClient := c

	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return err
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	return rpcClient.WalletPassphrase(walletpwd, timeout_second)
}

func SetWalletPwd(walletpwdold string, walletpwdnew string, c *rpc.Client) (err error) {

	rpcClient := c

	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return err
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	return rpcClient.WalletPassphraseChange(walletpwdold, walletpwdnew)
}

func SendFromAccount(Account string, toaddress string, amount float64, c *rpc.Client) (txHs string, err error) {

	addrto, err := btcutil.DecodeAddress(toaddress, chainParams)
	if err != nil {
		return "", err
	}

	amountbtc, err := btcutil.NewAmount(amount)
	if err != nil {
		return "", err
	}

	rpcClient := c
	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return "", err
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	tx, err := rpcClient.SendFrom(Account, addrto, amountbtc)
	if err != nil {
		return "", err
	}

	return tx.String(), err
}

func GetSendResult(txHs string, c *rpc.Client) (amount float64, fee float64, err error) {
	rpcClient := c

	hasTx, err := chainhash.NewHashFromStr(txHs)
	if err != nil {
		return 0, 0, err
	}

	if c == nil {
		rpcClient, err = GetConnect()
		if err != nil {
			return 0, 0, err
		}
		defer func() { ReleaseConnect(rpcClient) }()
	}

	rst, err := rpcClient.GetTransaction(hasTx)

	amount = 0
	fee = 0

	if err == nil {
		if rst.Confirmations < 6 {
			return 0, 0, errors.New("confirmations is too little")
		}
		amount = rst.Details[0].Amount
		fee = rst.Fee
	}

	return amount, fee, err
}

func testAPI(c *rpc.Client) {

	//SetAccount("", "", nil)
	err := UnlockWallet("123", 3600, c)

	//v, _ := GetAddressByAccount("wallet_create_by_shawn", c)

	// v2, _ := GetBalenceByAddress("2MyTd5vt3PVpCHsYSGeyjULL8QWkykMkw5U", c)
	// v3, _ := GetBalenceByAccount("wallet_create_by_shawn", c)

	// fmt.Printf("GetBalenceByAddress: %v\n", v2)
	// fmt.Printf("GetBalenceByAccount: %v\n", v3)

	// ListTransaction("wallet_create_by_shawn", c)

	// return

	// //ListTransaction("wallet_create_by_shawn", c)

	//	bc, _ := GetBalenceByAccount("wallet_create_by_shawn", c)

	//	fmt.Printf("GetBalenceByAccount %v\n", bc)
	//err := SendFromAccount("test2", "2MyTd5vt3PVpCHsYSGeyjULL8QWkykMkw5U", 10, c)

	txHs, err := SendFromAccount("wallet_create_by_shawn", "2NBvd43tvdFUziVUEAcwfdVSC7ZVAiTLYhZ", 10, c) //SendFromAccount("wallet_create_by_shawn", "2NBvd43tvdFUziVUEAcwfdVSC7ZVAiTLYhZ", 1, c)
	//SendFromAccount("test2", "2MyTd5vt3PVpCHsYSGeyjULL8QWkykMkw5U", 10, c) //

	if err != nil {
		fmt.Printf("SendFromAccount: %v", err)
		return
	}

	_, _, err = GetSendResult(txHs, c)

	fmt.Printf("SendFromAccount: %v\n", err)

	// err = SendFromAccount("wallet_create_by_shawn", "2NBvd43tvdFUziVUEAcwfdVSC7ZVAiTLYhZ", 2, c)

	// fmt.Printf("SendFromAccount: %v\n", err)

	// // address, err := GetAddressByAccount("test2", c)
	// // if err != nil {
	// // 	fmt.Printf("GetAddressByAccount: %v", err)
	// // 	return
	// // }

	// // fmt.Printf("GetAddressByAccount %v", address)

	// // acc, err := GetAccountByAddress("2N6BytQyKz1PJoWpBCso1vY1ZAiY3fmxf1k", c)
	// // if err != nil {
	// // 	fmt.Printf("GetAccountByAddress: %v", err)
	// // 	return
	// // }

	// // fmt.Printf("GetAccountByAddress %v", acc)

	// // addr, err := CreateAddress("test2")
	// // if err != nil {
	// // 	fmt.Printf("NewAddress: %v", err)
	// // 	return
	// // }

	// // fmt.Printf("NewAddress %v", addr)

	// // err = SendFromAccount("wallet_create_by_shawn", "2NBvd43tvdFUziVUEAcwfdVSC7ZVAiTLYhZ", 5)
	// // if err != nil {
	// // 	fmt.Printf("SendFromAccount: %v", err)
	// // 	return
	// // }

	// // fmt.Printf("SendFromAccount succeed")

	// // amount, err := GetBalenceByAddress("2MyTd5vt3PVpCHsYSGeyjULL8QWkykMkw5U")

	// // if err != nil {
	// // 	fmt.Printf("GetBalenceByAddress: %v", err)
	// // 	return
	// // }

	// // fmt.Printf("GetBalenceByAddress: %v", amount)

	// // accounts, err := ListAccount()
	// // if err != nil {
	// // 	fmt.Printf("GetBalenceByAddress: %v", err)
	// // 	return
	// // }

	// // fmt.Printf("ListAccount: %v", accounts)
}

// func main() {
// 	c, err := GetConnect()
// 	if err != nil {
// 		fmt.Errorf("GetConnect: %v", err)
// 		return
// 	}

// 	testAPI(c)
// 	ReleaseConnect(c)
// }
