// Deprecated: This package has moved into go-libp2p as a sub-package, github.com/libp2p/go-libp2p/p2p/host/blank.
package blankhost

import (
	blankhost "github.com/libp2p/go-libp2p/p2p/host/blank"

	"github.com/libp2p/go-libp2p-core/connmgr"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/network"
)

// BlankHost is the thinnest implementation of the host.Host interface
// Deprecated: use github.com/libp2p/go-libp2p/p2p/host/blank.BlankHost instead.
type BlankHost = blankhost.BlankHost

// Deprecated: use github.com/libp2p/go-libp2p/p2p/host/blank.Option instead.
type Option = blankhost.Option

// Deprecated: use github.com/libp2p/go-libp2p/p2p/host/blank.WithConnectionManager instead.
func WithConnectionManager(cmgr connmgr.ConnManager) Option {
	return blankhost.WithConnectionManager(cmgr)
}

// Deprecated: use github.com/libp2p/go-libp2p/p2p/host/blank.WithEventBus instead.
func WithEventBus(eventBus event.Bus) Option {
	return blankhost.WithEventBus(eventBus)
}

// Deprecated: use github.com/libp2p/go-libp2p/p2p/host/blank.NewBlankHost instead.
func NewBlankHost(n network.Network, options ...Option) *BlankHost {
	return blankhost.NewBlankHost(n, options...)
}
