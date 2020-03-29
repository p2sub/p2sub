# Introduction

Event-driven Architecture might be a good choice to build a distributed system, instead of messing whole system with many duplicated requests. We could create an unified channel to deliver events and its data. When we try to explore pub/sub pattern, we are also find out it has its own problems.

* Could not verify publishers
* Could not create a verifiable private channel
* Could not verify message ownership
* There are a centralized service or queue, it's become single point of failure
* Recover from error states
* Acknowledge of states

In the try to improve pub/sub pattern, we created P2SUB. It is a brand new distributed pub/sub pattern with many advance features waiting for you to explore.

## Features

* **Zero Configuration**: _The configuration will be done magically (self-signed certificate)_
* **Cryptography Guarantee**: _Based on elliptic curve digital signature (ED25519)_
* **Byzantine Fault Tolerant**: _Fast and robust distributed system with fault tolerance (BFT)_
* **End-to-End Encryption**: _All data transfer will be encrypted (AES128 + Diffie–Hellman key exchange)_
* **Verifiable Publisher**: _All publisher will sign their messages_
* **Transaction Based**: _All message will be treated like transaction_
* **No Single Point of Failure**: _Instead of centralized queue or service, we use blockchain as its backbone_

## Development

You could generate your own node identities by run this command

```bash
$ go run ./cmd/configuration/configuration.go
```

Configuration will be appear in `./conf.d`

**note:** After alpha stage we won't use configuration file anymore.

## License

P2SUB is licensed under [Apache License 2.0](https://github.com/chiro-hiro/p2sub/blob/master/LICENSE)
