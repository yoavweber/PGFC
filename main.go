package main

import (
	"context"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"io/ioutil"
)

const repoPath = "./.ipfs"

func main() {
	// Basic ipfsnode setup
	r, err := fsrepo.Open(repoPath) // Attempts to open IPFS repo
	if err != nil {
		if fsrepo.IsInitialized(repoPath) { // Checks if error comes from uninitialized repo
			panic(err)
		} else {

			// Create a config with default options and a 2048 bit key
			cfg, err := config.Init(ioutil.Discard, 2048)
			if err != nil {
				panic(err)
			}
			// Initialize repo at path
			err = fsrepo.Init(repoPath, cfg)
			if err != nil {
				panic(err)
			}
		}

	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &core.BuildCfg{
		Repo:   r,
		Online: true,
	}

	nd, err := core.NewNode(ctx, cfg)

	if err != nil {
		panic(err)
	}

	fmt.Println(nd)
}