package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

//Client handle client side
type Client struct {
	newConnects  chan net.Conn
	DeadConnects chan net.Conn
	MapConnects  map[net.Conn]bool
	MapReceiver  map[string]net.Conn
	buffers      chan []byte
}

const (
	channelSize = 2
	bufferSize  = 1024
)

//CreateClient create new client
func CreateClient(protocol string, bindSock string) *Client {
	return &Client{
		newConnects:  make(chan net.Conn, channelSize),
		DeadConnects: make(chan net.Conn, channelSize),
		buffers:      make(chan []byte, channelSize),
		MapConnects:  make(map[net.Conn]bool)}
}

//Dial start dailing
func (s *Client) Dial(network, address string) {
	go func() {
		connect, err := net.Dial(network, address)
		if err != nil {
			panic(err)
		}
		s.newConnects <- connect
	}()
}

//HandleLoop Client main loop
func (s *Client) HandleLoop(handler func([]byte)) {
	for {
		select {
		case connect := <-s.newConnects:
			s.MapConnects[connect] = true
			log.Println("[NEW CONNECT]", connect.RemoteAddr(), "->", connect.LocalAddr())
			go func() {
				buf := make([]byte, bufferSize)
				for {
					nbyte, err := connect.Read(buf)
					if err != nil {
						s.DeadConnects <- connect
						break
					} else {
						fragment := make([]byte, nbyte)
						copy(fragment, buf[:nbyte])
						log.Println("[RECEIVE DATA]", connect.RemoteAddr(), "->", connect.LocalAddr())
						s.buffers <- fragment
					}
				}
			}()
		case deadConnect := <-s.DeadConnects:
			log.Println("[CLOSE CONNECT]", deadConnect.RemoteAddr(), "->", deadConnect.LocalAddr())
			err := deadConnect.Close()
			if err != nil {
				log.Println("[ERROR]", err)
			}
			delete(s.MapConnects, deadConnect)
		case buffer := <-s.buffers:
			go handler(buffer)
			/*
				log.Printf("[RECEIVE DATA] Data:\n %s", hex.Dump(buffer))
				for connect, i := range s.mapConnects {
					if i {
						go func(connect net.Conn) {
							totalWritten := 0
							for totalWritten < len(buffer) {
								writtenThisCall, err := connect.Write(buffer[totalWritten:])
								if err != nil {
									s.DeadConnects <- connect
									break
								}
								totalWritten += writtenThisCall
							}
						}(connect)
					} else {
						s.DeadConnects <- connect
					}
				}
			*/
		}
	}
}

func main() {
	inData := make(chan string, 1024)
	receiveData := make(chan string, 1024)
	// connect to this socket
	conn, _ := net.Dial("tcp", "[::1]:6262")

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
				conn, _ = net.Dial("tcp", "[::1]:6262")
			}
		}
	}
}
