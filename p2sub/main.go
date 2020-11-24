// Copyright 2019 P2Sub Authors
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

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"github.com/p2sub/p2sub/keypair"
)

func main() {
	// Init and parse configurations
	Init()

	// Main context
	ctx := context.Background()

	// Create multiaddress from given string
	bindPort := conf.GetBindPort()
	bindHost := conf.GetBindHost()
	bindStr := fmt.Sprintf("/ip4/%s/tcp/%d", bindHost, bindPort)
	sugar.Debugf("Bind address: %s", bindStr)
	sourceMultiAddr, err := multiaddr.NewMultiaddr(bindStr)
	if err != nil {
		sugar.Panic(err)
	}

	// Generate or load existing key pair
	nodeConfigFile := conf.GetKeyFile()
	nodeKey := new(keypair.KeyPair)
	if _, err := os.Stat(nodeConfigFile); err != nil {
		// Create a new key pair
		nodeKey, err = keypair.New()
		if err != nil {
			panic(err)
		}
		sugar.Debugf("Save key to file: %s", nodeConfigFile)
		nodeKey.SaveToFile(nodeConfigFile)
	} else {
		// Load key from json file if existed
		nodeKey, err = keypair.LoadFromFile(nodeConfigFile)
		sugar.Debugf("Load key from file: %s", nodeConfigFile)
		if err != nil {
			panic(err)
		}
	}

	//Setup host with key
	nodeID, _ := nodeKey.GetID()
	sugar.Debugf("Setup host with given private key, node ID: %s", nodeID)
	prvKey := nodeKey.GetPrivateKey()
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}

	// Start new gossip pub sub
	myPubsub, err := pubsub.NewGossipSub(
		ctx,
		host,
		pubsub.WithPeerExchange(true),
	)

	// Start direct connect if direct connect was set
	directConnection := conf.GetDirectConnect()
	if directConnection != "" {
		sugar.Infof("Boot node is: %s", directConnection)
		// Turn the destination into a multiaddr.
		mAddr, err := multiaddr.NewMultiaddr(directConnection)
		if err != nil {
			log.Fatalln(err)
		}

		// Extract the peer ID from the multiaddr.
		info, err := peer.AddrInfoFromP2pAddr(mAddr)
		if err != nil {
			log.Fatalln(err)
		}

		// Connect to node
		host.Connect(ctx, *info)
	}

	// Detect other nodes by domain
	domain := conf.GetDomain()
	if domain != "" {
		// Start a DHT, for use in peer discovery. We can't just make a new DHT
		// client because we want each peer to maintain its own local copy of the
		// DHT, so that the bootstrapping node of the DHT can go down without
		// inhibiting future peer discovery.
		kademliaDHT, err := dht.New(ctx, host)
		if err != nil {
			panic(err)
		}

		// Bootstrap the DHT. In the default configuration, this spawns a Background
		// thread that will refresh the peer table every five minutes.
		sugar.Debug("Bootstrapping the DHT")
		if err = kademliaDHT.Bootstrap(ctx); err != nil {
			panic(err)
		}

		// Let's connect to the bootstrap nodes first. They will tell us about the
		// other nodes in the network.
		var wg sync.WaitGroup
		for _, peerAddr := range dht.DefaultBootstrapPeers {
			peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := host.Connect(ctx, *peerinfo); err != nil {
					sugar.Warn(err)
				} else {
					sugar.Infof("Connection established with bootstrap node: %v", *peerinfo)
				}
			}()
		}
		wg.Wait()

		// We use a rendezvous point `domain` to announce our location.
		// This is like telling your friends to meet you at the Eiffel Tower.
		domain := conf.GetDomain()
		sugar.Info("Announcing ourselves...")
		routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
		discovery.Advertise(ctx, routingDiscovery, domain)
		sugar.Debug("Successfully announced!")

		// Now, look for others who have announced
		// This is like your friend telling you the location to meet you.
		sugar.Debug("Searching for other peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, domain)
		if err != nil {
			panic(err)
		}

		for curPeer := range peerChan {
			if curPeer.ID == host.ID() {
				continue
			}

			sugar.Debugf("Connecting to: %s", curPeer.ID.Pretty())
			err := host.Connect(ctx, curPeer)

			if err != nil {
				//sugar.Warnf("Connection failed: %v", err)
				continue
			}

			sugar.Infof("Connected to: %s", curPeer.ID.Pretty())
		}
	}

	myPubsub.GetTopics()

	topic, _ := myPubsub.Join("hello")
	helloWorld, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	for {
		time.AfterFunc(time.Duration(rand.Intn(10))*time.Second, func() {
			topic.Publish(ctx, []byte(nodeConfigFile))
		})
		msg, err := helloWorld.Next(ctx)
		if err != nil {
			panic(err)
		}
		sugar.Debugf("Topic: %s from: %s data: %s", topic.String(), msg.GetFrom().String(), string(msg.GetData()))
	}

}
