package app

import (
	"net"
)

func (this *AppBase) AddConn(c *ClientConn) {
	this.Conns[this.ConnLast] = c
	this.ConnLast++
}

func (this *AppBase) DelConn(id int) {
	if _, ok := this.Conns[id]; ok {
		if len(this.Conns[id].Name) > 0 {
			delete(this.RemoteSvr, this.Conns[id].Name)
		}
		delete(this.Conns, id)
	}
}

func (this *AppBase) GetConnById(id int) *ClientConn {
	if v, ok := this.Conns[id]; ok {
		return v
	}
	return nil
}

func (this *AppBase) GetConnByName(name string) *ClientConn {
	if v, ok := this.RemoteSvr[name]; ok {
		return v
	}
	return nil
}

func (this *AppBase) RegMsgFunc(id int, f MsgFunc) {
	this.MsgProc[id] = f

	if id > this.MsgProcCount {
		this.MsgProcCount = id
	}
}

func (this *AppBase) Listen(name, net_type, address string, onRet ConnRetFunc) {
	if len(address) == 0 || len(address) == 0 || len(net_type) == 0 {
		onRet("listen failed", name, 0, "listen failed")
		return
	}

	// 打开本地TCP侦听
	serverAddr, err := net.ResolveTCPAddr(net_type, address)

	if err != nil {
		onRet("listen failed", name, 0, "Listen Start : port failed: '"+address+"' "+err.Error())
		return
	}

	listener, err := net.ListenTCP(net_type, serverAddr)
	if err != nil {
		onRet("listen failed", name, 0, "TcpSerer ListenTCP: "+err.Error())
		return
	}

	ln := new(ListenConn)
	ln.InitListen(name, net_type, address, listener)
	this.Listener[name] = ln

	onRet("listen ok", name, 0, "")

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			if !onRet("accept failed", name, 0, "TcpSerer Accept: "+err.Error()) {
				break
			}
			continue
		}
		c := new(ClientConn)
		c.InitClient(this.ConnLast, conn)
		this.AddConn(c)
		onRet("accept ok", "", c.Id, "")

		go this.ConnProc(c, onRet)
	}
}

func (this *AppBase) Connect(name, net_type, address string, onRet ConnRetFunc) {
	if len(address) == 0 || len(net_type) == 0 || len(name) == 0 {
		onRet("connect failed", name, 0, "listen failed")
		return
	}

	// 打开本地TCP侦听
	remoteAddr, err := net.ResolveTCPAddr(net_type, address)

	if err != nil {
		onRet("connect failed", name, 0, "Connect Start : port failed: '"+address+"' "+err.Error())
		return
	}

	conn, err := net.DialTCP(net_type, nil, remoteAddr)
	if err != nil {
		onRet("connect failed", name, 0, "Connect dialtcp failed: '"+address+"' "+err.Error())
	} else {
		c := new(ClientConn)
		c.InitClient(this.ConnLast, conn)
		c.Name = name
		this.RemoteSvr[name] = c
		this.AddConn(c)

		onRet("connect ok", name, c.Id, "")
		go this.ConnProc(c, onRet)
	}
}

func (this *AppBase) ConnProc(c *ClientConn, onRet ConnRetFunc) {

	for {
		c.Stream.Seek(0)
		err := c.Msg.ReadData(c.Conn)

		if err == nil {

			c.Stream.Seek(MaxHeader)
			msg_code := c.Stream.ReadU2()

			if msg_code >= 0 && msg_code <= this.MsgProcCount && this.MsgProc[msg_code] != nil {
				this.MsgProc[msg_code](c)
			}

		} else {
			onRet("read failed", c.Name, c.Id, err.Error())
			break
		}
	}

	onRet("pre close", c.Name, c.Id, "")

	err := c.Conn.Close()
	if err != nil {
		onRet("close failed", c.Name, c.Id, err.Error())
	} else {
		onRet("close ok", c.Name, c.Id, "")
	}

	GetApp().DelConn(c.Id)
}
