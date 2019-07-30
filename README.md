# Stellar Anchor

Stellar Anchor is a web server that implements the Stellar Ecosystem Proposals. The purpose for this project to provide a foundation for creating Stellar Applications. This web server can represent a fake financial institution for your Stellar wallet to connect to.

![img](https://travis-ci.com/drdgvhbh/StellarAnchor.svg?branch=master)

## Features

Protocols

- [x] [SEP-0001](https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0001.md) — stellar.toml
- [ ] [SEP-0006](https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0006.md) — Anchor/Client interoperability
- [ ] [SEP-0009](https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0009.md) — Standard KYC / AML fields
- [x] [SEP-0010](https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0010.md) — Stellar Web Authentication
- [x] [SEP-0012](https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0012.md) — Anchor/Client customer info transfer

Assets

- [ ] Ethereum (ETH)
- [x] Lumens (XLM)
- [ ] Bitcoin (BTC)

## Building from Source

```sh
git clone https://github.com/drdgvhbh/StellarAnchor.git
cd StellarAnchor
go mod vendor
make start
```

## Running the Tests

```sh
make test-unit
make test-e2e
make cover
make coveralls
```
