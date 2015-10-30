package app

import (
	"github.com/toophy/##[AppName]##/help"
	lua "github.com/toophy/gopher-lua"
)

// 事件 : lua使用的通用事件
type Event_from_lua struct {
	help.Evt_base
	module   string     // lua模块名
	function string     // lua函数名
	param    lua.LValue // 参数(table)
}

// 事件执行
func (this *Event_from_lua) Exec() bool {
	// 当前线程调用-> 执行这个事件
	GetApp().Tolua_Common_Param(this.module, this.function, this.param)
	return true
}
