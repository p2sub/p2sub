package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func main() {
	exit := false
	inData := make(chan string, 1024)
	for {
		go func() {
			// read in input from stdin
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("> ")
			text, _ := reader.ReadString('\n')
			// send to socket
			inData <- text
		}()
		select {
		case stdinData := <-inData:
			stdinData = stdinData[:len(stdinData)-1]
			log.Print(hex.Dump([]byte(stdinData)))
			if stdinData == "exit" {
				exit = true
			}
			break
		}
		if exit {
			break
		}
	}
}
