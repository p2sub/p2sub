// Copyright 2019 Trần Anh Dũng <chiro@fkguru.com>
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

package transaction

import (
	"crypto/ed25519"
	"time"

	"github.com/p2sub/p2sub/address"
)

//New transaction
func New(from *address.Address, to *address.Address, flag Flag, method Method, data []byte) *Transaction {
	var toAddress []byte
	var dataBytes []byte
	if to == nil {
		toAddress = nil
	} else {
		toAddress = to.GetAddress()
	}
	if data == nil {
		dataBytes = []byte("")
	} else {
		dataBytes = data
	}
	return &Transaction{
		Signature: make([]byte, ed25519.SignatureSize),
		Flag:      flag,
		Method:    method,
		From:      from.GetAddress(),
		To:        toAddress,
		Time:      uint64(time.Now().Unix()),
		Length:    uint32(len(dataBytes)),
		Data:      dataBytes}
}

//NewBroardcast transaction
func NewBroardcast(from *address.Address, flag Flag, method Method, data []byte) *Transaction {
	return New(from, nil, flag|Broadcast, method, data)
}

//NewPrivate transaction
func NewPrivate(from *address.Address, to *address.Address, flag Flag, method Method, data []byte) *Transaction {
	return New(from, to, flag|Private, method, data)
}
