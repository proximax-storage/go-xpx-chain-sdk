package websocket

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"golang.org/x/net/websocket"
)

func newMessagePublisher(conn *websocket.Conn) messagePublisher {
	return &catapultWebsocketMessagePublisher{
		conn: conn,
	}
}

type messagePublisher interface {
	PublishSubscribeMessage(uid string, path pathType) error
	PublishUnsubscribeMessage(uid string, path pathType) error
}

type catapultWebsocketMessagePublisher struct {
	conn *websocket.Conn
}

func (p *catapultWebsocketMessagePublisher) PublishSubscribeMessage(uid string, path pathType) error {
	dto := &sdk.SubscribeDTO{
		Uid:       uid,
		Subscribe: string(path),
	}

	if err := websocket.JSON.Send(p.conn, dto); err != nil {
		return errors.Wrap(err, "error publishing subscribe message into websocket connection")
	}

	return nil
}

func (p *catapultWebsocketMessagePublisher) PublishUnsubscribeMessage(uid string, path pathType) error {
	dto := &sdk.UnsubscribeDTO{
		Uid:         uid,
		Unsubscribe: string(path),
	}

	if err := websocket.JSON.Send(p.conn, dto); err != nil {
		return errors.Wrap(err, "error publishing unsubscribe message into websocket connection")
	}

	return nil
}
