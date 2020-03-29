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
	"encoding/hex"

	"github.com/btcsuite/btcutil/base58"
)

//New Generate new address from random
func New() *Address {
	//If nil was set then crypto/rand will be used
	if vk, pk, err := ed25519.GenerateKey(nil); err == nil {
		return &Address{publicKey: vk,
			privateKey: pk}
	}
	return nil
}

//FromBase58Address Import address from base58 address string for verification only
func FromBase58Address(address string) *Address {
	if vk := base58.Decode(address); len(vk) == ed25519.PublicKeySize {
		return &Address{publicKey: vk}
	}
	return nil
}

//FromSeed generate address from seed
func FromSeed(seed []byte) *Address {
	pk := ed25519.NewKeyFromSeed(seed)
	if vk, ok := pk.Public().(ed25519.PublicKey); ok {
		return &Address{privateKey: pk,
			publicKey: vk}
	}
	return nil
}

//FromPublicKey create address from public key
func FromPublicKey(vk []byte) *Address {
	if vk != nil && len(vk) == ed25519.PublicKeySize {
		return &Address{publicKey: vk}
	}
	return nil
}

//FromHexSeed gernate from hex string of seed
func FromHexSeed(seed string) *Address {
	if len(seed)%2 == 1 {
		seed = "0" + seed
	}
	if s, e := hex.DecodeString(seed); e == nil {
		return FromSeed(s)
	}
	return nil
}

//FromPrivateKey create new address from private key
func FromPrivateKey(pk []byte) *Address {
	a := Address{
		privateKey: pk,
	}
	a.publicKey = a.GetAddress()
	return &a
}
