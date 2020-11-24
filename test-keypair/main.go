package main

import (
	"fmt"

	"github.com/p2sub/p2sub/config"
)

func main() {
	/*
		if _, err := os.Stat("./test.json"); err != nil {
			myKeyPair, _ := keypair.New()
			myKeyPair.SaveToFile("./test.json")
			fmt.Println(myKeyPair.GetPublicKey().Raw())
		} else {
			myKeyPair := keypair.LoadFromFile("./test.json")
			fmt.Println(myKeyPair.GetPublicKey().Raw())
		}*/
	cfg := config.GetConfig()
	cfg.Set("test", true)
	cfg.Set("test1", "Hello")
	fmt.Println(cfg.GetBool("test"))
	fmt.Println(cfg.GetBytes("test1"))
}
