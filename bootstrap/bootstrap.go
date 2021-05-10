package bootstrap

import (
	"PGFS/global"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

func addBootstrap(addresses []string) error {
	if fsrepo.IsInitialized(global.RepoPath) { // Checks if repo is initialized
		nodeRepo, err := fsrepo.Open(global.RepoPath) // Opens repo
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

func removeBootstrap() {

}

func clearBootstrap() error {
	if fsrepo.IsInitialized(global.RepoPath) { // Checks if repo is initialized
		nodeRepo, err := fsrepo.Open(global.RepoPath) // Opens repo
		if err != nil {
			return fmt.Errorf("failed to open repo when adding bootstrap: %s", err)
		}

		var cfg *config.Config // Defines config file

		cfg, err = nodeRepo.Config() // Receives current config
		if err != nil {
			return fmt.Errorf("failed to open repo when adding bootstrap: %s", err)
		}

		cfg.Bootstrap = nil

		err = nodeRepo.SetConfig(cfg) // Sets current config with updated config
		if err != nil {
			return fmt.Errorf("failed to set config when adding bootstrap: %s", err)
		}
		return nil
	} else {
		return fmt.Errorf("cannot add bootstrap to an uninitialized node")
	}
}
