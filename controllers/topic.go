package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
	"path"
	"strings"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController)Get()  {
	this.Data["IsTopic"] = true
	this.TplName ="topic.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)

	topics,err := models.GetAllTopics("","",false)
	if err != nil {
		panic(err)
	}
	this.Data["Topic"]=topics
}

func (this *TopicController)Post() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	//	表单解析
	tid := this.Input().Get("tid")
	title := this.Input().Get("title")
	content := this.Input().Get("content")
	category := this.Input().Get("category")
	lable := this.Input().Get("lable")

	//获取附件
	_,fh,err := this.GetFile("attachment")
	if err != nil {
		panic(err)
	}
	var attachment string
	if fh != nil {
		attachment = fh.Filename
		beego.Info(attachment)
		err = this.SaveToFile("attament", path.Join("attachment", attachment))
		if err != nil {
			panic(err)
		}
		if len(tid) == 0 {
			err = models.AddTopic(title, category, lable, content, attachment)
		} else {
			err = models.ModiyTopic(tid, title, category, lable, content, attachment)
		}
		if err != nil {
			panic(err)
		}
		this.Redirect("/topic", 302)
	}
}

func (this *TopicController) Add() {
	if  !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	this.TplName = "topic_add.html"
	this.Data["IsLogin"] = true
}

func (this *TopicController) Delete() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	err := models.DeleteTopic(this.Input().Get("tid"))
	if err != nil {
		panic(err)
	}
	this.Redirect("/topic", 302)
}

func (this *TopicController) Modify() {
	if  !checkAccount(this.Ctx){
		this.Redirect("/login",302)
		return
	}
	this.TplName = "topic_modify.html"
	tid := this.Input().Get("tid")
	topic,err := models.GetTopic(tid)
	if err != nil {
		this.Redirect("/", 302)
		return
	}
	this.Data["Topic"] = topic
	this.Data["Tid"] = tid
	this.Data["IsLogin"] = true
}

func (this *TopicController) View() {
	this.TplName = "topic_view.html"

	reqUrl := this.Ctx.Request.RequestURI
	i := strings.LastIndex(reqUrl, "/")
	tid := reqUrl[i+1:]
	topic, err := models.GetTopic(tid)
	if err != nil {
		this.Redirect("/", 302)
		return
	}
	this.Data["Topic"] = topic
	this.Data["Lables"] = strings.Split(topic.Lables, " ")

	replies, err := models.GetAllReplies(tid)
	if err != nil {
		return
	}

	this.Data["Replies"] = replies
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}















