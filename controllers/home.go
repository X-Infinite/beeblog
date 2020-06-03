package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["home"] = true
	this.TplName="home.html"
	this.Data["IsLogin"]= checkAccount(this.Ctx)


	topics,err :=models.GetAllTopics(
		this.Input().Get("cate"),this.Input().Get("lable"),true)
	if err != nil {
		return
	}
	this.Data["Topics"] = topics
	categories,err := models.GetAllCategories()
	if err != nil {
		panic(err)
	}
	this.Data["Categories "] = categories
}

