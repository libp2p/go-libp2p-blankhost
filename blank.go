package blankhost

import (
	"context"
	"io"

	logging "github.com/ipfs/go-log"
	core "github.com/libp2p/go-libp2p-core"
	connmgr "github.com/libp2p/go-libp2p-core/connmgr"
	network "github.com/libp2p/go-libp2p-core/network"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	mstream "github.com/multiformats/go-multistream"
)

var log = logging.Logger("blankhost")

// BlankHost is the thinnest implementation of the host.Host interface
type BlankHost struct {
	n    core.Network
	mux  *mstream.MultistreamMuxer
	cmgr connmgr.ConnManager
}

func NewBlankHost(n core.Network) *BlankHost {
	bh := &BlankHost{
		n:    n,
		cmgr: &connmgr.NullConnMgr{},
		mux:  mstream.NewMultistreamMuxer(),
	}

	n.SetStreamHandler(bh.newStreamHandler)
	return bh
}

var _ core.Host = (*BlankHost)(nil)

func (bh *BlankHost) Addrs() []ma.Multiaddr {
	addrs, err := bh.n.InterfaceListenAddresses()
	if err != nil {
		log.Debug("error retrieving network interface addrs: ", err)
		return nil
	}

	return addrs
}

func (bh *BlankHost) Close() error {
	return bh.n.Close()
}

func (bh *BlankHost) Connect(ctx context.Context, pi core.PeerAddrInfo) error {
	// absorb addresses into peerstore
	bh.Peerstore().AddAddrs(pi.ID, pi.Addrs, pstore.TempAddrTTL)

	cs := bh.n.ConnsToPeer(pi.ID)
	if len(cs) > 0 {
		return nil
	}

	_, err := bh.Network().DialPeer(ctx, pi.ID)
	return err
}

func (bh *BlankHost) Peerstore() pstore.Peerstore {
	return bh.n.Peerstore()
}

func (bh *BlankHost) ID() core.PeerID {
	return bh.n.LocalPeer()
}

func (bh *BlankHost) NewStream(ctx context.Context, p core.PeerID, protos ...core.ProtocolID) (core.Stream, error) {
	s, err := bh.n.NewStream(ctx, p)
	if err != nil {
		return nil, err
	}

	var protoStrs []string
	for _, pid := range protos {
		protoStrs = append(protoStrs, string(pid))
	}

	selected, err := mstream.SelectOneOf(protoStrs, s)
	if err != nil {
		s.Close()
		return nil, err
	}

	selpid := protocol.ID(selected)
	s.SetProtocol(selpid)
	bh.Peerstore().AddProtocols(p, selected)

	return s, nil
}

func (bh *BlankHost) RemoveStreamHandler(p core.ProtocolID) {
	bh.Mux().RemoveHandler(string(p))
}

func (bh *BlankHost) SetStreamHandler(pid core.ProtocolID, handler network.StreamHandler) {
	bh.Mux().AddHandler(string(pid), func(p string, rwc io.ReadWriteCloser) error {
		is := rwc.(core.Stream)
		is.SetProtocol(protocol.ID(p))
		handler(is)
		return nil
	})
}

func (bh *BlankHost) SetStreamHandlerMatch(pid core.ProtocolID, m func(string) bool, handler network.StreamHandler) {
	bh.Mux().AddHandlerWithFunc(string(pid), m, func(p string, rwc io.ReadWriteCloser) error {
		is := rwc.(core.Stream)
		is.SetProtocol(protocol.ID(p))
		handler(is)
		return nil
	})
}

// newStreamHandler is the remote-opened stream handler for core.Network
func (h *BlankHost) newStreamHandler(s core.Stream) {

	protoID, handle, err := h.Mux().Negotiate(s)
	if err != nil {
		log.Warning("protocol mux failed: %s", err)
		s.Close()
		return
	}

	s.SetProtocol(protocol.ID(protoID))

	go handle(protoID, s)
}

// TODO: i'm not sure this really needs to be here
func (bh *BlankHost) Mux() protocol.Switch {
	return bh.mux
}

// TODO: also not sure this fits... Might be better ways around this (leaky abstractions)
func (bh *BlankHost) Network() core.Network {
	return bh.n
}

func (bh *BlankHost) ConnManager() connmgr.ConnManager {
	return bh.cmgr
}
