package main

import (
	"PGFS/global"
	"PGFS/node"
	"PGFS/peers"
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"log"
	"strconv"
	"time"
)

func main() {

	/*
		Initializes the context
	*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/*
		Spawns a node
	*/
	log.Println("Spawning node on " + global.RepoPath)
	node, err := spawnNode(ctx, false)
	if err != nil {
		panic(err)
	}
	log.Println("Node spawned on " + global.RepoPath)

	log.Println("")

	/*
		Node identity information
	*/
	id, err := peers.GetPeerID()
	if err != nil {
		panic(err)
	}
	log.Println("* Identity information:")
	log.Println("* PeerID: " + id)

	log.Println("")

	/*
		Connects to bootstrap node/s
	*/

	go peers.ConnectToPeers(ctx, node)

	/*
		Retrieves the peer list
	*/
	peerList, err := peers.ListAllPeers(node, ctx)
	log.Println("? Peer list:")
	if len(peerList) > 0 {
		for i := range peerList {
			var peer icore.ConnectionInfo
			peer = peerList[i]
			log.Println("? Peer #" + strconv.Itoa(i) + ": " + peer.ID().Pretty())
			log.Println("?  Address: " + peer.Address().String())
			log.Println("?  Direction: " + peer.Direction().String())
		}
	} else {
		log.Println("!  No peers found :(")
	}

	time.Sleep(10 * time.Second)

}

/*
	Spawns a node
*/
func spawnNode(ctx context.Context, isServer bool) (icore.CoreAPI, error) {
	if err := node.SetupPlugins(""); err != nil {
		return nil, err
	}

	// Checks if repo is initialized
	if !fsrepo.IsInitialized(global.RepoPath) {
		// Initializes repo in repoPath
		if err := node.RepoInit(); err != nil {
			return nil, err
		}
	}

	// Opens the repo
	nodeRepo, err := fsrepo.Open(global.RepoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo on node creation: %s", err)
	}

	// Spawns an IPFS node
	return node.CreateNode(ctx, nodeRepo, isServer)

}
