package main

import (
	"PGFS/content"
	"PGFS/global"
	"PGFS/node"
	"PGFS/peers"
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"log"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Spawn a node using a temporary path, creating a temporary repo for the run
	log.Println("Spawning node on " + global.RepoPath)
	node, err := spawnNode(ctx, false)
	if err != nil {
		panic(err)
	}

	// Node identity information
	log.Println("Node spawned on " + global.RepoPath + "\nIdentity information:")
	key, _ := node.Key().Self(ctx)
	log.Println(" PeerID: " + key.ID().Pretty() + "\n Path: " + key.Path().String())

	var bootstrapNodes = []string {
		"/ip4/server/tcp/4001/ipfs/QmNrv9UcFRhG6ToxcSAdvNBkPZZv3Yp8xvAFixnLHzCLow",
		//"/ip4/10.22.201.110/tcp/4001/ipfs/QmXxHTom3PepoW3VrGDvmU89EKah8qRSsoZ1WLBki7w63i",
		//10.212.139.99
	}

	go peers.ConnectToPeers(ctx, node, bootstrapNodes)

	peerList, err := peers.ListAllPeers(node, ctx)
	log.Println(peerList)

	addContentPath := global.ContentPath + "test.txt"
	cid, err := content.AddContent(addContentPath, node, ctx)
	if err != nil {
		log.Println(err)
	}

	addContentPath = global.ContentPath + "2021-05-12-XCT-XXX-01.igc"
	cid, err = content.AddContent(addContentPath, node, ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Content added with CID: " + cid)

	time.Sleep(10*time.Second)

	/*

		cid = "QmS98pgfsLTc91kjHDzb5V9nCwXbs9pe2h2Kj8zsVCXEmR"

		filePath, err := content.GetContent(cid, node, ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println("Content with CID: " + cid + "\nreceived and written to " + filePath)



		filePath, err := getContent(cid, node, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Content with CID: " + cid + "\nreceived and written to " + filePath)




	 */



}

// Spawns a node
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





