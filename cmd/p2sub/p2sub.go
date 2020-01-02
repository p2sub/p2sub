package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/p2sub/p2sub/address"
	"github.com/p2sub/p2sub/transaction"
)

func main() {
	start := time.Now()
	//7zGCDka9k2cooRWPTBtPjLQMsLE5UdhoFUwzaMyw7DkQ
	sender := address.FromHexSeed("6578f93ce65b0c9d3bb578adc61d0092a62f340f9c342c9dd747731308ca32e5")
	message := "Hello world, I'm Chiro"
	tx1 := transaction.NewBroardcast(sender,
		transaction.MakeFlag(transaction.Broadcast, transaction.Sync, transaction.Ack),
		transaction.Ping,
		[]byte(message))
	tx1.Sign(sender)
	fmt.Println("Signed transaction:")
	tx1.Debug()
	fmt.Printf("Tx1: %x\n", tx1.Serialize())
	tx2 := transaction.Unserialize(tx1.Serialize())
	println("Is tx2 was signed properly:", tx2.Verify())
	fmt.Printf("Tx2: %x\n", tx2.Serialize())
	tx2.Debug()
	fmt.Println("Is the same?", bytes.Compare(tx1.Serialize(), tx2.Serialize()) == 0)
	fmt.Printf("Took: %s\n", time.Since(start))
}
