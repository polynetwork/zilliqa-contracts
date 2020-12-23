package main

import (
	"fmt"
	"log"
)

func main() {
	deployer := &Deployer{
		PrivateKey: "e887faa2e702daa055e59ff9f94d2af9ded1b308fc30935bbc1b63dabbfb8b11",
		Host:       "http://54.213.251.102:5555",
		ProxyPath:  "../contracts/ZilCrossChainManagerProxy.scilla",
		ImplPath:   "../contracts/ZilCrossChainManager.scilla",
	}

	proxy, impl, err := deployer.Deploy()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("proxy = ", proxy)
	fmt.Println("impl = ", impl)

}
