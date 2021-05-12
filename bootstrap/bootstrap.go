package bootstrap

import (
	"PGFS/global"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

/*
	Adds a given bootstrap/s to the list defined in the config file located at RepoPath
*/
func AddBootstrap(addresses []string) error {
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

/*
	Removes a given bootstrap from the list defined in the config file located at RepoPath
*/
func RemoveBootstrap(address string) error {
	if fsrepo.IsInitialized(global.RepoPath) { // Checks if repo is initialized
		nodeRepo, err := fsrepo.Open(global.RepoPath) // Opens repo
		if err != nil {
			return fmt.Errorf("failed to open repo when removing bootstrap: %s", err)
		}

		var cfg *config.Config // Defines config file

		cfg, err = nodeRepo.Config() // Receives current config
		if err != nil {
			return fmt.Errorf("failed to open repo when removing bootstrap: %s", err)
		}

		if len(cfg.Bootstrap) > 0 {
			for i := range cfg.Bootstrap {
				if cfg.Bootstrap[i] == address {
					cfg.Bootstrap = append(cfg.Bootstrap[:i], cfg.Bootstrap[i+1:]...) // Removes bootstrap/s from the config
					break
				} else if i == len(cfg.Bootstrap) {
					return fmt.Errorf("bootstrap " + address + " wasn't removed, address not found in bootstrap list")
				}
			}
		} else {
			return fmt.Errorf("bootstrap " + address + " wasn't removed, bootstrap list empty")
		}

		err = nodeRepo.SetConfig(cfg) // Sets current config with updated config
		if err != nil {
			return fmt.Errorf("failed to set config when removing bootstrap: %s", err)
		}
		return nil
	} else {
		return fmt.Errorf("cannot remove bootstrap from an uninitialized node")
	}
}

/*
	Clears the bootstrap list of the config file located at RepoPath
*/
func ClearBootstrap() error {
	if fsrepo.IsInitialized(global.RepoPath) { // Checks if repo is initialized
		nodeRepo, err := fsrepo.Open(global.RepoPath) // Opens repo
		if err != nil {
			return fmt.Errorf("failed to open repo when clearing bootstrap list: %s", err)
		}

		var cfg *config.Config // Defines config file

		cfg, err = nodeRepo.Config() // Receives current config
		if err != nil {
			return fmt.Errorf("failed to open repo when clearing bootstrap list: %s", err)
		}

		cfg.Bootstrap = nil

		err = nodeRepo.SetConfig(cfg) // Sets current config with updated config
		if err != nil {
			return fmt.Errorf("failed to set config when clearing bootstrap list: %s", err)
		}
		return nil
	} else {
		return fmt.Errorf("cannot clear the bootstrap list on an uninitialized node")
	}
}

/*
	Retrieves a list of the bootstraps defined in the config
*/
func GetBootstrapList() ([]string, error) {
	if fsrepo.IsInitialized(global.RepoPath) { // Checks if repo is initialized
		nodeRepo, err := fsrepo.Open(global.RepoPath) // Opens repo
		if err != nil {
			return nil, fmt.Errorf("failed to open repo when clearing bootstrap list: %s", err)
		}

		var cfg *config.Config // Defines config file

		cfg, err = nodeRepo.Config() // Receives current config
		if err != nil {
			return nil, fmt.Errorf("failed to open repo when clearing bootstrap list: %s", err)
		}

		return cfg.Bootstrap, nil
	} else {
		return nil, fmt.Errorf("cannot clear the bootstrap list on an uninitialized node")
	}
}
