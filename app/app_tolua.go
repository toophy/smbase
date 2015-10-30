package app

import (
	lua "github.com/toophy/gopher-lua"
)

// 获取用Lua类型封装结构指针  *LUserData
func (this *AppBase) GetLUserData(n string, a interface{}) *lua.LUserData {

	ud := this.luaState.NewUserData()
	ud.Value = a
	this.luaState.SetMetatable(ud, this.luaState.GetTypeMetatable(n))

	return ud
}

// 调用Lua函数 : 没有参数, 没有返回值
func (this *AppBase) Tolua_Common(m string, f string) {
	// 捕捉异常
	defer func() {
		if r := recover(); r != nil {
			this.LogFatal("AppBase:Tolua_Common (" + m + "," + f + ") : " + r.(error).Error())
		}
	}()

	// 调用Lua脚本函数
	if err := this.luaState.CallByParam(lua.P{
		Fn:      this.luaState.GetFunction(m, f), // 调用的Lua函数
		NRet:    0,                               // 返回值的数量
		Protect: true,                            // 保护?
	}); err != nil {
		panic(err)
	}

	return
}

// 调用Lua函数 : 有参数, 有返回值
func (this *AppBase) Tolua_Common_Param_Ret(m string, f string, t lua.LValue) (ret lua.LValue) {
	// 捕捉异常
	defer func() {
		if r := recover(); r != nil {
			ret = nil
			this.LogFatal("AppBase:Tolua_Common_Param_Ret (" + m + "," + f + ") : " + r.(error).Error())
		}
	}()

	if t == nil {
		t = &this.luaNilTable
	}

	// 调用Lua脚本函数
	if err := this.luaState.CallByParam(lua.P{
		Fn:      this.luaState.GetFunction(m, f), // 调用的Lua函数
		NRet:    1,                               // 返回值的数量
		Protect: true,                            // 保护?
	}, t); err != nil {
		panic(err)
	}

	// 处理Lua脚本函数返回值
	ret = this.luaState.Get(-1)
	this.luaState.Pop(1)
	return
}

// 调用Lua函数 : 只有参数
func (this *AppBase) Tolua_Common_Param(m string, f string, t lua.LValue) {
	// 捕捉异常
	defer func() {
		if r := recover(); r != nil {
			this.LogFatal("AppBase:Tolua_Common_Param (" + m + "," + f + ") : " + r.(error).Error())
		}
	}()

	if t == nil {
		t = &this.luaNilTable
	}

	// 调用Lua脚本函数
	if err := this.luaState.CallByParam(lua.P{
		Fn:      this.luaState.GetFunction(m, f), // 调用的Lua函数
		NRet:    0,                               // 返回值的数量
		Protect: true,                            // 保护?
	}, t); err != nil {
		panic(err)
	}

	return
}
