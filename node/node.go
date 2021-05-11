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

func RepoInit() error {

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return fmt.Errorf("failed to initialize repo config: %s", err)
	}

	// Assigns custom bootstrap to repo
	cfg.Bootstrap = nil

	cfg.Swarm.AddrFilters = nil


	// Create the repo with the config
	err = fsrepo.Init(global.RepoPath, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize repo: %s", err)
	}


	swarmKeyPath := "./global/swarm.key"
	if _, err := os.Stat(swarmKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("failed to locate swarm.key at %s: %s", swarmKeyPath, err)
	}
	swarmKeyIn, err := os.Open(swarmKeyPath)
	if err != nil {
		return fmt.Errorf("could not open ./global/swarm.key: %s", err)
	}

	swarmKeyOut, err := os.Create(global.RepoPath + "swarm.key")
	if err != nil {
		return fmt.Errorf("could not create %s swarm.key: %s", global.RepoPath, err)
	}

	_, err = io.Copy(swarmKeyOut, swarmKeyIn)
	if err != nil {
		return fmt.Errorf("could not copy swarm.key to new repo location: %s", err)
	}

	return nil
}

// Creates an IPFS node and returns its coreAPI
func CreateNode(ctx context.Context, repo repo.Repo, isServer bool) (icore.CoreAPI, error) {


	// Build configurations of the node
	nodeOptions := &core.BuildCfg{
		Online:  true,
		Repo: repo,
	}

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
