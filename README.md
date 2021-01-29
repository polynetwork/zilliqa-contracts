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

The `ZilCrossChainManager` contract is the main contract of the cross chain infrastructure between Zilliqa and Poly chain.

## Roles and Privileges

The table below describes the roles and privileges that this contract defines:

| Role | Description & Privileges|                                    
| --------------- | ------------------------------------------------- |
| `admin`         | The administrator of the contract.  `admin` is a multisig wallet contract (i.e., an instance of `Wallet`).    |
| `book keepers`         | Book keepers of Poly chain which can submit cross chain transactions from Poly chain to Zilliqa|

## Data Types

The contract defines and uses several custom ADTs that we describe below:

1. Error Data Type:

```ocaml
type Error =
  | ContractFrozenFailure
  | ConPubKeysAlreadyInitialized
  | ErrorDeserializeHeader
  | NextBookersIllegal
  | SignatureVerificationFailed
  | HeaderLowerOrBookKeeperEmpty
  | InvalidMerkleProof
  | IncorrectMerkleProof
  | MerkleProofDeserializeFailed
  | AddressFormatMismatch
  | WrongTransaction
  | TransactionAlreadyExecuted
  | TransactionHashInvalid
  | AdminValidationFailed
  | ProxyValidationFailed
  | StagingAdminValidationFailed
  | StagingAdminNotExist
```

## Immutable Parameters

The table below lists the parameters that are defined at the contract deployment time and hence cannot be changed later on.

| Name | Type | Description |                                    
| ---------------      | ----------|-                                         |
| `this_chain_id`         | `Uint64` | The identifier of Zilliqa in Poly Chain. |
| `init_proxy_address` | `ByStr20` | The initial address of the `ZilCrossChainManagerProxy` contract.  |
| `init_admin`  | `ByStr20` |  The initial admin of the contract.  |

## Mutable Fields

The table below presents the mutable fields of the contract and their initial values. 

| Name        | Type       | Initial Value                           | Description                                        |
| ----------- | --------------------|--------------- | -------------------------------------------------- |
|`paused` | `ByStr20` | `True` | A flag to record the paused status of the contract. Certain transitions in the contract cannot be invoked when the contract is paused. |
| `conKeepersPublicKeyList` | `List ByStr20` |  `Nil {ByStr20}` | List of public key of consensus book Keepers. |
| `curEpochStartHeight` | `Uint32` |  `Uint32 0` | Current Epoch Start Height of Poly chain block. |
| `zilToPolyTxHashMap` | `Map Uint256 ByStr32` |  `Emp Uint256 ByStr32` | A map records transactions from Zilliqa to Poly chain. |
| `zilToPolyTxHashIndex` | Uint256` |  `Uint256 0` | Record the length of aboving map. |
| `fromChainTxExist` | Map Uint64 (Map ByStr32 Unit)` |  Emp Uint64 (Map ByStr32 Unit)` |Record the from chain txs that have been processed. |
| `contractadmin` | `ByStr20` |  `init_admin` | Address of the administrator of this contract. |




# ZilCrossChainManagerProxy Contract Specification

# LockProxy Contract Specification

# More on corss chain infrastructure

- [polynetwork](https://github.com/polynetwork/poly)
- [zilliqa-relayer](https://github.com/Zilliqa/zilliqa-relayer)