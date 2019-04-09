package sdk

import "golang.org/x/net/websocket"

func newMessagePublisher(conn *websocket.Conn) messagePublisher {
	return &catapultWebsocketMessagePublisher{
		conn: conn,
	}
}

type messagePublisher interface {
	PublishSubscribeMessage(uid string, path string) error
	PublishUnsubscribeMessage(uid string, path string) error
}

type catapultWebsocketMessagePublisher struct {
	conn *websocket.Conn
}

func (p *catapultWebsocketMessagePublisher) PublishSubscribeMessage(uid string, path string) error {
	dto := &subscribeDTO{
		Uid:       uid,
		Subscribe: path,
	}

	if err := websocket.JSON.Send(p.conn, dto); err != nil {
		return err
	}

	return nil
}

func (p *catapultWebsocketMessagePublisher) PublishUnsubscribeMessage(uid string, path string) error {
	dto := &unsubscribeDTO{
		Uid:         uid,
		Unsubscribe: path,
	}

	if err := websocket.JSON.Send(p.conn, dto); err != nil {
		return err
	}

	return nil
}
