package main

import (
	"beeblog/controllers"
	"beeblog/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
)

func init() {
	models.RegisterDB()
}

func main() {
	orm.Debug = true

	//同步数据库
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}

	//路由
	beego.Router("/",&controllers.HomeController{})
	beego.Router("/category",&controllers.CategoryController{})
	beego.Router("/login",&controllers.LoginController{})
	beego.Router("/topic",&controllers.TopicController{})
	beego.Router("/reply",&controllers.ReplyController{})
	beego.Router("/reply/add",&controllers.ReplyController{})
	beego.Router("/reply/delete",&controllers.ReplyController{})
	beego.AutoRouter(&controllers.TopicController{})

	_ = os.Mkdir("attachment", os.ModePerm)
	beego.Router("/attachment/:all",&controllers.AttachmentController{})

	beego.Run()
}

