package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type CategoryController struct {
	beego.Controller
}

func (this *CategoryController) Get()  {

	op := this.Input().Get("op")

	switch op {
	case "add":
		name :=this.Input().Get("name")
		if len(name) == 0 {
			break
		}

		err := models.AddCategory(name)
		if err != nil {
			return
		}
		this.Redirect("/category",302)
		return
	case "del":
		id := this.Input().Get("id")
		if len(id) == 0 {
			break
		}

		err := models.DeleteCategory(id)
		if err != nil {
			return
		}
		this.Redirect("/category",302)
		return
		
	}

	this.Data["IsCategory"] =true
	this.TplName ="category.html"
	this.Data["IsLogin"] =checkAccount(this.Ctx)

	var err error
	this.Data["categories"],err = models.GetAllCategories()
	if err != nil {
		return
	}


}