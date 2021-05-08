package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
)

func SpawnNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {
	if err := setupPlugins(""); err != nil {
		return nil, err
	}

	// Checks if repo is initialized
	if !fsrepo.IsInitialized(repoPath) {
		// Initializes repo in repoPath
		if err := repoInit(repoPath); err != nil {
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

func repoInit(repoPath string) error {

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

func createNode(ctx context.Context, repo repo.Repo) (icore.CoreAPI, error) {

	// Build configurations of the node
	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTServerOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
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

func addBootstrap(repoPath string, addresses []string) error {
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
