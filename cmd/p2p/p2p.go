package main

import (
	"net"
	"strings"

	"github.com/p2sub/p2sub/logger"
	"github.com/p2sub/p2sub/network"
	"github.com/p2sub/p2sub/p2p"
)

var sugar = logger.GetSugarLogger()

//BroadCast messages to network
func handler(p *p2p.Peer, data []byte) {
	sugar.Info("Active peers:", len(p.ActivePeers))
	for connect, i := range p.ActivePeers {
		if i {
			go func(connect net.Conn) {
				totalWritten := 0
				for totalWritten < len(data) {
					writtenBytes, err := connect.Write(data[totalWritten:])
					if err != nil {
						p.DeadConnections <- connect
						break
					}
					totalWritten += writtenBytes
				}
				sugar.Info("Sent data:", connect.LocalAddr(), "->", connect.RemoteAddr())
			}(connect)
		} else {
			p.DeadConnections <- connect
		}
	}
}

func main() {
	p := p2p.CreatePeer("tcp", "[::1]:6262")
	go p.Listen()
	addresses, _ := network.ListAddress()
	for _, ip := range addresses {
		if strings.Index(ip, ":") >= 0 {
			p.Connect("tcp", "["+ip+"]:6262")
		} else {
			p.Connect("tcp", ip+":6262")
		}
	}
	defer p.HandleLoop(handler)
}
