package subs

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	Topic string

	Path struct {
		chanelName string
		address    *sdk.Address
	}
)

func PathFromWsMessageInfo(info *sdk.WsMessageInfo) *Path {
	return &Path{
		chanelName: info.ChannelName,
		address:    info.Address,
	}
}

func NewPath(chanelName string, address *sdk.Address) *Path {
	return &Path{
		chanelName: chanelName,
		address:    address,
	}
}

func (t *Path) String() string {
	if t.address == nil {
		return t.chanelName
	}

	return t.chanelName + "/" + t.address.String()
}

func (t *Path) Topic() Topic {
	return Topic(t.chanelName)
}
