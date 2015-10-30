package controllers

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"io"
	"io/ioutil"
	"net/http"
	"unsafe"
)

func ReverseProxyYouxi(url string, ctx *context.Context) error {

	var err error
	var resp *http.Response
	if ctx.Request.Method == "POST" {
		if ctx.Request.ParseForm() != nil {
			return errors.New("ctx request parse form failed")
		}
		resp, err = http.PostForm(url, ctx.Request.Form)
	} else {
		resp, err = http.Get(url)
	}
	if err != nil {
		return err
	}
	contentType := resp.Header.Get("Content-Type")
	if len(contentType) > 0 {
		ctx.ResponseWriter.Header().Set("Content-Type", contentType)
	}

	io.Copy(ctx.ResponseWriter, resp.Body)

	return nil
}

// 反向代理
func ReverseProxy(new_host string, ctx *context.Context) error {

	ctx.Request.URL.Host = new_host
	ctx.Request.URL.Scheme = "http"

	//ctx.Request.Host = new_host
	ctx.Request.RequestURI = ""

	resp, err := http.DefaultClient.Do(ctx.Request)

	if err != nil {
		println("proxy err 1 : " + err.Error())
		return errors.New("proxy err 1 : " + err.Error())
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		for _, vv := range v {
			ctx.ResponseWriter.Header().Add(k, vv)
		}
	}
	for _, c := range resp.Cookies() {
		ctx.ResponseWriter.Header().Add("Set-Cookie", c.Raw)
	}
	ctx.ResponseWriter.WriteHeader(resp.StatusCode)
	result, err_body := ioutil.ReadAll(resp.Body)
	if err_body != nil && err_body != io.EOF {
		println("proxy err 2 : " + err_body.Error())
	}

	ctx.ResponseWriter.Write(result)

	return nil
}

// 主动断开连接
func CloseConnect(ret string, ctrl unsafe.Pointer) {
	this := (*beego.Controller)(ctrl)
	if len(ret) > 0 {
		this.Ctx.WriteString(ret)
	} else {
		this.Ctx.WriteString("\n")
	}
	this.StopRun()
	return
}

// 主动断开连接
func CloseConnectNoData(ctrl unsafe.Pointer) {
	this := (*beego.Controller)(ctrl)
	this.Ctx.WriteString("\n")
	this.StopRun()
	return
}
