package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
)

const repoPath = "./.ipfs/"
const contentPath = "./content/"

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Spawn a node using a temporary path, creating a temporary repo for the run
	fmt.Println("Spawning node on " + repoPath)
	node, err := spawnNode(ctx)
	if err != nil {
		panic(err)
	}

	// Node identity information
	fmt.Println("Node spawned on " + repoPath + "\nIdentity information:")
	key, _ := node.Key().Self(ctx)
	fmt.Println(" PeerID: " + key.ID().Pretty() + "\n Path: " + key.Path().String())

	var bootstrapNodes = []string {
		"/ip4/10.212.137.178/tcp/4001/p2p/12D3KooWDHfFVgZqgRBQRDkYVm9hV8KE6EiaDLTgobHWYn7M62tq",
	}

	go connectToPeers(ctx, node, bootstrapNodes)


	/*

	cid := "QmRJE9bXiKrR3EZSmw9xYC1dGL8oK5NjTYkCGATw4Gq8rn"

	filePath, err := getContent(cid, node, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Content with CID: " + cid + "\nreceived and written to " + filePath)



	 */


	addContentPath := contentPath + "test.txt"
	cid, err := addContent(addContentPath, node, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Content added with CID: " + cid)


	/*

	filePath, err := getContent(cid, node, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Content with CID: " + cid + "\nreceived and written to " + filePath)



	 */



}

// Spawns a node
func spawnNode(ctx context.Context) (icore.CoreAPI, error) {
	if err := setupPlugins(""); err != nil {
		return nil, err
	}

	// Checks if repo is initialized
	if !fsrepo.IsInitialized(repoPath) {
		// Initializes repo in repoPath
		if err := repoInit(); err != nil {
			return nil, err
		}
	}

	// Opens the repo
	nodeRepo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo on node creation: %s", err)
	}

	// Spawns an IPFS node
	return createNode(ctx, nodeRepo)

}





