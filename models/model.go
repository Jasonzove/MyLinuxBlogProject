package models

import ("github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)
//表的设计
type User struct {
	Id int
	UserName string
	Passwd string
	Articles[]*Article `orm:"rel(m2m)"` //阅读的文章 用户表和文章表是多对多的关系 一个用户可以看多篇文章
}
//文章表和文章类型表是一对多
type Article struct {
	Id int			  `orm:"pk;auto"`
	Title string      `orm:"size(20)"`   //文章标题
	Content string    `orm:"size(500)"`   //内容
	Img string        `orm:"size(50);null"`   //图片 路径
	Time time.Time    `orm:"type(datetime);auto_now_add"`  //发布时间
	Count int         `orm:"default(0)"` //阅读量
	ArticleType *ArticleType `orm:"rel(fk)"`//设置外键  rel reverse是成对出现的 //一对多是加外键，多对多是见了一个关系表
											// 文章类型和文章是一对多的关系，这里是文章就是一
											//这里其实是一个article的一个外键，就是articletype和article之间的一个关系字段
	Users[]*User `orm:"reverse(many)"` //读者 一片文章可以被多人看
}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`
	Articles[]*Article `orm:"reverse(many)"`//文章类型和文章是一对多的关系，这里是文章是多
}


func init()  {
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/newsWeb?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}