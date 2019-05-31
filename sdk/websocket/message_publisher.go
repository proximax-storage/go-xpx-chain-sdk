package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"sync"
)

func newMessagePublisher(conn *websocket.Conn) MessagePublisher {
	return &catapultWebsocketMessagePublisher{
		conn: conn,
	}
}

type MessagePublisher interface {
	PublishSubscribeMessage(uid string, path Path) error
	PublishUnsubscribeMessage(uid string, path Path) error
	SetConn(conn *websocket.Conn)
}

type catapultWebsocketMessagePublisher struct {
	sync.Mutex
	conn *websocket.Conn
}

func (p *catapultWebsocketMessagePublisher) PublishSubscribeMessage(uid string, path Path) error {
	p.Lock()
	defer p.Unlock()

	dto := &subscribeDTO{
		Uid:       uid,
		Subscribe: string(path),
	}

	if err := p.conn.WriteJSON(dto); err != nil {
		return errors.Wrap(err, "publishing subscribe message into websocket connection")
	}

	return nil
}

func (p *catapultWebsocketMessagePublisher) PublishUnsubscribeMessage(uid string, path Path) error {
	p.Lock()
	defer p.Unlock()

	dto := &unsubscribeDTO{
		Uid:         uid,
		Unsubscribe: string(path),
	}

	if err := p.conn.WriteJSON(dto); err != nil {
		return errors.Wrap(err, "publishing unsubscribe message into websocket connection")
	}

	return nil
}

func (p *catapultWebsocketMessagePublisher) SetConn(conn *websocket.Conn) {
	p.conn = conn
}

type subscribeDTO struct {
	Uid       string `json:"uid"`
	Subscribe string `json:"subscribe"`
}

type unsubscribeDTO struct {
	Uid         string `json:"uid"`
	Unsubscribe string `json:"unsubscribe"`
}

type wsConnectionResponse struct {
	Uid string `json:"uid"`
}
