package asset

import "perun.network/go-perun/channel"

type Register interface {
	GetAsset(name string) channel.Asset
	GetName(asset channel.Asset) string
	GetAllAssets() []channel.Asset
}
