package content

import (
	"PGFS/global"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	igc "github.com/marni/goigc"
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

	someFile, dir, err := WrapContent(filePath) // Wraps content and opens it as a file

	cidFile, err := node.Unixfs().Add(ctx, someFile) // Adds the Unixfs Node to the nodes datastore

	someFile.Close() // closes file

	if err != nil { // Content could not be added, error occured
		osErr := os.RemoveAll(dir) // Attempts to remove temp dir
		if osErr != nil {
			panic(fmt.Errorf("failed to remove temp package wrapper: %s", osErr)) // error occured whilst removing temp dir, PANIC!!!
		}
		return "", fmt.Errorf("could not add File: %s", err)
	}

	err = os.RemoveAll(dir)
	if err != nil {
		return "", fmt.Errorf("failed to remove temp package wrapper: %s", err)
	}

	return cidFile.Cid().String(), nil
}

/*
	Wraps the content attaching a metadata.txt file to it for searching purposes
*/
func WrapContent(filePath string) (files.Node, string, error) {

	stamp := int(time.Now().UnixNano() / int64(time.Millisecond))     // Creates unique stamp
	time.Sleep(1 * time.Millisecond)                                  // Assure uniqueness
	dir := global.TempContentPath + "pkg" + strconv.Itoa(stamp) + "/" // Temp wrapper dir path

	err := os.Mkdir(dir, 0755) // Makes temp package dir on dir path
	if err != nil {
		return nil, dir, fmt.Errorf("failed creating wrapper package: %s", err)
	}

	fileIn, err := os.Open(filePath) // Opens the given file
	if err != nil {
		return nil, dir, fmt.Errorf("failed opening given file: %s", err)
	}

	fileOut, err := os.Create(dir + "data.igc") // Creates a data.igc file in the package directory
	if err != nil {
		fileIn.Close()
		return nil, dir, fmt.Errorf("failed creating data.igc file: %s", err)
	}

	_, err = io.Copy(fileOut, fileIn) // Copies the file into the temporary package directory
	if err != nil {
		fileIn.Close() // Closes file
		fileOut.Close()
		return nil, dir, fmt.Errorf("failed copying file to wrapper: %s", err)
	}

	fileIn.Close() // Closes file
	fileOut.Close()

	track, err := igc.ParseLocation(filePath) // getting the location igc
	fileData := fmt.Sprintf("Pilot: %s, gliderType: %s, date: %s",
		track.Pilot, track.GliderType, track.Date.String())

	err = ioutil.WriteFile(dir+"metadata.txt", []byte(fileData), 0755) // Creates a metadata file with the location

	if err != nil {
		return nil, dir, fmt.Errorf("failed creating metadata file: %s", err)
	}

	someFile, err := getUnixfsNode(dir) // Creates a Unixfs Node of the given filePath
	if err != nil {
		return nil, dir, fmt.Errorf("could not get File: %s", err)
	}

	return someFile, dir, nil
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
