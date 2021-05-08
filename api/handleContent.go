package api

import (
	"context"
	"fmt"
	"os"

	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

func AddContent(filePath string, node icore.CoreAPI, ctx context.Context) (string, error) {

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

func GetContent(cid string, node icore.CoreAPI, contentFolder string, ctx context.Context) (string, error) {

	var cidPath = path.New(cid)
	outputPath := contentFolder + cid // File output path

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
