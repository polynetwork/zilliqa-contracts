package main

import (
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"log"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/crosschain/polynetwork"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
)

func main() {
	privateKey := "e887faa2e702daa055e59ff9f94d2af9ded1b308fc30935bbc1b63dabbfb8b11"
	deployer := &Deployer{
		PrivateKey:    privateKey,
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

	l := &polynetwork.LockProxy{
		Addr:       lockProxy,
		Wallet:     wallet,
		Client:     client,
		ChainId:    chainID,
		MsgVersion: msgVersion,
	}

	tester := &Tester{p: p, l: l}
	tester.InitGenesisBlock()
	//tester.ChangeBookKeeper()
	tester.VerifierHeaderAndExecuteTx()

	// dummy ethereum contract address here
	ethLockProxy := "0x05f4a42e251f2d52b8ed15e9fedaacfcef1fad27"
	_, err3 := l.BindProxyHash("1", ethLockProxy)
	if err3 != nil {
		log.Fatalln(err3.Error())
	}

	_, err4 := l.BindAssetHash("0x0000000000000000000000000000000000000000", "1", ethLockProxy)
	if err4 != nil {
		log.Fatalln(err4.Error())
	}

	_, err5 := l.Lock("0x0000000000000000000000000000000000000000", "1", "0xd3573e0daa110b5498c54e93b66681fc0e0ff911", "100")
	if err5 != nil {
		log.Fatalln(err5.Error())
	}

	pubKey := keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true)
	address := keytools.GetAddressFromPublic(pubKey)

	_, err7 := l.SetManager("0x" + address)
	if err7 != nil {
		log.Fatalln(err7.Error())
	}

	// toAssetHash 0x05f4a42e251f2d52b8ed15e9fedaacfcef1fad27
	// toAddressHash 0xd3573e0daa110b5498c54e93b66681fc0e0ff911
	// amount 100
	// txData 0x1405f4a42e251f2d52b8ed15e9fedaacfcef1fad2714d3573e0daa110b5498c54e93b66681fc0e0ff9110000000000000000000000000000000000000000000000000000000000000064
	_, err6 := l.Unlock("0x1405f4a42e251f2d52b8ed15e9fedaacfcef1fad2714d3573e0daa110b5498c54e93b66681fc0e0ff9110000000000000000000000000000000000000000000000000000000000000064", "0x05f4a42e251f2d52b8ed15e9fedaacfcef1fad27", "1")
	if err6 != nil {
		log.Fatalln(err6.Error())
	}
}
