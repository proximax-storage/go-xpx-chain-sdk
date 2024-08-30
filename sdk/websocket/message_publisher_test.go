package websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/websocket"
)

func Test_catapultWebsocketMessagePublisher_PublishSubscribeMessage(t *testing.T) {
	type fields struct {
		conn *websocket.Conn
	}
	type args struct {
		uid  string
		path string
	}

	s := httptest.NewServer(http.HandlerFunc(echoError))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				conn: ws,
			},
			args: args{
				uid:  "123123123",
				path: "test-path",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &catapultWebsocketMessagePublisher{
				conn: tt.fields.conn,
			}

			err := p.PublishSubscribeMessage(tt.args.uid, tt.args.path)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func Test_catapultWebsocketMessagePublisher_PublishUnsubscribeMessage(t *testing.T) {
	type fields struct {
		conn *websocket.Conn
	}
	type args struct {
		uid  string
		path string
	}

	s := httptest.NewServer(http.HandlerFunc(echoError))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				conn: ws,
			},
			args: args{
				uid:  "123123123",
				path: "test-path",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &catapultWebsocketMessagePublisher{
				conn: tt.fields.conn,
			}
			err := p.PublishUnsubscribeMessage(tt.args.uid, tt.args.path)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

var upgr = websocket.Upgrader{}

func echoError(w http.ResponseWriter, r *http.Request) {
	c, err := upgr.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	err = c.WriteMessage(1, []byte("asddasdasd"))
	if err != nil {
		panic(err)
	}
}
