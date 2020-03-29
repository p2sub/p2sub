// Copyright 2019 P2Sub Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package address

import (
	"crypto/ed25519"

	"github.com/btcsuite/btcutil/base58"
)

//Address node indentity
type Address struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

//Address constant
const (
	SignatureSize   = 64
	NullMessageSize = 0
)

//GetPrivateKey Get private key to export configuration
func (a *Address) GetPrivateKey() []byte {
	if a.privateKey != nil {
		return a.privateKey
	}
	return nil
}

//GetSeed Get seed  to export configuration
func (a *Address) GetSeed() []byte {
	if a.privateKey != nil {
		return a.privateKey.Seed()
	}
	return nil
}

//GetAddress get public key of key pair
func (a *Address) GetAddress() []byte {
	if a.publicKey != nil && len(a.publicKey) == ed25519.PublicKeySize {
		return []byte(a.publicKey)
	} else if a.privateKey != nil && len(a.privateKey) == ed25519.PrivateKeySize {
		if publicKey, ok := a.privateKey.Public().(ed25519.PublicKey); ok {
			a.publicKey = publicKey
		}
		return []byte(a.publicKey)
	}
	return nil
}

//GetBase58Address get base58 encoded of public key
func (a *Address) GetBase58Address() string {
	if vk := a.GetAddress(); vk != nil {
		return base58.Encode(a.publicKey)
	}
	return "<nil>"
}

//IsSignKey is this address contain private key and able to sign transaction
func (a *Address) IsSignKey() bool {
	return a.privateKey != nil
}

//Sign sign a message
func (a *Address) Sign(message []byte) []byte {
	if a.IsSignKey() && len(message) > NullMessageSize {
		signature := ed25519.Sign(a.privateKey, message)
		signedMessage := make([]byte, len(message)+ed25519.SignatureSize)
		copy(signedMessage[:ed25519.SignatureSize], signature[:])
		copy(signedMessage[ed25519.SignatureSize:], message[:])
		return signedMessage
	}
	return nil
}

//Verify signed a message
func (a *Address) Verify(signedMessage []byte) bool {
	if vk := a.GetAddress(); vk != nil && len(signedMessage) > SignatureSize {
		signature := make([]byte, ed25519.SignatureSize)
		message := make([]byte, len(signedMessage)-ed25519.SignatureSize)
		copy(signature[:], signedMessage[:ed25519.SignatureSize])
		copy(message[:], signedMessage[ed25519.SignatureSize:])
		return ed25519.Verify(vk, message, signature)
	}
	return false
}
