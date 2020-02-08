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
	"reflect"
	"time"

	"github.com/p2sub/p2sub/address"
	"github.com/p2sub/p2sub/logger"
	"github.com/p2sub/p2sub/packer"
)

//Transaction basic t in the system
type Transaction struct {
	Signature []byte
	Flag      Flag
	Method    Method
	From      []byte
	To        []byte
	Time      uint64
	Length    uint32
	Data      []byte
}

//Basic transaction size without data
const (
	TxBroadcastSize int = ed25519.SignatureSize + ed25519.PublicKeySize + 16
	TxPrivateSize       = ed25519.SignatureSize + 2*ed25519.PublicKeySize + 16
)

//Flag flag for t
type Flag uint

//Transaction identity type
const (
	Sync Flag = Flag((1 << iota) & 0xffff)
	Ack
	AckAck
	Not
	Rst
	//Two of these flags won't be co-exist
	Private
	Broadcast
)

//Method using in
type Method uint

//Method type
const (
	Invalid Method = Method(iota)
	Ping
	Pong
	Pub
	RequestKey
	ExchangeKey
	Gosship
	GosshipAboutGosship
	Authentication
)

//MakeFlag create flag value by turning on each bit
func MakeFlag(flags ...Flag) Flag {
	t := uint16(0)
	for i := 0; i < len(flags); i++ {
		t |= uint16(flags[i])
	}
	return Flag(t)
}

//GetFlag get transaction flag
func (t *Transaction) GetFlag() Flag {
	return t.Flag
}

//GetMethod get packate method
func (t *Transaction) GetMethod() Method {
	return t.Method
}

//IsFlag check for flag positive
func (t *Transaction) IsFlag(flag Flag) bool {
	return flag&t.Flag == flag
}

//IsMethod check for flag positive
func (t *Transaction) IsMethod(methodType Method) bool {
	return t.GetMethod() == methodType
}

//Debug string
func (t *Transaction) Debug() {
	sugar := logger.GetSugarLogger()
	logger.HexDump("Transaction's signature:", t.Signature)
	sugar.Debugf("Flag: %x (BroadCast=%t, Private=%t, Sync=%t)",
		t.Flag,
		t.IsFlag(Broadcast),
		t.IsFlag(Private),
		t.IsFlag(Sync))
	sugar.Debugf("Method: %d", t.Method)
	sugar.Debugf("From: %x", t.From)
	if t.IsFlag(Private) {
		sugar.Debugf("To: %x", t.To)
	}
	sugar.Debugf("Time: %s", time.Unix(int64(t.Time), 0))
	sugar.Debugf("Length: %d", t.Length)
	logger.HexDump("Data", t.Data)
}

//Serialize a t to serialized data
func (t *Transaction) Serialize() []byte {
	s := packer.NewSerialize()
	s.Write(t.Signature,
		uint16(t.Flag), uint16(t.Method),
		t.From,
		t.To,
		t.Time,
		t.Length,
		t.Data)
	return s.Bytes()
}

//Sign sign t by sender address
func (t *Transaction) Sign(sender *address.Address) []byte {
	if sender.IsSignKey() {
		signedMessage := sender.Sign(t.Serialize()[ed25519.SignatureSize:])
		if signedMessage != nil {
			t.Signature = signedMessage[:ed25519.SignatureSize]
			return signedMessage
		}
	}
	return nil
}

//Verify signed t
func (t *Transaction) Verify() bool {
	if t != nil && t.From != nil {
		senderAddress := address.FromPublicKey(t.From)
		return senderAddress.Verify(t.Serialize())
	}
	return false
}

//Unserialize transform serialized data to t
func Unserialize(rawTx []byte) *Transaction {
	//RawTx was too small
	if len(rawTx) >= TxBroadcastSize {
		u := packer.NewUnserialize(rawTx)
		var to []byte
		var data []byte
		signature := u.Pop(reflect.Slice, ed25519.SignatureSize).([]byte)
		flag := u.Pop(reflect.Uint16).(uint16)
		method := u.Pop(reflect.Uint16).(uint16)
		from := u.Pop(reflect.Slice, ed25519.PublicKeySize).([]byte)
		if flag&uint16(Private) > 0 {
			//Transaction is private but have size smaller
			//than basic private transaction
			if u.Size() < TxPrivateSize {
				return nil
			}
			to = u.Pop(reflect.Slice, ed25519.PublicKeySize).([]byte)
		} else {
			to = nil
		}
		time := u.Pop(reflect.Uint64).(uint64)
		length := u.Pop(reflect.Uint32).(uint32)
		//Remaining bytes should larger than required
		if length <= uint32(u.Len()) {
			if b, ok := u.Pop(reflect.Slice, int(length)).([]byte); ok {
				data = b
			} else {
				data = nil
			}
			//Return t
			return &Transaction{
				Signature: signature,
				Flag:      Flag(flag),
				Method:    Method(method),
				From:      from,
				To:        to,
				Time:      time,
				Length:    length,
				Data:      data}
		}
	}
	return nil
}
