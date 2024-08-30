package subs

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	Topic string

	Path struct {
		chanelName Topic
		address    *sdk.Address
	}
)

func PathFromWsMessageInfo(info *sdk.WsMessageInfo) *Path {
	return &Path{
		chanelName: Topic(info.ChannelName),
		address:    info.Address,
	}
}

func NewPath(chanelName Topic, address *sdk.Address) *Path {
	return &Path{
		chanelName: chanelName,
		address:    address,
	}
}

func (t *Path) String() string {
	if t.address == nil {
		return string(t.chanelName)
	}

	return string(t.chanelName) + "/" + t.address.Address
}

func (t *Path) Topic() Topic {
	return Topic(t.chanelName)
}

func (t *Path) Address() *sdk.Address {
	if t.address == nil {
		return &sdk.Address{}
	}

	return t.address
}
