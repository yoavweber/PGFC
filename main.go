package main

import (
	"context"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
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

	/*
	// Peer list
	var list []icore.ConnectionInfo
	list, _ = node.Swarm().Peers(ctx)
	fmt.Println(list)
	*/


	addContentPath := contentPath + "test.txt"
	cid, err := addContent(addContentPath, node, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Content added with CID: " + cid)


	filePath, err := getContent(cid, node, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Content with CID: " + cid + "\nreceived and written to " + filePath)

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

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

func repoInit() error {

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return fmt.Errorf("failed to initialize repo config: %s", err)
	}

	// Assigns custom bootstrap to repo
	// TODO: Add bootstrap NTNU node as a bootstrap node
	cfg.Bootstrap = nil

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize repo: %s", err)
	}

	return nil
}

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repo repo.Repo) (icore.CoreAPI, error) {

	// Build configurations of the node
	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		// There is also an option called: libp2p.DHTServerOption
		Repo: repo,
	}

	// Creates new node
	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create new node: %s", err)
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

func addContent(filePath string, node icore.CoreAPI, ctx context.Context) (string, error) {

	someFile, err := getUnixfsNode(filePath)
	if err != nil {
		return "", fmt.Errorf("could not get File: %s", err)
	}

	cidFile, err := node.Unixfs().Add(ctx, someFile)
	if err != nil {
		return "", fmt.Errorf("could not add File: %s", err)
	}

	return cidFile.Cid().String(), nil
}

func getUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func getContent(cid string, node icore.CoreAPI, ctx context.Context) (string, error) {

	var cidPath = path.New(cid)
	outputPath := contentPath + cid // File output path

	fileNode, err := node.Unixfs().Get(ctx, cidPath) // Gets the node associated with the CID given
	if err != nil {
		return "", fmt.Errorf("could not get file with CID: %s", err)
	}

	err = files.WriteTo(fileNode, outputPath) // Writes fetched file to output path
	if err != nil {
		return "", fmt.Errorf("could not write out the fetched CID: %s", err)
	}

	return outputPath, nil // returns output path
}

func addBootstrap(addresses []string) error {
	if fsrepo.IsInitialized(repoPath) { // Checks if repo is initialized
		nodeRepo, err := fsrepo.Open(repoPath) // Opens repo
		if err != nil {
			return fmt.Errorf("failed to open repo when adding bootstrap: %s", err)
		}

		var cfg *config.Config // Defines config file

		cfg, err = nodeRepo.Config() // Receives current config
		if err != nil {
			return fmt.Errorf("failed to open repo when adding bootstrap: %s", err)
		}

		cfg.Bootstrap = append(cfg.Bootstrap, addresses...) // Appends bootstrap/s onto config

		err = nodeRepo.SetConfig(cfg) // Sets current config with updated config
		if err != nil {
			return fmt.Errorf("failed to set config when adding bootstrap: %s", err)
		}
		return nil
	} else {
		return fmt.Errorf("cannot add bootstrap to an uninitialized node")
	}

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
