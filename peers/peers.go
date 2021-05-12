package peers

import (
	"PGFS/global"
	"context"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"log"
	"sync"
)

func ConnectToPeers(ctx context.Context, ipfs icore.CoreAPI, peers []string) error {
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


func ListAllPeers(node icore.CoreAPI, ctx context.Context) ([]icore.ConnectionInfo, error) {

	// Peer list
	var list []icore.ConnectionInfo

	list, err := node.Swarm().Peers(ctx) // Swarm peers

	return list, err
}

func GetIdentity(node icore.CoreAPI, ctx context.Context) (string,error) {
	// Node identity information
	if fsrepo.IsInitialized(global.RepoPath) { // Checks if repo is initialized
		nodeRepo, err := fsrepo.Open(global.RepoPath) // Opens repo

		var cfg *config.Config // Defines config file

		cfg, err = nodeRepo.Config() // Receives current config
		if err != nil {
			return "", fmt.Errorf("failed to open repo when adding bootstrap: %s", err)
		}

		return cfg.Identity.PeerID, nil
	} else {
		return "", fmt.Errorf("cannot add bootstrap to an uninitialized node")
	}

}