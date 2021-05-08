package peers

import (
	"PGFS/api"
	"context"
	"fmt"
	"log"
	"sync"

	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

const repoPath = "./.ipfs/"
const contentPath = "./content/"

func Peer1() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Spawn a node using a temporary path, creating a temporary repo for the run
	fmt.Println("Spawning node on " + repoPath)
	node, err := api.SpawnNode(ctx, repoPath)
	if err != nil {
		panic(err)
	}

	// Node identity information
	fmt.Println("Node spawned on " + repoPath + "\nIdentity information:")
	key, _ := node.Key().Self(ctx)
	fmt.Println(" PeerID: " + key.ID().Pretty() + "\n Path: " + key.Path().String())

	// var bootstrapNodes = []string{
	// 	"/ip4/10.212.137.178/tcp/4001/p2p/12D3KooWDHfFVgZqgRBQRDkYVm9hV8KE6EiaDLTgobHWYn7M62tq",
	// }

	// go connectToPeers(ctx, node, bootstrapNodes)

	addContentPath := contentPath + "test.igc"
	cid, err := api.AddContent(addContentPath, node, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Content added with CID: " + cid)

	fmt.Println("enetring loop")

}

func connectToPeers(ctx context.Context, ipfs icore.CoreAPI, peers []string) error {
	var wg sync.WaitGroup
	peerInfos := make(map[peer.ID]*peerstore.PeerInfo, len(peers))
	for _, addrStr := range peers {
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		pii, err := peerstore.InfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		pi, ok := peerInfos[pii.ID]
		if !ok {
			pi = &peerstore.PeerInfo{ID: pii.ID}
			peerInfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(peerInfos))
	for _, peerInfo := range peerInfos {
		go func(peerInfo *peerstore.PeerInfo) {
			defer wg.Done()
			err := ipfs.Swarm().Connect(ctx, *peerInfo)
			if err != nil {
				log.Printf("failed to connect to %s: %s", peerInfo.ID, err)
			}
		}(peerInfo)
	}
	wg.Wait()
	return nil
}
