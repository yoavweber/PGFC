package node

import (
	"PGFS/global"
	"context"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

/*
	Sets up the plugins
*/
func SetupPlugins(externalPluginsPath string) error {
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

/*
	Initializes the repository at RepoPath
*/
func RepoInit() error {

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return fmt.Errorf("failed to initialize repo config: %s", err)
	}

	// Removes all default bootstraps from list
	// We connect to our bootstrap/s in main.go
	cfg.Bootstrap = nil

	// Removes address filters allowing for "local discovery" mode (file sharing through localhost)
	// This should be changed to the default setting if one wishes to allow the node to discover nodes with IP address outside of it's own network
	cfg.Swarm.AddrFilters = nil

	// Create the repo with the config
	err = fsrepo.Init(global.RepoPath, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize repo: %s", err)
	}

	// Checks if swarm.key exists on it's defined SwarmKeyPath
	if _, err := os.Stat(global.SwarmKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("failed to locate swarm.key at %s: %s", global.SwarmKeyPath, err)
	}
	// Current location of swarm.key
	swarmKeyIn, err := os.Open(global.SwarmKeyPath)
	if err != nil {
		return fmt.Errorf("could not open ./global/swarm.key: %s", err)
	}

	// New location of swarm.key (In newly initialized repository)
	swarmKeyOut, err := os.Create(global.RepoPath + "swarm.key")
	if err != nil {
		return fmt.Errorf("could not create %s swarm.key: %s", global.RepoPath, err)
	}

	// Copying swarm.key from old to new location
	_, err = io.Copy(swarmKeyOut, swarmKeyIn)
	if err != nil {
		return fmt.Errorf("could not copy swarm.key to new repo location: %s", err)
	}

	return nil
}

/*
	Creates an IPFS node and returns its coreAPI (client)
*/
func CreateNode(ctx context.Context, repo repo.Repo, isServer bool) (icore.CoreAPI, error) {

	// Build configurations of the node
	nodeOptions := &core.BuildCfg{
		Online: true,
		Repo:   repo,
	}

	// Checks if node wishes to be a DHT Server and sets routing appropriately
	if isServer {
		nodeOptions.Routing = libp2p.DHTServerOption
	} else {
		nodeOptions.Routing = libp2p.DHTOption
	}

	// Creates new node
	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create new node: %s", err)
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}
