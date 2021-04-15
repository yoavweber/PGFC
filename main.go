package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

func main() {
	// Basic ipfsnode setup
	r, err := fsrepo.Open("~/.ipfs")
	if err != nil {
		panic(err)
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