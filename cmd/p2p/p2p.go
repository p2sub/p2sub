package main

import (
	"github.com/p2sub/p2sub/p2p"
)

func main() {
	p := p2p.CreatePeer("tcp", ":6262")
	go p.Listen()
	p.Connect("tcp", "10.215.1.37:6262")
	defer p.HandleLoop()
}
