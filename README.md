# ProximaX Sirius Blockchain Golang SDK

Official ProximaX Sirius Blockchain SDK Library in Golang.

The ProximaX Sirius Catapult Chain Golang SDK works as a lightweight Golang library for interacting with the Sirius Blockchain. It provides a complete library set coverage, and supports synchronous and asynchronous requests. 

<p align="center">
    <img src="./doc/catapult-nem2-sdk-go.jpg">
</p>

## Getting Started

All catapult stuff starts from importing sdk:

```go
import "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
```

Then create a Catapult network configuration

For *Testnet*:
```go
conf, err := sdk.NewConfig("http://localhost:3000",sdk.Testnet)
```
Or for *Mainnet*:
```go
conf, err := sdk.NewConfig("http://localhost:3000",sdk.Mainnet)
```

Construct a new REST Catapult client
```go
client := sdk.NewClient(nil, conf)
```

Use this client to get current blockchain height

```go
chainHeight, err := client.Blockchain.GetChainHeight(context.Background())
```

A [Context](https://golang.org/pkg/context/) type is the first argument in any service method for specifying
deadlines, cancelation signals, and other request-scoped values

## Wiki / Examples

For more examples, go to our [wiki](https://github.com/proximax-storage/go-xpx-chain-sdk/wiki)

## API docs

If you want to refer to go docs in markdown, check out api [docs](./api.md).

## Core Contributors

 + [@Wondertan](https://github.com/Wondertan)
 + [@ruslanBik4](https://github.com/ruslanBik4)
 + [@slackve](https://github.com/slackve)
 + [@brambear](https://github.com/alvin-reyes)
 + [@carlocayos](https://github.com/carlocayos)


## Acknowledgement

This library is a fork of nem2-go-sdk which was also contributed by ProximaX Limited to NEM Foundation. 

## [Contribution](CONTRIBUTING.md)

We'd love to get more people involved in the project. Please feel free to raise any issues or PR and we'll review your contribution.
    
Copyright (c) 2019 ProximaX Limited
