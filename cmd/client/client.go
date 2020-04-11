package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	inData := make(chan string, 1024)
	receiveData := make(chan string, 1024)
	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:6001")

	for {
		go func() {
			// read in input from stdin
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("> ")
			text, _ := reader.ReadString('\n')
			// send to socket
			inData <- text
		}()
		go func() {
			// listen for reply
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err == nil {
				receiveData <- message
			}
		}()
		select {
		case receive := <-receiveData:
			fmt.Print("\n< " + receive)
		case send := <-inData:
			_, e := fmt.Fprintf(conn, send+"\n")
			if e != nil {
				conn.Close()
				conn, _ = net.Dial("tcp", "127.0.0.1:6001")
			}
		}
	}
}
