package main

import (
	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"log"
)

func main() {
	deployer := &Deployer{
		PrivateKey: "e887faa2e702daa055e59ff9f94d2af9ded1b308fc30935bbc1b63dabbfb8b11",
		Host:       "http://54.213.251.102:5555",
		ProxyPath:  "../contracts/ZilCrossChainManagerProxy.scilla",
		ImplPath:   "../contracts/ZilCrossChainManager.scilla",
	}
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(deployer.PrivateKey)
	client := provider.NewProvider(deployer.Host)
	proxy, impl, err := deployer.Deploy(wallet, client)
	if err != nil {
		log.Fatalln(err.Error())
	}

	p := &Proxy{
		ProxyAddr: proxy,
		ImplAddr:  impl,
		Wallet:    wallet,
		Client:    client,
	}

	err1 := p.UpgradeTo()
	if err1 != nil {
		log.Fatalln(err1.Error())
	}

	err2 := p.Unpause()
	if err2 != nil {
		log.Fatalln(err2.Error())
	}

}
