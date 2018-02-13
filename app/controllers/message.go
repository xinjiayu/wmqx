package controllers

import (
	"rmqc/container"
	"github.com/valyala/fasthttp"
	"rmqc/message"
	"rmqc/app/service"
	"rmqc/app"
)

type MessageController struct {
	BaseController
}

// return MessageController
func NewMessageController() *MessageController {
	return &MessageController{}
}

// add a message
func (this *MessageController) Add(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	name := this.GetCtxString(ctx, "name")
	comment := this.GetCtxString(ctx, "comment")
	durable := this.GetCtxBool(ctx, "durable")
	isNeedToken := this.GetCtxBool(ctx, "is_need_token")
	mode := this.GetCtxString(ctx, "mode")
	token := this.GetCtxString(ctx, "token")

	if name == "" || comment == ""{
		this.jsonError(ctx, "param require!", nil)
		return
	}
	if (mode != "fanout") && (mode != "topic") && (mode != "direct") {
		this.jsonError(ctx, "param error!", nil)
		return
	}

	// check message is exists
	ok := container.Ctx.QMessage.IsExistsMessage(name)
	if ok == true {
		this.jsonError(ctx, "message "+name+" is exist", nil)
		return
	}

	msg := &message.Message{
		Consumers  : []*message.Consumer{},
		Durable     : durable,
		IsNeedToken : isNeedToken,
		Mode        : mode,
		Name        : name,
		Token       : token,
		Comment     : comment,
	}

	err := service.NewMQ().DeclareExchange(name, mode, durable)
	if err != nil {
		app.Log.Error("Add message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "add message failed: "+err.Error(), nil)
		return
	}

	err = container.Ctx.QMessage.AddMessage(msg)
	if err != nil {
		app.Log.Error("Add message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "add message failed"+err.Error(), nil)
		return
	}

	app.Log.Info("Add message "+name+" success!")
	this.jsonSuccess(ctx, "success", nil)
}

// update a message
func (this *MessageController) Update(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	name := this.GetCtxString(ctx, "name")
	comment := this.GetCtxString(ctx, "comment")
	durable := this.GetCtxBool(ctx, "durable")
	isNeedToken := this.GetCtxBool(ctx, "is_need_token")
	mode := this.GetCtxString(ctx, "mode")
	token := this.GetCtxString(ctx, "token")

	if name == "" || comment == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}
	if (mode != "fanout") && (mode != "topic") && (mode != "direct") {
		this.jsonError(ctx, "param error!", nil)
		return
	}

	// check message is exists
	ok := container.Ctx.QMessage.IsExistsMessage(name)
	if ok == false {
		this.jsonError(ctx, "message "+name+" not exist", nil)
		return
	}

	msg := &message.Message{
		Durable     : durable,
		IsNeedToken : isNeedToken,
		Mode        : mode,
		Name        : name,
		Token       : token,
		Comment     : comment,
	}

	err := service.NewMQ().DeclareExchange(name, mode, durable)
	if err != nil {
		app.Log.Error("Update message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "update message failed: "+err.Error(), nil)
		return
	}

	err = container.Ctx.QMessage.UpdateMessageByName(name, msg)
	if err != nil {
		app.Log.Error("Update message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "update message failed: "+err.Error(), nil)
		return
	}

	app.Log.Info("Update message "+name+" success!")

	this.jsonSuccess(ctx, "success", nil)
}

// delete a message
func (this *MessageController) Delete(ctx *fasthttp.RequestCtx) {
	r := this.AccessToken(ctx)
	if r != true {
		this.jsonError(ctx, "token error", nil)
		return
	}

	name := this.GetCtxString(ctx, "name")
	if name == "" {
		this.jsonError(ctx, "param require!", nil)
		return
	}

	// check message is exists
	ok := container.Ctx.QMessage.IsExistsMessage(name)
	if ok == false {
		this.jsonError(ctx, "message "+name+" not exist", nil)
		return
	}

	err := service.NewMQ().DeleteExchange(name)
	if err != nil {
		app.Log.Error("Delete message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "delete message failed: "+err.Error(), nil)
		return
	}

	err = container.Ctx.QMessage.DeleteMessageByName(name)
	if err != nil {
		app.Log.Error("Delete message "+name+" failed: "+err.Error())
		this.jsonError(ctx, "delete message failed: "+err.Error(), nil)
		return
	}

	app.Log.Info("Delete message "+name+" success!")

	this.jsonSuccess(ctx, "success", nil)
}

// get message status
func (messageController *MessageController) Status(ctx *fasthttp.RequestCtx) {

}