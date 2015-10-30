package routers

import (
	"github.com/astaxie/beego"
	"github.com/toophy/##[AppName]##/controllers"
)

func init() {
	beego.AutoRouter(&controllers.LocalController{})
}
