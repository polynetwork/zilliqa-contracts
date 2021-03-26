# Zilliqa cross chain smart contract

# Table of Content

- [Overview](#overview)
- [ZilCrossChainManager Contract Specification](#zilcrosschainmanager-contract-specification)
- [ZilCrossChainManagerProxy Contract Specification](#zilcrosschainmanagerproxy-contract-specification)
- [LockProxy Contract Specification](#lockproxy-contract-specification)
- [More on cross chain infrastructure](#more-on-cross-chain-infrastructure)

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
| `zilToPolyTxHashIndex` | `Uint256` |  `Uint256 0` | Record the length of aboving map. |
| `fromChainTxExist` | `Map Uint64 (Map ByStr32 Unit)` |  `Emp Uint64 (Map ByStr32 Unit)` |Record the from chain txs that have been processed. |
| `contractadmin` | `ByStr20` |  `init_admin` | Address of the administrator of this contract. |

## Transitions 

Note that some of the transitions in the `ZilCrossChainManager` contract takes `initiator` as a parameter which as explained above is the caller that calls the `ZilCrossChainManagerProxy` contract which in turn calls the `ZilCrossChainManager` contract. 

> Note: No transition in the `ZilCrossChainManager` contract can be invoked directly. Any call to the `ZilCrossChainManager` contract must come from the `ZilCrossChainManagerProxy` contract.

All the transitions in the contract can be categorized into three categories:

* **Housekeeping Transitions:** Meant to facilitate basic admin-related tasks.
* **Crosschain Transitions:** The transitions that related to cross chain tasks.

### Housekeeping Transitions

| Name        | Params     | Description | Callable when paused? | Callable when not paused? | 
| ----------- | -----------|-------------|:--------------------------:|:--------------------------:|
| `Pause` | `initiator: ByStr20`| Pause the contract temporarily to stop any critical transition from being invoked. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | :heavy_check_mark: |
| `Unpause` | `initiator: ByStr20`| Un-pause the contract to re-allow the invocation of all transitions. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.  | :heavy_check_mark: | :heavy_check_mark: |
| `UpdateAdmin` | `newAdmin: ByStr20, initiator: ByStr20` | Set a new `stagingcontractadmin` by `newAdmin`. <br>  :warning: **Note:** `initiator` must be the current `contractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: |
| `ClaimAdmin` | ` initiator: ByStr20` | Claim to be new `contract admin`. <br>  :warning: **Note:** `initiator` must be the current `stagingcontractadmin` of the contract.| :heavy_check_mark: | :heavy_check_mark: 

### Crosschain Transitions

| Name        | Params     | Description | Callable when paused? | Callable when not paused? | 
| ----------- | -----------|-------------|:--------------------------:|:--------------------------:|
| `InitGenesisBlock` | `rawHeader: ByStr, pubkeys: List Pubkey`| Sync Poly chain genesis block header to smart contrat. | <center>:x:</center> | :heavy_check_mark: |
| `ChangeBookKeeper` | `rawHeader: ByStr, pubkeys: List Pubkey, sigList: List Signature`| Change Poly chain consensus book keeper. | <center>:x:</center> | :heavy_check_mark: |
| `CrossChain` | `toChainId: Uint64, toContract: ByStr, method: ByStr, txData: ByStr`| ZRC2 token cross chain to other blockchain. this function push tx event to blockchain. | <center>:x:</center> | :heavy_check_mark: |
| `VerifyHeaderAndExecuteTx` | `proof: Proof, rawHeader: ByStr, headerProof: Proof, curRawHeader: ByStr, headerSig: List Signature`| Verify Poly chain header and proof, execute the cross chain tx  from Poly chain to Zilliqa. | <center>:x:</center> | :heavy_check_mark: |

# ZilCrossChainManagerProxy Contract Specification

`ZilCrossChainManagerProxy` contract is a relay contract that redirects calls to it to the `ZilCrossChainManager` contract.

## Roles and Privileges

The table below describes the roles and privileges that this contract defines:

| Role | Description & Privileges|                                    
| --------------- | ------------------------------------------------- |
| `init_admin`           | The initial admin of the contract which is usually the creator of the contract. `init_admin` is also the initial value of admin. |
| `admin`    | Current `admin` of the contract initialized to `init_admin`. Certain critical actions can only be performed by the `admin`, e.g., changing the current implementation of the `ZilCrossChainManager` contract. |
|`initiator` | The user who calls the `ZilCrossChainManagerProxy` contract that in turn calls the `ZilCrossChainManager` contract. |


## Immutable Parameters

The table below lists the parameters that are defined at the contract deployment time and hence cannot be changed later on.

| Name | Type | Description |
|--|--|--|
|`init_crosschain_manager`| `ByStr20` | The address of the `ZilCrossChainManager` contract. |
|`init_admin`| `ByStr20` | The address of the admin. |

## Mutable Fields

The table below presents the mutable fields of the contract and their initial values.

| Name | Type | Initial Value |Description |
|--|--|--|--|
|`crosschain_manager`| `ByStr20` | `init_crosschain_manager` | Address of the current implementation of the `ZilCrossChainManager` contract. |
|`admin`| `ByStr20` | `init_owner` | Current `admin` of the contract. |

## Transitions

All the transitions in the contract can be categorized into two categories:
- **Housekeeping Transitions** meant to facilitate basic admin related tasks.
- **Relay Transitions** to redirect calls to the `ZilCrossChainManager` contract.

### Housekeeping Transitions

| Name | Params | Description |
|--|--|--|
|`UpgradeTo`| `new_crosschain_manager : ByStr20` |  Change the current implementation address of the `ZilCrossChainManager` contract. <br> :warning: **Note:** Only the `admin` can invoke this transition|
|`ChangeProxyAdmin`| `newAdmin : ByStr20` |  Change the current `stagingadmin` of the contract. <br> :warning: **Note:** Only the `admin` can invoke this transition.|
|`ClaimProxyAdmin` | `` |  Change the current `admin` of the contract. <br> :warning: **Note:** Only the `stagingadmin` can invoke this transition.|

### Relay Transitions

These transitions are meant to redirect calls to the corresponding `ZilCrossChainManager`
contract. While redirecting, the contract may prepare the `initiator` value that
is the address of the caller of the `ZilCrossChainManagerProxy` contract. The signature of
transitions in the two contracts is exactly the same expect the added last
parameter `initiator` for the `ZilCrossChainManager` contract.

| Transition signature in the `ZilCrossChainManagerProxy` contract  | Target transition in the `ZilCrossChainManager` contract |
|--|--|
|`Pause()` | `Pause(initiator : ByStr20)` |
|`UnPause()` | `UnPause(initiator : ByStr20)` |
|`UpdateAdmin(newAdmin: ByStr20)` | `UpdateAdmin(admin: ByStr20, initiator : ByStr20)`|
|`ClaimAdmin()` | `ClaimAdmin(initiator : ByStr20)`|
|`InitGenesisBlock(rawHeader: ByStr, pubkeys: List Pubkey)` | `InitGenesisBlock(rawHeader: ByStr, pubkeys: List Pubkey)`|
|`ChangeBookKeeper(rawHeader: ByStr, pubkeys: List Pubkey, sigList: List Signature)` | `ChangeBookKeeper(rawHeader: ByStr, pubkeys: List Pubkey, sigList: List Signature)`|
| `CrossChain(toChainId: Uint64, toContract: ByStr, method: ByStr, txData: ByStr)` | ` CrossChain(toChainId: Uint64, toContract: ByStr, method: ByStr, txData: ByStr)`|
| `VerifyHeaderAndExecuteTx(proof: Proof, rawHeader: ByStr, headerProof: Proof, curRawHeader: ByStr, headerSig: List Signature)` | `VerifyHeaderAndExecuteTx(proof: Proof, rawHeader: ByStr, headerProof: Proof, curRawHeader: ByStr, headerSig: List Signature)`|

# LockProxy Contract Specification

`LockProxy` is a contract that allows people to lock ZRC2 tokens and native zils to get corresponding tokens in target chain (e.g. ERC20 in ethereum) and vise versa.

## Roles and Privileges

The table below describes the roles and privileges that this contract defines:

| Role | Description & Privileges|                                    
| --------------- | ------------------------------------------------- |
| `init_admin`           | The initial admin of the contract which is usually the creator of the contract. `init_admin` is also the initial value of admin. |
| `admin`    | Current `admin` of the contract initialized to `init_admin`. Certain critical actions can only be performed by the `admin`. |
| `init_manager_proxy` | The initial cross chain manager proxy address. |
| `init_manager` | The initial cross chain manager address. |

## Immutable Parameters

The table below lists the parameters that are defined at the contract deployment time and hence cannot be changed later on.

| Name | Type | Description |
|--|--|--|
|`init_admin`| `ByStr20` | The address of the admin. |
|`init_manager_proxy`| `ByStr20` | The initial cross chain manager proxy address. |
|`init_manager`| `ByStr20` | The initial cross chain manager address. |

## Mutable Fields

The table below presents the mutable fields of the contract and their initial values.

| Name | Type | Initial Value |Description |
|--|--|--|--|
|`contractadmin`| `ByStr20` | `init_owner` | Current `admin` of the contract. |
|`manager`| `ByStr20` | `init_manager` | Address of the current `ZilCrossChainManager` contract. |
|`manager_proxy`| `ByStr20` | `init_manager_proxy` | Address of the current `ZilCrossChainManagerProxy` contract. |

## Transitions

| Name | Params | Description |
|--|--|--|
|`Lock`| `fromAssetHash: ByStr20, toChainId: Uint64, toAddress: ByStr, amount: Uint128` | Invoked by the user, a certin amount tokens will be locked in the proxy contract the invoker/msg.sender immediately, then the same amount of tokens will be unloked from target chain proxy contract at the target chain with chainId later.|
|`Unlock`| `txData: ByStr, fromContractAddr: ByStr, fromChainId: Uint64` | Invoked by the Zilliqa crosschain management contract, then mint a certin amount of tokens to the designated address since a certain amount was burnt from the source chain invoker.|


# More on cross chain infrastructure

- [polynetwork](https://github.com/polynetwork/poly)
- [zilliqa-relayer](https://github.com/Zilliqa/zilliqa-relayer)
