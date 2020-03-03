package main

import (
	"fmt"
	"log"

	"github.com/p2sub/p2sub/leveldb"
)

func main() {
	var insertStatus int
	var listItems []leveldb.Item
	insertItem := leveldb.Item {
		Key: []byte("key-2"),
		Value: []byte("Value 1"),
	}
	lvlManager, err := leveldb.NewManager("./db", nil)
	defer lvlManager.Close()
	if err != nil {
		log.Fatal(err)
	}
	insertStatus, err = lvlManager.AddNewItem(insertItem)
	if err != nil {
		log.Fatal(err)
	}
	if insertStatus == 1 {
		fmt.Println("Successfully insert")
	}
	// Get all items from database
	listItems, err = lvlManager.GetAllItems()
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range listItems {
		fmt.Printf("Key: %s | Value: %s\n", string(item.Key), string(item.Value))
	}
}