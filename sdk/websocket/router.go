package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/handlers"
)

func NewRouter(uid string, publisher MessagePublisher, topicHandlers TopicHandlersStorage) Router {
	router := messageRouter{
		uid:               uid,
		topicHandlers:     topicHandlers,
		messageInfoMapper: messageInfoMapperFn(MapMessageInfo),
		messagePublisher:  publisher,
		dataCh:            make(chan []byte, 1024),
	}

	// TODO: Close channel
	go router.run()

	return &router
}

type Router interface {
	RouteMessage([]byte)
	SetUid(string)
	Close()
}

type messageRouter struct {
	uid               string
	messagePublisher  MessagePublisher
	messageInfoMapper MessageInfoMapper
	topicHandlers     TopicHandlersStorage
	dataCh            chan []byte
}

func (r *messageRouter) run() {
	for {
		select {
		case m, ok := <-r.dataCh:
			if !ok {
				return
			}

			messageInfo, err := r.messageInfoMapper.MapMessageInfo(m)
			if err != nil {
				panic(errors.Wrap(err, "getting message info"))
			}

			handler := r.topicHandlers.GetHandler(Path(messageInfo.ChannelName))
			if handler == nil {
				fmt.Println("getting topic handler from topic handlers storage")
				continue
			}

			if ok := handler.Handle(messageInfo.Address, m); !ok {
				if err := r.messagePublisher.PublishUnsubscribeMessage(r.uid, Path(handler.Format(messageInfo))); err != nil {
					fmt.Println(err, "unsubscribing from topic")
					continue
				}
			}
		}
	}
}

func (r *messageRouter) RouteMessage(m []byte) {
	r.dataCh <- m
}

func (r *messageRouter) SetUid(uid string) {
	r.uid = uid
}

func (r *messageRouter) Close() {
	close(r.dataCh)
}

func MapMessageInfo(m []byte) (*sdk.WsMessageInfo, error) {
	var messageInfoDTO sdk.WsMessageInfoDTO
	if err := json.Unmarshal(m, &messageInfoDTO); err != nil {
		return nil, errors.Wrap(err, "unmarshaling message info data")
	}

	return messageInfoDTO.ToStruct()
}

type MessageInfoMapper interface {
	MapMessageInfo([]byte) (*sdk.WsMessageInfo, error)
}

type messageInfoMapperFn func([]byte) (*sdk.WsMessageInfo, error)

func (p messageInfoMapperFn) MapMessageInfo(m []byte) (*sdk.WsMessageInfo, error) {
	return p(m)
}

type TopicHandler struct {
	Topic
	handlers.Handler
}

type TopicHandlersStorage interface {
	HasHandler(path Path) bool
	GetHandler(path Path) *TopicHandler
	SetTopicHandler(path Path, handler *TopicHandler)
}

type topicHandlersMap map[Path]*TopicHandler

type topicHandlers struct {
	sync.RWMutex
	h topicHandlersMap
}

func (h *topicHandlers) HasHandler(path Path) bool {
	h.RLock()
	defer h.RUnlock()
	_, ok := h.h[path]
	return ok
}

func (h *topicHandlers) GetHandler(path Path) *TopicHandler {
	h.RLock()
	defer h.RUnlock()
	val, ok := h.h[path]
	if !ok {
		return nil
	}

	return val
}

func (h *topicHandlers) SetTopicHandler(path Path, handler *TopicHandler) {
	h.Lock()
	defer h.Unlock()
	h.h[path] = handler
}

type Topic interface {
	Format(info *sdk.WsMessageInfo) Path
}

type topicFormatFn func(info *sdk.WsMessageInfo) Path

func (f topicFormatFn) Format(info *sdk.WsMessageInfo) Path {
	return f(info)
}

func formatPlainTopic(info *sdk.WsMessageInfo) Path {
	return Path(fmt.Sprintf("%s/%s", Path(info.ChannelName), info.Address.Address))
}

func formatBlockTopic(_ *sdk.WsMessageInfo) Path {
	return pathBlock
}
