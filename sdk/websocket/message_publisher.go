package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type MessagePublisher interface {
	PublishSubscribeMessage(uid string, path string) error
	PublishUnsubscribeMessage(uid string, path string) error
	SetConn(conn *websocket.Conn)
}

type catapultWebsocketMessagePublisher struct {
	m    sync.Mutex
	conn *websocket.Conn
}

func newMessagePublisher(conn *websocket.Conn) MessagePublisher {
	return &catapultWebsocketMessagePublisher{
		m:    sync.Mutex{},
		conn: conn,
	}
}

func (p *catapultWebsocketMessagePublisher) PublishSubscribeMessage(uid string, path string) error {
	dto := &subscribeDTO{
		Uid:       uid,
		Subscribe: path,
	}

	p.m.Lock()
	err := p.conn.WriteJSON(dto)
	p.m.Unlock()
	if err != nil {
		return errors.Wrap(err, "publishing subscribe message into websocket connection")
	}

	return nil
}

func (p *catapultWebsocketMessagePublisher) PublishUnsubscribeMessage(uid string, path string) error {
	dto := &unsubscribeDTO{
		Uid:         uid,
		Unsubscribe: path,
	}

	p.m.Lock()
	err := p.conn.WriteJSON(dto)
	p.m.Unlock()
	if err != nil {
		return errors.Wrap(err, "publishing unsubscribe message into websocket connection")
	}

	return nil
}

func (p *catapultWebsocketMessagePublisher) SetConn(conn *websocket.Conn) {
	p.m.Lock()
	p.conn = conn
	p.m.Unlock()
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
