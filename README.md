# Ethereum Client

## Install

```sh
go get -u github.com/nextyio/eth-client
go get -u github.com/kardianos/govendor

govendor fetch -tree github.com/ethereum/go-ethereum/crypto/secp256k1

govendor install +local
```
