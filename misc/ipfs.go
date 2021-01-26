package main

import (
	"context"
	"fmt"
	"path/filepath"

	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	icore "github.com/ipfs/interface-go-ipfs-core"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)


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

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

// Spawns a node on the default repo location, if the repo exists
func spawnDefault(ctx context.Context) (icore.CoreAPI, error) {
	defaultPath, err := config.PathRoot()
	if err != nil {
		// shouldn't be possible
		return nil, err
	}
	fmt.Println("repo: ", defaultPath)

	if err := setupPlugins(defaultPath); err != nil {
		return nil, err

	}

	return createNode(ctx, defaultPath)
}


func main() {
	/// --- Part I: Getting a IPFS node running
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	// Spawn a node using the default path (~/.ipfs), assuming that a repo exists there already
	ipfs, err := spawnDefault(ctx)
	if err != nil {
		fmt.Println("No IPFS repo available on the default path: ", err.Error())
		panic(err)
	}

	// 字符串写入ipfs
	someContent := files.NewBytesFile([]byte("test2")) 
	cidFile, err := ipfs.Unixfs().Add(ctx, someContent)
	if err != nil {
		panic(fmt.Errorf("Could not add File: %s", err))
	}
	fmt.Println("cid: ", cidFile.String())

	//outputBasePath := "./data/"
	//outputPathFile := outputBasePath + strings.Split(cidFile.String(), "/")[2]

	// 获取文件
	rootNodeFile, err := ipfs.Unixfs().Get(ctx, cidFile)
	if err != nil {
		panic(fmt.Errorf("Could not get file with CID: %s", err))
	}

	// 文件大小
	size, _ := rootNodeFile.Size()

	// 读出内容
	longBuf := make([]byte, size)
	if _, err := rootNodeFile.(files.File).Read(longBuf); err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("content: %v\n", string(longBuf))

	//err = files.WriteTo(rootNodeFile, outputPathFile)
	//if err != nil {
	//	panic(fmt.Errorf("Could not write out the fetched CID: %s", err))
	//}

}