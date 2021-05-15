module PGFS

go 1.15

require (
	github.com/ipfs/go-ipfs v0.8.0
	github.com/ipfs/go-ipfs-config v0.12.0
	github.com/ipfs/go-ipfs-files v0.0.8
	github.com/ipfs/interface-go-ipfs-core v0.4.0
	github.com/libp2p/go-libp2p-core v0.8.5
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/marni/goigc v0.1.0
	github.com/multiformats/go-multiaddr v0.3.1
)

replace github.com/ipfs/go-ipfs => github.com/yoavweber/go-pgfs v0.8.1-0.20210515185130-524bcb2b929b
