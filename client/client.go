package client

import (
	"perun.network/go-perun/wire"
	"perun.network/perun-demo-tui/asset"
)

type DemoClient interface {
	Subject

	// DisplayName returns the name of the client (e.g. "Alice").
	DisplayName() string

	// DisplayAddress returns the address of the client (favourably in a way that can be copied to a block explorer).
	DisplayAddress() string

	// WireAddress returns the wire address of the client. It is used to e.g. open a channel with the client.
	WireAddress() wire.Address

	// OpenChannel should propose and open a channel with the specified peer.
	// The amount specifies the amount of funds per asset that this client should deposit into the channel.
	// Note: The Demo currently also does not include any visualization of
	// proposals. One option is to have the peer deposit the same amount of funds and implement the ProposalHandler
	// accordingly.
	OpenChannel(address wire.Address, amount map[asset.TUIAsset]float64)

	// SendPaymentToPeer should send a payment to the peer of the client in the client's currently open channel.
	// Note: For now we assume that the client only has one open channel.
	// TODO: Move this functionality into a `Channel` interface.
	SendPaymentToPeer(amounts map[asset.TUIAsset]float64)

	// Settle should settle the client's currently open channel.
	// Note: For now we assume that the client only has one open channel.
	// TODO: Move this functionality into a `Channel` interface.
	Settle()

	// HasOpenChannel returns true iff the client has an open channel.
	HasOpenChannel() bool

	// GetOpenChannelAssets returns the assets of the client's currently open channel.
	GetOpenChannelAssets() []asset.TUIAsset
}
