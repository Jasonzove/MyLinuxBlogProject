package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"classsOne/models"
	"math"
	"strconv"

)

type ArticleController struct{
	beego.Controller
}
//处理下拉旷改变
func (this*ArticleController)HandleSelect()  {
	typename := this.GetString("select")
	//beego.Info(typename)
	//处理数据
	if typename == ""{
		beego.Info("下拉框传递数据失败")
		return 
	}
	//查询数据
	o := orm.NewOrm()
	var articles []models.Article
	o.QueryTable("article").RelatedSel("ArticleType").Filter("ArticleType__TypeName",typename).All(&articles)
	//beego.Info(articles)

}

//文章列表页
func (this*ArticleController)ShowArticleList()  {
	beego.Info("Start ShowArticleList!")
	//在router之前加过session过滤函数了，此处就不再需要session过滤
	//userName := this.GetSession("userName")
	////GetSession的返回值是interface，interface判断应该用nil判断
	//if userName == nil {
	//	this.Redirect("/", 302)//如果没有session 回到登录界面
	//	return
	//}
	//1查询
	o:=orm.NewOrm()
	qs := o.QueryTable("Article")//指定查询的表 也可以直接使用对象作为表名，例如var article Article o.QueryTable(article)
	var articles []models.Article
	//_,err := qs.All(&articles)//select&from article
	//if err!=nil {
	//	beego.Info("查询错误")
	//	return
	//}
	//beego.Info(articles[0])

	//当通过ShowArticle回车进入界面的时候，会直接路由到这个函数，进行相应的数据加载，
	// 但是由于没有界面，所有的GetString操作得到的值都是空的

	//根据类型获取数据
	typeName := this.GetString("select")
	var count int64

	//pageIndex1 := 1
	pageIndex := this.GetString("pageIndex")//当前第几页
	beego.Info("pageIndex is " + pageIndex)

	pageIndex1,err := strconv.Atoi(pageIndex) //如果不存在就会有错误，就不会把值付给pageIndex1，就不会覆盖默认值
	if err != nil {
		pageIndex1 = 1
		pageIndex = "1"
	}
	//获取总数据
	//count,err := qs.RelatedSel("ArticleType").Count()//返回数据条目数 qs指定查询的表是Article，就是查询Article表的条目shu
	//if err !=nil {
	//	beego.Info("查询错误")
	//	return
	//}

	//beego.Info("Start ShowArticleList!")
	pageSize := 1//每页显示几条数据
	start := pageSize*(pageIndex1-1)//每页索引的起始位置
	//RelatedSel("ArticleType")在查询的时候关联的表ArticleType有数据的话才能查询出来。没有数据的话原表中article中有的数据也不会被查询出来
	//qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)//查询每页数据并放在articles中//pageSize 一页显示多少 2，start起始位置


	//获取类型数据
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]=types

	if typeName == ""{
		count,_ = qs.Count()
		if count != 0 && len(types ) > 0{
			typeName = types[0].TypeName
		} else {
			beego.Info("types is null or article table is null")
			this.TplName="index.html"
			return
		}
	}

	//limit 查询部分数据 第一个参数为查询的数量，第二个参数为开始查询的位置
	count,_ = qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()

	//获取总页数
	pageCount := float64(count)/float64(pageSize)
	pageCount1 := math.Ceil(pageCount)//向上取整，获取一个和我最近但是比我大的整数

	//首页末页处理
	FirstPage := false
	if pageIndex1 == 1 {
		FirstPage = true
	}
	LastPage := false
	if pageIndex1 == int(pageCount1) {
		LastPage=true
	}



	//beego.Info(typename)
	//处理数据 如果下拉框是空的就意味着没有类型，就要把所有的数据都获取
	if typeName == ""{
		beego.Info("下拉框传递数据失败")
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)
	}else {
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}

	userName := this.GetSession("userName")
	this.Data["userName"] = userName
	this.Data["typeName"] = typeName
	beego.Info("count=",count)
	this.Data["FirstPage"] = FirstPage
	this.Data["LastPage"] = LastPage
	this.Data["count"]=count
	this.Data["pageCount"] = pageCount1
	this.Data["pageIndex"] = pageIndex
	this.Data["articles"]= articles

	//2把数据传递给试图
	this.Layout="layout.html"
	this.LayoutSections=make(map[string]string)
	this.LayoutSections["indexHead"] = "indexHead.html"
	this.TplName="index.html"

}

func (this*ArticleController)ShowArticleContent(){
	//获取数据
	id,err:=this.GetInt("articleId")//获取id号
	//数据校验
	if err != nil{
		beego.Info("传递的链接错误")
	}

	//操作数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id

	o.Read(&article)

	//修改阅读量
	article.Count += 1

	//多对多插入读者
	//1获取操作对象
	//article := models.Article{Id:id}
	//2获取多对多操作对象 获取操作对象中多对多关系的要插入的字段
	m2m := o.QueryM2M(&article,"Users")
	//3获取插入对象
	userName := this.GetSession("userName")

	user := models.User{}
	user.UserName = userName.(string)
	o.Read(&user,"UserName")
	//多对多插入
	_, err = m2m.Add(&user)
	if err != nil {
		beego.Info("插入失败")
		return
	}
	o.Update(&article)

	//o.LoadRelated(&article,"Users")
	//o.QueryTable("Article").RelatedSel("User").Filter("Users__User__UserName",userName.(string)).
		//Distinct().Filter("Id",id).One(&article)
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)
	beego.Info(article)
	//返回视图页面
	this.Data["users"] = users
	this.Data["article"] = article
	this.Layout="layout.html"
	this.LayoutSections=make(map[string]string)
	this.LayoutSections["contentHead"] = "head.html"
	this.TplName = "content.html"

}

func (this*ArticleController)ShowAddArticle(){
	//查询类型数据，传递到视图中
	o:=orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]=types
	this.TplName="add.html"
}

//1那数据
//2判断数据
//3插入数据
//4返回视图
func (this*ArticleController)HandleAddArticle() {
	//1那数据
	artiName:= this.GetString("articleName")
	artiContent := this.GetString("content")
	f,h,err := this.GetFile("uploadname")

	defer f.Close()
	//1判断文件格式
	ext:=path.Ext(h.Filename)//获取文件后缀
	if ext!=".jpg" && ext != ".png" && ext!= ".jpeg"{
		beego.Info("上传格式不正确")
		return
	}
	//2判断文件大小
	if h.Size>5000000 {
		beego.Info("文件太大不允许上传")
		return
	}
	//3不能重名
	fileName := time.Now().Format("2006-01-02 15:04:05")//时间格式的固定字符串

	this.SaveToFile("uploadname","./static/img/"+fileName+ext)
	if err!=nil {
		beego.Info("上传文件失败")
		return
	}
	//3插入数据
	//1获取orm对象
	o:= orm.NewOrm()
	//2创建一个插入对象
	article := models.Article{}
	//3赋值
	article.Title = artiName
	article.Content = artiContent
	article.Img = "./static/img/"+fileName+ext
	article.Time = time.Now()

	//给article对象赋值
	//获取到下拉框传递的类型数据
	typename := this.GetString("select")
	if typename == ""{
		beego.Info("下拉框数据错误")
		return
	}
	//获取type对向
	var artiType models.ArticleType
	artiType.TypeName=typename
	beego.Info("typename 是 " + typename)
	err = o.Read(&artiType,"TypeName")
	if err != nil {
		beego.Info("获取类型错误")
	}
	article.ArticleType = &artiType

	//4插入
	_,err = o.Insert(&article)
	if err != nil {
		beego.Info("插入数据失败"+err.Error() )
		return
	}
	//4返回视图
	this.Redirect("/Article/ShowArticle",302)
}
//1URL传值
//2执行删除操作
//删除文章
func (this*ArticleController)HandleDelete(){
	id,_:=this.GetInt("id")
	//1.orm对象
	o:=orm.NewOrm()
	//2要有删除对象
	article := models.Article{Id:id}
	//3delete
	_,err := o.Delete(&article)
	if err != nil {
		beego.Info("删除数据失败"+err.Error())

	}
	this.Redirect("/Article/ShowArticle",302)

}

func (this*ArticleController)ShowUpdate()  {
	//获取数据
	id,err:=this.GetInt("id")//获取id号,与GetString()一样 //如果用GetString()可以用strconv.Atoi(id)转换成int
	//数据校验
	if err != nil{
		beego.Info("传递的链接错误")
		return
	}
	//查询操作
	o := orm.NewOrm()
	var article models.Article
	article.Id = id

	err = o.Read(&article)//主键查询
	if err != nil {
		beego.Info("查询错误")
		return
	}
	//把数据传递给试图
	this.Data["article"]=article//通过这一步给前端界面设置{{.article }}对象前端可以通过这个对象的字段给相应的内容赋值
	this.TplName ="update.html"
}
func (this*ArticleController)HandleUpdate()  {
	//1拿数据 通过this.Data["article"]=article设置给界面的数据来获取相应的字段值， 包括界面已经修改的article的字段的值
	name:=this.GetString("articleName")
	content:= this.GetString("content")
	id,_ := this.GetInt("id")
	//
	//
	if name=="" || content==""{
		beego.Info("更新数据失败")
		return
	}
	f,h,err := this.GetFile("uploadname")
	var filename string
	var ext string
	if h.Filename == "" {
		beego.Info("No上传文件")
	}else {
		if err!=nil {
			beego.Info("上传文件成功")
			return
		}
		defer f.Close()//如果没有关闭文件，这个文件的文件流就会存在在内存中，造成内存的浪费
		//1,判断大小
		if h.Size>50000 {
			beego.Info("图片太大")
			return
		}
		ext=path.Ext(h.Filename)
		if ext!=".jpg"&&ext!=".png"&&ext!=".jpeg" {
			beego.Info("上传文件类型错误")
			return
		}
		//3.防止文件名重复
		filename=time.Now().Format("2016-01-02-15:04:05")
		this.SaveToFile("uploadname","./static/img/"+filename+ext)
	}

	//更新操作 取数据库中对应的主键id的数据记录
	o:=orm.NewOrm()
	article:=models.Article{Id:id}
	err = o.Read(&article)
	if err != nil {
		beego.Info("要更新的文章不存在")
		return
	}
	//重新赋值
	article.Title=name
	article.Content=content
	article.Img="./static/img/"+filename+ext
	//将对应id的记录更新的值更新到数据库中
	_,err = o.Update(&article)
	if err != nil {
		beego.Info("更新失败")
		return
	}
	// 跳转，重定向界面
	this.Redirect("/Article/ShowArticle",302)
}

func (this*ArticleController)ShowAddType()  {
	//读取类型表，显示数据
	o:=orm.NewOrm()
	var artiType[] models.ArticleType
	_,err := o.QueryTable("ArticleType").All(&artiType)
	if err != nil {
		beego.Info("查询类型错wu" + err.Error())
	}
	this.Data["types"]=artiType

 	this.TplName="addType.html"

}

func (this*ArticleController)HandleAddType(){
	//1获取数据
	typename := this.GetString("typeName")
	//2判断数据
	if typename == "" {
		beego.Info("添加类型数据为空")
		return
	}

	//3执行插入操作
	o:=orm.NewOrm()
	var artiType models.ArticleType
	artiType.TypeName=typename
	_,err := o.Insert(&artiType)
	if err!=nil {
		beego.Info("插入文章类型失败")
		return
	}
	//4展示视图
	this.Redirect("/Article/AddArticleType",302)
}

func (this*ArticleController)Logout()  {
	//1删除登陆状态
	this.DelSession("userName")
	//2跳转登陆页面
	this.Redirect("/",302)

}