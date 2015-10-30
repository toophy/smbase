package app

import (
	"net"
)

// 客户端连接
type ClientConn struct {
	Name    string
	Type    string
	Address string
	Id      int
	Sid     int
	Msg     Ty_net_msg
	Stream  Ty_msg_stream
	Conn    *net.TCPConn
}

func (c *ClientConn) InitClient(id int, con *net.TCPConn) {
	c.Type = "Null"
	c.Address = con.RemoteAddr().String()
	c.Id = id
	c.Conn = con
	c.Msg.InitNetMsg()
	c.Stream.InitMsgStream(&c.Msg)
}

func (c *ClientConn) IsNull() bool {
	return c.Type == "Null"
}

// 侦听服务
type ListenConn struct {
	Name    string
	Type    string
	Address string
	Conn    *net.TCPListener
}

func (c *ListenConn) InitListen(name, net_type, address string, con *net.TCPListener) {
	c.Name = name
	c.Type = net_type
	c.Address = address
	c.Conn = con
}
