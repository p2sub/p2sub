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
	"encoding/hex"
	"errors"
	"net"
	"strings"

	"github.com/p2sub/p2sub/logger"
)

//Peer handle connections between peer
type Peer struct {
	listener          net.Listener
	bindAddress       string
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
	ChannelSize = 2
	BufferSize  = 1024
)

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ipAddr := ip.String(); ipAddr == "" || ipAddr == "<nil>" {
				continue
			} else {
				return ipAddr, nil
			}
		}
	}
	return "", errors.New("Could find external IP address")
}

func getPort(address string) string {
	if i := strings.LastIndex(address, ":"); i >= 0 {
		return address[i+1:]
	}
	return "<nil>"
}

//CreatePeer create a new peer
func CreatePeer(network, address string) *Peer {
	sugar := logger.GetSugarLogger()
	listener, err := net.Listen(network, address)
	ip, e := externalIP()
	port := getPort(address)
	if err != nil {
		sugar.Fatal("Not able to listen:", address, " protocol:", network, " error:", err)
	}
	if e != nil {
		sugar.Fatal("Not able get external IP:", e)
	}
	if port == "<nil>" || port == "" {
		sugar.Fatal("Invalid listent port")
	}
	p := &Peer{listener: listener,
		bindAddress:     ip + ":" + port,
		NewConnections:  make(chan net.Conn, ChannelSize),
		DeadConnections: make(chan net.Conn, ChannelSize),
		ActivePeers:     make(map[net.Conn]bool),
		DataChannel:     make(chan []byte, ChannelSize)}
	sugar.Infof("Listening network: %s IP: %s port: %s", network, ip, port)
	return p
}

//HandleLoop peer main loop
func (p *Peer) HandleLoop() {
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
			hexDump := hex.Dump(receivedData)
			sugar.Debugf("Received  %d bytes of data", len(receivedData))
			sugar.Debugf("Dumped data:\n%s", hexDump[:len(hexDump)-1])
			p.BroadCast(receivedData)
		}
	}
}

//Connect to another peer
func (p *Peer) Connect(network, address string) {
	sugar := logger.GetSugarLogger()
	if p.bindAddress == address {
		sugar.Infof("Skip loop connect -> %s", p.bindAddress)
		return
	}
	connect, err := net.Dial(network, address)
	if err == nil {
		p.NewConnections <- connect
	} else {
		sugar.Error("Not able to connect:", connect.RemoteAddr(), "X", connect.LocalAddr(), err)
	}
}

//BroadCast messages to network
func (p *Peer) BroadCast(data []byte) {
	sugar := logger.GetSugarLogger()
	sugar.Info("Active peers:", len(p.ActivePeers))
	for connect, i := range p.ActivePeers {
		if i {
			go func(connect net.Conn) {
				totalWritten := 0
				for totalWritten < len(data) {
					writtenThisCall, err := connect.Write(data[totalWritten:])
					if err != nil {
						p.DeadConnections <- connect
						break
					}
					totalWritten += writtenThisCall
				}
				sugar.Info("Sent data:", connect.LocalAddr(), connect.RemoteAddr())
			}(connect)
		} else {
			p.DeadConnections <- connect
		}
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
