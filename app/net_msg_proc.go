package app

import ()

func on_c2g_login(c *ClientConn) {
	if c.Id > 0 {
		name := c.Stream.ReadStr()
		GetApp().LogInfo(name)
	}
}
