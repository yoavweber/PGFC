package content

import (
	"PGFS/global"
	"context"
	"fmt"
	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"os"
)

/*
	Gets the file for a given CID and stores it in the ContentPath
*/
func GetContent(cid string, node icore.CoreAPI, ctx context.Context) (string, error) {

	var cidPath = path.New(cid)            // Creates path file for the CID
	outputPath := global.ContentPath + cid // File output path

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

/*
	Adds a file for a given filePath to the IPFS network
	storing file within own datastore (found at RepoPath) and notifying the network of it's existence
*/
func AddContent(filePath string, node icore.CoreAPI, ctx context.Context) (string, error) {

	someFile, err := getUnixfsNode(filePath) // Creates a Unixfs Node of the given filePath
	if err != nil {
		return "", fmt.Errorf("could not get File: %s", err)
	}

	cidFile, err := node.Unixfs().Add(ctx, someFile) // Adds the Unixfs Node to the nodes datastore
	if err != nil {
		return "", fmt.Errorf("could not add File: %s", err)
	}

	return cidFile.Cid().String(), nil
}

/*
	Gets the UnixfsNode for a given path
	? Used in the AddContent() function
*/
func getUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path) // Creates a file info of path
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st) // Creates Unixfs Node file
	if err != nil {
		return nil, err
	}

	return f, nil
}
