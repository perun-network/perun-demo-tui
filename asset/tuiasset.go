package asset

import "perun.network/go-perun/channel"

type TUIAsset struct {
	channel.Asset
	Name string
}

func MakeTUIAsset(asset channel.Asset, name string) TUIAsset {
	return TUIAsset{Asset: asset, Name: name}
}
