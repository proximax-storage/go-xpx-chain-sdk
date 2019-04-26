package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/handlers"
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
}

type messageRouter struct {
	uid               string
	messagePublisher  MessagePublisher
	messageInfoMapper MessageInfoMapper
	topicHandlers     TopicHandlersStorage
}

func (r messageRouter) RouteMessage(m []byte) {
	messageInfo, err := r.messageInfoMapper.MapMessageInfo(m)
	if err != nil {
		panic(errors.Wrap(err, "getting message info"))
	}

	handler := r.topicHandlers.GetHandler(Path(messageInfo.ChannelName))
	if handler == nil {
		panic(errors.Wrap(ErrUnsupportedMessageType, "getting topic handler from topic handlers storage"))
	}

	if ok := handler.Handle(messageInfo.Address, m); !ok {
		if err := r.messagePublisher.PublishUnsubscribeMessage(r.uid, Path(handler.Format(messageInfo))); err != nil {
			panic(errors.Wrap(err, "unsubscribing from topic"))
		}
	}

	return
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

type topicHandlers map[Path]*TopicHandler

func (h topicHandlers) HasHandler(path Path) bool {
	_, ok := h[path]
	return ok
}

func (h topicHandlers) GetHandler(path Path) *TopicHandler {
	val, ok := h[path]
	if !ok {
		return nil
	}

	return val
}

func (h topicHandlers) SetTopicHandler(path Path, handler *TopicHandler) {
	h[path] = handler
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
