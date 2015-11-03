package app

import (
	"fmt"
)

func Main_go() {
	// 重要go协程加入这个玩意, 防止崩溃
	// 出自 : http://www.open-open.com/lib/view/open1422970537748.html
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		var st = func(all bool) string {
	// 			// Reserve 1K buffer at first
	// 			buf := make([]byte, 512)

	// 			for {
	// 				size := runtime.Stack(buf, all)
	// 				// The size of the buffer may be not enough to hold the stacktrace,
	// 				// so double the buffer size
	// 				if size == len(buf) {
	// 					buf = make([]byte, len(buf)<<1)
	// 					continue
	// 				}
	// 				break
	// 			}

	// 			return string(buf)
	// 		}
	// 		lib.Log4e("panic:" + toString(err) + "\nstack:" + st(false))
	// 	}
	// }()

	RegMsgProc()

	//go GetApp().Listen("main_listen", "tcp", ":8001", OnListenRet)
}

func OnListenRet(typ string, name string, id int, info string) bool {
	name_fix := name
	if len(name_fix) == 0 {
		name_fix = fmt.Sprintf("Conn[%d]", id)
	}

	switch typ {
	case "listen failed":
		GetApp().LogFatal("%s : Listen failed[%s]", name_fix, info)

	case "listen ok":
		GetApp().LogInfo("%s : Listen(0.0.0.0:%d) ok.", name_fix, 8001)

	case "accept failed":
		GetApp().LogFatal(info)
		return false

	case "accept ok":
		GetApp().LogDebug("%s : Accept ok", name_fix)

	case "connect failed":
		GetApp().LogError("%s : Connect failed[%s]", name_fix, info)

	case "connect ok":
		GetApp().LogDebug("%s : Connect ok", name_fix)

	case "read failed":
		GetApp().LogError("%s : Connect read[%s]", name_fix, info)

	case "pre close":
		GetApp().LogDebug("%s : Connect pre close", name_fix)

	case "close failed":
		GetApp().LogError("%s : Connect close failed[%s]", name_fix, info)

	case "close ok":
		GetApp().LogDebug("%s : Connect close ok.", name_fix)
	}

	return true
}
