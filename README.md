# Introduction

Event-driven Architecture is a good idea to build an distributed system, instead of messing whole system with many duplicated requests. We could create a unique channel to deliver events and its data. When we try to explore pub/sub channel pattern, we are also find out, it has it own problems.

* Could not verify publisher
* Could not create a verifiable private channel
* There are a centralized service or queue, it's become single point of failure
* Recover from error state

In the try to improve pub/sub pattern, we created P2SUB is a brand new distributed pub/sub channel with many advance features waiting you to explore

## Features

* **Zero Configuration**: _The configuration will be done magically (self-signed certificate)_
* **Cryptography Guarantee**: _Based on elliptic curve digital signature (ED25519)_
* **Byzantine Fault Tolerant**: _Fast and robust distributed system with fault tolerance (BFT)_
* **End-to-End Encryption**: _All data transfer will be encrypted (AES128)_
* **Verifiable Publisher**: _All publisher will sign their messages_
* **Transaction Based**: _All message will be treated like transaction_
* **No Single Point of Failure**: _Instead of centralized queue or service, we use blockchain as its backbone_

## License

P2SUB is licensed under [Apache License 2.0](https://github.com/chiro-hiro/p2sub/blob/master/LICENSE)
