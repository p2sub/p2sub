// Copyright 2019 Trần Anh Dũng <chiro@fkguru.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package p2p

import (
	"net"
	"strings"

	"github.com/p2sub/p2sub/logger"
	"github.com/p2sub/p2sub/network"
)

//Peer handle connections between peer
type Peer struct {
	listener          net.Listener
	bindAddress       string
	addresses         []string
	NewConnections    chan net.Conn
	DeadConnections   chan net.Conn
	ActivePeers       map[net.Conn]bool
	ActiveConnections map[string]net.Conn
	DataChannel       chan []byte
}

//SmartCondition type
type SmartCondition map[string]func()

//Basic configuration
const (
	ChannelSize uint = 2
	BufferSize       = 1024
)

func getIndexInSlice(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

func getIPAddress(address string) string {
	if i := strings.LastIndex(address, ":"); i >= 0 {
		return strings.Trim(address[:i], "[]")
	}
	return "<nil>"
}

func getPort(address string) string {
	if i := strings.LastIndex(address, ":"); i >= 0 {
		return address[i+1:]
	}
	return "<nil>"
}

//CreatePeer create a new peer
func CreatePeer(proto, address string) *Peer {
	sugar := logger.GetSugarLogger()
	listener, err := net.Listen(proto, address)
	port := getPort(address)
	if err != nil {
		sugar.Fatal("Not able to listen:", address, " protocol:", proto, " error:", err)
	}
	if port == "<nil>" || port == "" {
		sugar.Fatal("Invalid listent port")
	}
	p := &Peer{listener: listener,
		bindAddress:     address,
		NewConnections:  make(chan net.Conn, ChannelSize),
		DeadConnections: make(chan net.Conn, ChannelSize),
		ActivePeers:     make(map[net.Conn]bool),
		DataChannel:     make(chan []byte, ChannelSize)}
	if addresses, err := network.ListAddress(); err == nil {
		p.addresses = addresses
	} else {
		sugar.Fatal("Not able to list all IP addresses:", err)
	}
	sugar.Infof("Listening protocol %s bind to (%s)", proto, address)
	return p
}

//HandleLoop peer main loop
func (p *Peer) HandleLoop(handler func(p *Peer, data []byte)) {
	sugar := logger.GetSugarLogger()
	for {
		select {
		case connect := <-p.NewConnections:
			p.ActivePeers[connect] = true
			sugar.Info("New connection:", connect.RemoteAddr(), "->", connect.LocalAddr())
			go func() {
				buf := make([]byte, BufferSize)
				for {
					nbyte, err := connect.Read(buf)
					if err != nil {
						p.DeadConnections <- connect
						break
					} else {
						chunk := make([]byte, nbyte)
						copy(chunk, buf[:nbyte])
						sugar.Info("Received data:", connect.RemoteAddr(), "->", connect.LocalAddr())
						p.DataChannel <- chunk
					}
				}
			}()
		case deadConnect := <-p.DeadConnections:
			sugar.Info("Close connect", deadConnect.RemoteAddr(), "->", deadConnect.LocalAddr())
			err := deadConnect.Close()
			if err != nil {
				sugar.Error("Could not close connect:", err)
			}
			delete(p.ActivePeers, deadConnect)
		case receivedData := <-p.DataChannel:
			sugar.Debugf("Received  %d bytes of data", len(receivedData))
			logger.HexDump("Dumped data:", receivedData)
			// Trigger handler
			handler(p, receivedData)
		}
	}
}

//Connect to another peer
func (p *Peer) Connect(network, address string) {
	sugar := logger.GetSugarLogger()
	bindIPAddress := getIPAddress(p.bindAddress)
	bindPort := getPort(p.bindAddress)
	targetIPAddress := getIPAddress(address)
	targetPort := getPort(address)
	if targetPort == bindPort {
		if (bindIPAddress == targetIPAddress) ||
			(getIndexInSlice(p.addresses, targetIPAddress) >= 0) {
			sugar.Infof("Skip loop connect -> %s", address)
			return
		}
	}
	connect, err := net.Dial(network, address)
	if err == nil {
		p.NewConnections <- connect
	} else {
		sugar.Error("Not able to connect:", connect.RemoteAddr(), "X", connect.LocalAddr(), err)
	}
}

//Listen current peer
func (p *Peer) Listen() {
	sugar := logger.GetSugarLogger()
	go func() {
		for {
			connect, err := p.listener.Accept()
			if err == nil {
				p.NewConnections <- connect
			} else {
				sugar.Error("Not able accept connect:", connect.RemoteAddr(), "X", connect.LocalAddr(), err)
			}
		}
	}()
}
