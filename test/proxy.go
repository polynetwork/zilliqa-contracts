package main

import (
	"errors"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"strconv"
)

type Proxy struct {
	ProxyAddr string
	ImplAddr  string
	Wallet    *account.Wallet
	Client    *provider.Provider
}

func (p *Proxy) UpgradeTo() error {
	bech32Addr, err := bech32.ToBech32Address(p.ProxyAddr)
	if err != nil {
		return err
	}

	gasPrice, err1 := p.Client.GetMinimumGasPrice()
	if err1 != nil {
		return err1
	}

	args := []core.ContractValue{
		{
			"new_crosschain_manager",
			"ByStr20",
			"0x" + p.ImplAddr,
		},
	}

	contract := contract2.Contract{
		Address: bech32Addr,
		Signer:  p.Wallet,
		Provider: p.Client,

	}

	params := contract2.CallParams{
		Version:  strconv.FormatInt(int64(util.Pack(chainID, msgVersion)), 10),
		GasPrice: gasPrice,
		GasLimit: "40000",
		Amount: "0",
		SenderPubKey: util.EncodeHex(p.Wallet.DefaultAccount.PublicKey),
	}

	tx, err2 := contract.Call("UpgradeTo", args, params, true)
	if err2 != nil {
		return err2
	}
	tx.Confirm(tx.ID, 1000, 10, p.Client)
	if tx.Status == core.Confirmed {
		return nil
	}

	return errors.New("upgrade failed")
}
