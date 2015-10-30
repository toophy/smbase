package controllers

import (
	"github.com/astaxie/beego"
	"unsafe"
)

type LocalController struct {
	beego.Controller
}

// 平台账号登录-返回成功,失败,登录地址
func (this *LocalController) Get() {
	CloseConnect("fuck you", unsafe.Pointer(this))
	return
}

// 平台账号登录-返回成功,失败,登录地址
func (this *LocalController) Post() {
	CloseConnect("fuck you", unsafe.Pointer(this))
	return
}

func (this *LocalController) Reload() {

	local_call := true
	if this.Ctx.Input.IP() != "127.0.0.1" ||
		this.Ctx.Input.Host() != "127.0.0.1" ||
		this.Ctx.Input.Domain() != "127.0.0.1" {
		local_call = false
	}

	if local_call {
		CloseConnect("ok", unsafe.Pointer(this))
		return
	}

	CloseConnect("failed", unsafe.Pointer(this))
	return
}
