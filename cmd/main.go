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
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Spawn a node using a temporary path, creating a temporary repo for the run
	fmt.Println("Spawning node on " + global.RepoPath)
	node, err := spawnNode(ctx, false)
	if err != nil {
		panic(err)
	}

	// Node identity information
	fmt.Println("Node spawned on " + global.RepoPath + "\nIdentity information:")
	key, _ := node.Key().Self(ctx)
	fmt.Println(" PeerID: " + key.ID().Pretty() + "\n Path: " + key.Path().String())

	var bootstrapNodes = []string {
		"/ip4/10.22.201.110/tcp/4001/p2p/QmdoF64z3SrpWEJphCq3pPq5z3y5y3Gu2KePhFti6R7L9i",
	}

	go peers.ConnectToPeers(ctx, node, bootstrapNodes)

	addContentPath := global.ContentPath + "2021-04-11-XCT-XXX-02.igc"
	cid, err := content.AddContent(addContentPath, node, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Content added with CID: " + cid)


	time.Sleep(30*time.Second)







/*




 */
	/*

	cid := "QmXV4rfFBLbsM1q2TyzmEM8Pkb7Pp3dHAYf737ScSJQZcv"

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





