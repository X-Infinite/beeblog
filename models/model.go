package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Category struct {
	Id              int64
	Title           string
	CreateTime      time.Time `orm:"index"`
	Views           int64       `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}
type User struct {
	Id int
	Username string
	Password string
}
type Topic struct {
	Id               int64
	Uid              int64
	Title            string
	Category        string
	Lables          string
	Content          string `orm:"size(5000)"`
	Attachment       string
	CreateTime       time.Time `orm:"index"`
	UpdateTime       time.Time `orm:"index"`
	Views            int64      `orm:"index"`
	Author           string
	ReplyTime        time.Time `orm:"index"`
	ReplyCount      int64
	ReplayLastUserId int64
}

type Comment struct {
	Id int64
	Tid int64
	Name string
	Content string   `orm:"size(1000)"`
	Created time.Time`orm:"index"`
}

func RegisterDB() {
	// register model
	orm.RegisterModel(new(Category), new(Topic),new(User),new(Comment))
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	// set default database
	err :=	orm.RegisterDataBase(
		"default",
		"mysql",
		"root:123456@tcp(47.114.0.30:3306)/test?charset=utf8",
		30)
	if err != nil {
		panic(err)
	}
	// create table
	//orm.RunSyncdb("default", false, true)
}

func AddCategory(name string) error {
	o := orm.NewOrm()
	cate := &Category{
		Title:      name,
		CreateTime: time.Now(),
		TopicTime:  time.Now(),
	}
	//查询数据
	qs := o.QueryTable("category")
	err := qs.Filter("titlie", name).One(cate)
	if err == nil {
		return err
	}
	//插入数据
	_, err = o.Insert(cate)
	if err != nil {
    return err
	}
	return nil
}

func DeleteCategory(id string) error {
	cid , err := strconv.ParseInt(id,10 ,64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()

	cate := &Category{Id:cid}
	_,err = o.Delete(cate)
	return err
}

func GetAllCategories()([]*Category, error) {
	o :=orm.NewOrm()

	cates :=make([]*Category,0)

	qs := o.QueryTable("category")

	_, err := qs.All(&cates)
	return cates , err
}


//文章的增删改查

func GetAllTopics(category ,lable string,isDesc bool) (topics []*Topic,err error){

	 o:= orm.NewOrm()
	 topics =make([]*Topic,0)
	 qs := o.QueryTable("topic")
	if  isDesc {
		if len(category) >0 {
			qs.Filter("category",category)
		}
		if len(lable) >0 {
			qs =qs.Filter("lables__contains","$"+lable+"#")
		}
		//_,err =qs.OrderBy("-created").All(&topics)
		
	}else {
		_,err = qs.All(&topics)
	}
	return  topics,err
}

func GetTopic(tid string) (*Topic,error)  {
	tidNum,err := strconv.ParseInt(tid,10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(Topic)
	qs :=o.QueryTable("topic")
	err = qs.Filter("id",tidNum).One(topic)
	if err != nil {
		return nil, err
	}
	topic.Views++
	_,err =o.Update(topic)
	topic.Lables =strings.Replace(strings.Replace(
		topic.Lables,"#"," ",-1),"$"," ",-1)
	return topic, nil
}

func AddTopic(title, category, lable, context, attactment string) error {
	lable ="$"+strings.Join(strings.Split(lable,""),"#$")+"#"
	o := orm.NewOrm()
	topic :=&Topic{
		Title: title,
		Category: category,
		Lables: lable,
		Content: context,
		Attachment: attactment,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	_,err := o.Insert(topic)
	if err != nil {
		return err
	}

	//更新分类统计
	cate := new(Category)
	qs :=o.QueryTable("category")
	err = qs.Filter("title",category).One(cate)
	if err == nil {
		cate.TopicCount++
		_,err=o.Update(cate)
	}
	return err
}

func ModiyTopic(tid, title, category, lable, content, attachment string) error {

	tidNum,err := strconv.ParseInt(tid,10, 64)
	if err != nil {
		return err
	}
	
	lable = "$"+strings.Join(strings.Split(lable,""),"#$")+"#"
	var oldCate,oldAttach string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate =topic.Category
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Category =category
		topic.Lables = lable
		topic.Content =content
		topic.Attachment =attachment
		topic.UpdateTime =time.Now()
		_,err = o.Update(topic)
		if err != nil {
			return err
		}
		
	}
	
	//更新分类统计
	if len(oldCate) >0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title",oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_,err = o.Update(cate)
		}
	}
	//删除旧的附件
	if len(oldAttach)>0 {
		os.Remove(path.Join("attachment",oldAttach))
	}
	cate := new(Category)
	qs :=o.QueryTable("category")
	err =qs.Filter("title",category).One(cate)
	if err == nil {
		cate.TopicCount++
		_,err =o.Update(cate)
	}
	
	return nil
}

func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	var oldCate string
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}

	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}
	return err
}

//回复

func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	o := orm.NewOrm()
	_, err = o.Insert(reply)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		_, err = o.Update(topic)
	}
	return err
}

func GetAllReplies(tid string) (replies []*Comment, err error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	replies = make([]*Comment, 0)

	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&replies)
	return replies, err
}

func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	var tidNum int64
	reply := &Comment{Id: ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
		_, err = o.Delete(reply)
		if err != nil {
			return err
		}
	}

	replies := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = replies[0].Created
		topic.ReplyCount = int64(len(replies))
		_, err = o.Update(topic)
	}
	return err
}
