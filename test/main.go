package main

import (
	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/crosschain/polynetwork"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"log"
)

func main() {
	deployer := &Deployer{
		PrivateKey:    "e887faa2e702daa055e59ff9f94d2af9ded1b308fc30935bbc1b63dabbfb8b11",
		Host:          "https://polynetworkcc3dcb2-5-api.dev.z7a.xyz",
		ProxyPath:     "../contracts/ZilCrossChainManagerProxy.scilla",
		ImplPath:      "../contracts/ZilCrossChainManager.scilla",
		LockProxyPath: "../contracts/LockProxy.scilla",
	}
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(deployer.PrivateKey)
	client := provider.NewProvider(deployer.Host)
	proxy, impl, lockProxy, err := deployer.Deploy(wallet, client)
	log.Printf("lock proxy address: %s\n", lockProxy)
	if err != nil {
		log.Fatalln(err.Error())
	}

	p := &polynetwork.Proxy{
		ProxyAddr:  proxy,
		ImplAddr:   impl,
		Wallet:     wallet,
		Client:     client,
		ChainId:    chainID,
		MsgVersion: msgVersion,
	}

	_, err1 := p.UpgradeTo()
	if err1 != nil {
		log.Fatalln(err1.Error())
	}

	_, err2 := p.Unpause()
	if err2 != nil {
		log.Fatalln(err2.Error())
	}

	tester := &Tester{p: p}
	tester.InitGenesisBlock()
	//tester.ChangeBookKeeper()
	tester.VerifierHeaderAndExecuteTx()

}
