## Introduction
Dione Protocol is an open-source L1 Blockchain which is a fork of Avalanche Blockchain Network.

## Dione Platform
Dione Protocol is a heterogeneous network of blockchains. As opposed to homogeneous networks, where all applications reside in the same chain, heterogeneous networks allow separate chains to be created for different applications.

The Primary Network is a special Subnet that contains all validators (including validators of any custom Subnets). A node can become a validator for the Primary Network by staking Dione tokens.

The Dione Protocol uses Avalanche Snowman Consensus Protocol as Primary network and Subnets consensus.

## Virtual Machines
Dione Protocol allows developers to implement their Dapp in the language of their choice using the Virtual Machine (VM) framework. VMs define blockchains, and serve as application-level blueprints for how blockchains are created.

Developers can implement VMs in any language, and use libraries and tech stacks that they're familiar with. Developers have fine control over the behavior of their blockchain, and can redefine the rules of a blockchain to fit any use-case they have.

## Validators
All Dione Protocol validators as members of the Dione Protocol primary network are required to run three VMs:

Coreth: Defines the Contract Chain (C-Chain); supports smart contract functionality and is EVM-compatible.
Platform VM: Defines the Platform Chain (P-Chain); supports operations on staking and Subnets.
Dione VM: Defines the Exchange Chain (X-Chain); supports operations on Dione Native Tokens.
Validators are able to install additional VMs on their node to validate additional Subnets in the Dione ecosystem. In exchange, validators receive staking rewards in the form of a reward token configured by Subnets.

## Consensus
Dione Protocol utilize Avalanche Snowman Consensus Protocol as Primary network and Subnets consensus.

### Consensus Algorithm
```
    preference := pizza
    consecutiveSuccesses := 0
    while not decided:
    ask k random people their preference
    if >= α give the same response:
        preference := response with >= α
        if preference == old preference:
        consecutiveSuccesses++
        else:
        consecutiveSuccesses = 1
    else:
        consecutiveSuccesses = 0
    if consecutiveSuccesses > β:
        decide(preference)
```

## Changelog
- Forked Avalanchego, coreth, subnets repos
- Updated Imports
- Updated Chains config
- Updated Readme files, Documents

## Refrences
[Avalanchego](https://github.com/ava-labs/avalanchego)
[Coreth](https://github.com/ava-labs/coreth)
[Subnets](https://github.com/ava-labs/subnet-evm)
[Avalanche Docs](https://docs.avax.network/)
[Snowman Consensus Algorithm Explained](https://docs.avax.network/overview/getting-started/avalanche-consensus#algorithm-explained)
