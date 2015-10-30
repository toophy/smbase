package app

import (
	"errors"
	"fmt"
	"io"
	"net"
)

// 消息包定义
// 包长度   : uint16  :  (本包所有数据长度)
// 包验证码 : uint8   :  本包验证码
// 包消息数 : uint8   :  本包容纳的消息数量
// 包内数据 : 字节流  :  字节流消息
//
// 消息定义
// 消息长度 : uint16  :  一条消息的长度
// 消息编号 : uint16  :  消息编号(用于选择消息处理函数)
// 消息内容 : 字节流  :  字节流消息内容

// 紧缩方案
// 1. 包长度, 验证码, 消息数,
// 1.1  包长度 => 2^20 , 验证码 => 2^5 , 消息数 => 2^7
//
// 2. 消息长度, 消息编号
// 2.1 消息长度 => 2^20 , 消息编号 => 2^12

const (
	MaxDataLen     = 5080
	MaxSendDataLen = 4000
	MaxHeader      = 2
)

// 网络消息体
type Ty_net_msg struct {
	Data []byte
	Len  int
}

func (t *Ty_net_msg) InitNetMsg() {
	t.Len = 0
	t.Data = make([]byte, MaxDataLen)
}

func (t *Ty_net_msg) PrintData() {
	fmt.Println(t.Data[:t.Len+2])
}

func (t *Ty_net_msg) ReadData(conn *net.TCPConn) error {

	t.Len = 0
	length, err := io.ReadFull(conn, t.Data[:2])
	if length != MaxHeader {
		GetApp().LogWarn("Packet header : %d != %d", length, MaxHeader)
		return err
	}
	if err != nil {
		return err
	}

	body_len := int(t.Data[1]) + (int(t.Data[0]) << 8)

	if body_len > (MaxDataLen - 2) {
		err = errors.New("Body too much")
		return err
	}

	t.Len = body_len + 2
	return t.ReadBody(conn)
}

func (t *Ty_net_msg) ReadBody(conn *net.TCPConn) error {

	length, err := io.ReadFull(conn, t.Data[2:t.Len])
	if length != (t.Len - 2) {
		GetApp().LogWarn("Packet length : %d != %d ", length, t.Len-2)
		return err
	}
	if err != nil {
		return err
	}
	// 注意 : 可以解密

	return nil
}

func (t *Ty_net_msg) Send(conn *net.TCPConn) error {
	if t.Len > MaxHeader && t.Len < MaxSendDataLen {

		t.Data[0] = byte((t.Len & 0xFF00) >> 8)
		t.Data[1] = byte(t.Len & 0xFF)

		_, err := conn.Write(t.Data[:MaxHeader+t.Len])
		if err != nil {
			GetApp().LogWarn(err.Error())
			return err
		}
	}

	return nil
}
