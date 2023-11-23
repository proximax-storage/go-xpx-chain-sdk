package health

import (
	"net"
)

type NodeTcpIo struct {
	nodeInfo *NodeInfo
	conn     net.Conn
}

func NewNodeTcpIo(info *NodeInfo) (*NodeTcpIo, error) {
	connection, err := net.Dial("tcp", info.Endpoint)
	if err != nil {
		return nil, err
	}

	return &NodeTcpIo{
		nodeInfo: info,
		conn:     connection,
	}, nil
}

func (nc *NodeTcpIo) Write(p Byter) (int, error) {
	return nc.conn.Write(p.Bytes())
}

func (nc *NodeTcpIo) Read(parser Parser, expectedSize int) error {
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
