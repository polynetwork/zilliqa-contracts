# Zilliqa corss chain smart contract

# Table of Content

- [Overview](#overview)
- [ZilCrossChainManager Contract Specification](#zilcrosschainmanager-contract-specification)
- [ZilCrossChainManagerProxy Contract Specification](#zilcrosschainmanagerproxy-contract-specification)
- [LockProxy Contract Specification](#lockproxy-contract-specification)

# Overview

The table blow summarizes the purpose of the contracts that polynetwork will use:

| Contract Name | File and Location | Description |
|--|--| --|
|ZilCrossChainManager| [`ZilCrossChainManager.scilla`](./contracts/ZilCrossChainManager.scilla)  | The main contract that keeps track of the book keepers of Poly chain, push cross chain transaction event to relayer and execute the cross chain transaction from Poly chain to Zilliqa.|
|ZilCrossChainManagerProxy| [`ZilCrossChainManagerProxy.scilla`](./contracts/ZilCrossChainManagerProxy.scilla)  | A proxy contract that sits on top of the ZilCrossChainManager contract. Any call to the `ZilCrossChainManager` contract must come from `ZilCrossChainManagerProxy`. This contract facilitates upgradeability of the `ZilCrossChainManager` contract in case a bug is found.|
|LockProxy| [`LockProxy.scilla`](./contracts/LockProxy.scilla)  | A application contract that allows people to lock ZRC2 tokens and native zils to get corresponding tokens in target chain (e.g. ERC20 in ethereum) and vise versa.|

# ZilCrossChainManager Contract Specification

# ZilCrossChainManagerProxy Contract Specification

# LockProxy Contract Specification