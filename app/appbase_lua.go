package app

import (
	"errors"
	lua "github.com/toophy/gopher-lua"
)

// 初始化LuaState, 可以用来 Reload LuaState
func (this *AppBase) reloadLuaState() error {

	if this.luaState != nil {
		this.luaState.Close()
		this.luaState = nil
	}

	this.luaState = lua.NewState()
	if this.luaState == nil {
		return errors.New("场景线程初始化Lua失败")
	}

	RegLua_all_appBase(this.luaState)

	// 注册公告变量-->本线程
	this.luaState.SetGlobal("ts", this.GetLUserData("AppBase", this))

	// 执行初始化脚本
	this.luaState.Require("data/app_init")

	// 加载所有 ##[AppName]## 文件夹里面的 *.lua 文件
	this.luaState.RequireDir("data/app")

	return nil
}

// !!!只能获取, 不准许保存指针, 获取LState
func (this *AppBase) GetLuaState() *lua.LState {
	return this.luaState
}

// lua投递事件
func (this *AppBase) PostEventFromLua(m string, f string, t uint64, p lua.LValue) bool {
	evt := &Event_from_lua{}
	evt.Init("", t)
	evt.module = m
	evt.function = f
	evt.param = p
	return this.PostEvent(evt)
}
