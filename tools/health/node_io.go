package health

import (
	"log"
	"net"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/tools/health/packets"
)

type NodeTcpIo struct {
	nodeInfo *NodeInfo
	conn     net.Conn
}

func NewNodeTcpIo(info *NodeInfo) (*NodeTcpIo, error) {
	log.Printf("Dialing %s=%v", info.Endpoint, info.IdentityKey)
	connection, err := net.DialTimeout("tcp", info.Endpoint, 5*time.Second)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	log.Printf("Connected to %s=%v", info.Endpoint, info.IdentityKey)

	return &NodeTcpIo{
		nodeInfo: info,
		conn:     connection,
	}, nil
}

func (nc *NodeTcpIo) Write(p packets.Byter) (int, error) {
	return nc.conn.Write(p.Bytes())
}

func (nc *NodeTcpIo) Read(parser packets.Parser, expectedSize int) error {
	offset, n := 0, 0
	var err error
	buf := make([]byte, expectedSize)
	for offset < len(buf) {
		n, err = nc.conn.Read(buf[offset:])
		if err != nil {
			return err
		}
		offset += n
	}

	return parser.Parse(buf)
}

func (nc *NodeTcpIo) Close() error {
	return nc.conn.Close()
}
