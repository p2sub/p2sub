package transaction_test

import (
	"testing"

	"github.com/p2sub/p2sub/address"
	"github.com/p2sub/p2sub/transaction"
)

func bytesCompare(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var c byte = 0
	for i := 0; c == 0 && i < len(a); i += 1 {
		c = a[i] ^ b[i]
	}
	return c == 0
}

func TestTransaction(t *testing.T) {
	from := address.New()
	to := address.New()
	tx := transaction.New(from, to, transaction.Private, transaction.Ping, []byte{0x55, 0x44})
	tx.Sign(from)
	serializedTx := tx.Serialize()
	recoverTx := transaction.Unserialize(serializedTx)
	result := true
	result = result && bytesCompare(recoverTx.From, from.GetAddress())
	result = result && bytesCompare(recoverTx.To, to.GetAddress())
	result = result && tx.Verify()
	if !result {
		t.Error("Recover tx was difference to original tx")
	}

}
