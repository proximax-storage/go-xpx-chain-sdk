package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"sync"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/handlers"
)

func NewRouter(uid string, publisher MessagePublisher, topicHandlers TopicHandlersStorage) Router {
	return &messageRouter{
		uid:               uid,
		topicHandlers:     topicHandlers,
		messageInfoMapper: messageInfoMapperFn(MapMessageInfo),
		messagePublisher:  publisher,
	}
}

type Router interface {
	RouteMessage([]byte)
	SetUid(string)
}

type messageRouter struct {
	uid               string
	messagePublisher  MessagePublisher
	messageInfoMapper MessageInfoMapper
	topicHandlers     TopicHandlersStorage
}

func (r *messageRouter) RouteMessage(m []byte) {
	messageInfo, err := r.messageInfoMapper.MapMessageInfo(m)
	if err != nil {
		panic(errors.Wrap(err, "getting message info"))
	}

	handler := r.topicHandlers.GetHandler(Path(messageInfo.ChannelName))
	if handler == nil {
		fmt.Println("getting topic handler from topic handlers storage")
		return
	}

	if ok := handler.Handle(messageInfo.Address, m); !ok {
		if err := r.messagePublisher.PublishUnsubscribeMessage(r.uid, Path(handler.Format(messageInfo))); err != nil {
			fmt.Println(err, "unsubscribing from topic")
			return
		}
	}

	return
}

func (r *messageRouter) SetUid(uid string) {
	r.uid = uid
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
	sync.Mutex
	h topicHandlersMap
}

func (h *topicHandlers) HasHandler(path Path) bool {
	h.Lock()
	defer h.Unlock()
	_, ok := h.h[path]
	return ok
}

func (h *topicHandlers) GetHandler(path Path) *TopicHandler {
	h.Lock()
	defer h.Unlock()
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
