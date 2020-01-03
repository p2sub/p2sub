package main

import (
	"strings"

	"github.com/p2sub/p2sub/network"
	"github.com/p2sub/p2sub/p2p"
)

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
	defer p.HandleLoop()
}
