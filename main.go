package main

import (
	"context"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"io/ioutil"
	"path/filepath"
)

const repoPath = "./.ipfs"

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

	/*

	// Peer list
	var list []icore.ConnectionInfo
	list, _ = node.Swarm().Peers(ctx)
	fmt.Println(list)

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

	// Spawns an IPFS node
	return createNode(ctx, repoPath)

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
func createNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {

	// Opens the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo on node creation: %s", err)
	}

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