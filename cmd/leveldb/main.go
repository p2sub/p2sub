package main

import (
	"fmt"
	"log"

	"github.com/p2sub/p2sub/leveldb"
)

func main() {
	manager, err := leveldb.New("./db", nil)
	addItem := leveldb.Item {
		Key: []byte("Item-1"),
		Value: []byte("Value-1"),
	}
	if err != nil {
		log.Println(err.Error())
	}
	defer manager.Close()
	manager.AddNewItem(addItem)
	item, err := manager.FindItemByKey([]byte("Item-2"))
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Printf("Key: %s | Value: %s\n", string(item.Key), string(item.Value))
}