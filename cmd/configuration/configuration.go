package main

import (
	"crypto/ed25519"
	"math/rand"

	"github.com/p2sub/p2sub/address"
	"github.com/p2sub/p2sub/configuration"
	"github.com/p2sub/p2sub/packer"
)

// Get random nonce number
func randomNonce() uint32 {
	return uint32(rand.Intn(0xffffffff))
}

// Use private key of notary to sign public identity of node
func signNodeID(notaryConf configuration.Config, nodeConf configuration.Config) []byte {
	nodeID := packer.NewSerialize()
	// Construct node identity
	nodeID.Write(nodeConf.NodeAddress.GetAddress(),
		nodeConf.Nonce,
		nodeConf.NodeType,
		notaryConf.NodeAddress.GetAddress(),
		nodeConf.Name)
	//Only get signature
	return notaryConf.NodeAddress.Sign(nodeID.Bytes())[:ed25519.SignatureSize]
}

// Verify node identity by using notary public key
func verifyNodeID(nodeConf configuration.Config) bool {
	nodeID := packer.NewSerialize()
	// Serialize node identity for signing
	nodeID.Write(nodeConf.NodeAddress.GetAddress(),
		nodeConf.Nonce,
		nodeConf.NodeType,
		nodeConf.Issuer,
		nodeConf.Name)
	nodeByteID := nodeID.Bytes()
	// Re construct signed message
	signedMessage := make([]byte, len(nodeByteID)+len(nodeConf.Signature))
	copy(signedMessage[ed25519.SignatureSize:], nodeByteID)
	copy(signedMessage[:ed25519.SignatureSize], nodeConf.Signature)
	notaryVk := address.FromPublicKey(nodeConf.Issuer)
	return notaryVk.Verify(signedMessage)
}

func main() {

	notary := configuration.Config{
		Name:          "notary1",
		BindPort:      "6001",
		BindHost:      "127.0.0.1",
		Nonce:         randomNonce(),
		NodeType:      configuration.NodeNotary,
		Signature:     nil,
		NodeAddress:   *address.New(),
		ConfigService: "127.0.0.1:6001",
	}

	master1 := configuration.Config{
		Name:          "master1",
		BindPort:      "6011",
		BindHost:      "127.0.0.1",
		Nonce:         randomNonce(),
		NodeType:      configuration.NodeMaster,
		Signature:     nil,
		Issuer:        notary.NodeAddress.GetAddress(),
		NodeAddress:   *address.New(),
		ConfigService: "127.0.0.1:6001",
	}
	// Assign signature value
	master1.Signature = signNodeID(notary, master1)
	if !verifyNodeID(master1) {
		panic("Master 1 was invalid")
	}

	master2 := configuration.Config{
		Name:          "master2",
		BindPort:      "6012",
		BindHost:      "127.0.0.1",
		Nonce:         randomNonce(),
		NodeType:      configuration.NodeMaster,
		Issuer:        notary.NodeAddress.GetAddress(),
		Signature:     nil,
		NodeAddress:   *address.New(),
		ConfigService: "127.0.0.1:6001",
	}
	// Assign signature value
	master2.Signature = signNodeID(notary, master2)
	if !verifyNodeID(master2) {
		panic("Master 2 was invalid")
	}

	master3 := configuration.Config{
		Name:          "master3",
		BindPort:      "6013",
		BindHost:      "127.0.0.1",
		Nonce:         randomNonce(),
		NodeType:      configuration.NodeMaster,
		Issuer:        notary.NodeAddress.GetAddress(),
		Signature:     nil,
		NodeAddress:   *address.New(),
		ConfigService: "127.0.0.1:6001",
	}
	// Assign signature value
	master3.Signature = signNodeID(notary, master3)
	if !verifyNodeID(master3) {
		panic("Master 3 was invalid")
	}

	//Export configuration
	notary.Export("./conf.d/notary.json")
	master1.Export("./conf.d/master1.json")
	master2.Export("./conf.d/master2.json")
	master3.Export("./conf.d/master3.json")

}
