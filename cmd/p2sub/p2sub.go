package main

import (
	"bytes"
	"time"

	"github.com/p2sub/p2sub/address"
	"github.com/p2sub/p2sub/logger"
	"github.com/p2sub/p2sub/transaction"
)

func main() {
	sugar := logger.GetSugarLogger()
	start := time.Now()
	//7zGCDka9k2cooRWPTBtPjLQMsLE5UdhoFUwzaMyw7DkQ
	sender := address.FromHexSeed("6578f93ce65b0c9d3bb578adc61d0092a62f340f9c342c9dd747731308ca32e5")
	message := "Hello world, I'm Chiro"
	tx1 := transaction.NewBroardcast(sender,
		transaction.MakeFlag(transaction.Sync, transaction.Ack),
		transaction.Ping,
		[]byte(message))
	logger.HexDump("Unsigned transaction 1:", tx1.Serialize())
	tx1.Sign(sender)
	sugar.Info("Signed transaction:")
	tx1.Debug()
	logger.HexDump("Signed transaction 1:", tx1.Serialize())
	tx2 := transaction.Unserialize(tx1.Serialize())
	sugar.Info("Is tx2 was signed properly: ", tx2.Verify())
	logger.HexDump("Received transaction 2:", tx2.Serialize())
	tx2.Debug()
	sugar.Info("Is the same? ", bytes.Compare(tx1.Serialize(), tx2.Serialize()) == 0)
	sugar.Infof("Took: %s", time.Since(start))
	/*
		var confs configuration.Configs
		confs = append(confs, configuration.ConfigItem{
			Name:      "chiro-node-0",
			PublicKey: "some-address-0",
			Signature: "test-signature-0",
		})
		confs = append(confs, configuration.ConfigItem{
			Name:      "chiro-node-1",
			PublicKey: "some-address-1",
			Signature: "test-signature-1",
		})
		if confs.Export("./test.json") {
			if conf := configuration.Import("./test.json"); conf != nil {
				sugar.Debug("Configuration", conf.ToString())
			}
		}*/
}
