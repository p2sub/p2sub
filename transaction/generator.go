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
