package main

import (
	"flag"
	"net"
	"os"

	"github.com/p2sub/p2sub/configuration"
	"github.com/p2sub/p2sub/logger"
	"github.com/p2sub/p2sub/p2p"
)

var sugar = logger.GetSugarLogger()

//Load configuration from file
func loadConfig() (conf *configuration.Config, err error) {
	configFile := flag.String("config", "", "Path to configuration file")
	flag.Parse()
	if *configFile == "" {
		flag.Usage()
		os.Exit(0)
	}
	sugar.Info("Load configuration from: ", *configFile)
	return configuration.Import(*configFile)
}

//Handle connect from peer-to-peer network
func master(p *p2p.Peer, data []byte) {
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

//Handle connect from peer-to-peer network
func notary(p *p2p.Peer, data []byte) {
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

//Bot current Node
func bootNode(conf *configuration.Config, handler func(p *p2p.Peer, data []byte)) {
	p2pNode := p2p.CreatePeer("tcp", conf.BindHost+":"+conf.BindPort)
	go p2pNode.Listen()
	defer p2pNode.HandleLoop(handler)
}

func main() {
	config, err := loadConfig()
	if err == nil {
		if config.NodeType == configuration.NodeNotary {
			bootNode(config, notary)
		} else {
			bootNode(config, master)
		}
	}
}
