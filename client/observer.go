package client

import (
	"github.com/google/uuid"
	"perun.network/go-perun/channel"
)

type Observer interface {
	UpdateState(string)
	UpdateBalance(string)
	GetID() uuid.UUID
}

type Subject interface {
	Register(observer Observer)
	Deregister(observer Observer)

	// NotifyAllState should be called whenever a state change occurs (see go-perun's client.Channel.OnUpdate).
	// The UpdateState method of all registered observers should be called with a text-representation of the new state.
	// Note: tview's DynamicColors are enabled.
	NotifyAllState(from, to *channel.State)

	// NotifyAllBalance should be called whenever the on-chain balance of the client changes. The client is responsible for
	// keeping track of its balance changes. The UpdateBalance method of all registered observers should be called with
	// a text representation of the new balance. tview's DynamicColors are enabled
	// (e.g."[green]420.1337[white] Ada").
	NotifyAllBalance()
}
