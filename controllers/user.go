package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"classsOne/models"
	"time"
)
type RegController struct {
	beego.Controller

}
func (this*RegController)ShowReg(){
	this.TplName = "register.html"
}

//1拿到浏览器传递的数据
//2数据处理
//3插入数据库
//4返回视图
func (this*RegController)HandleReg(){
	//1拿到浏览器传递的数据
	name := this.GetString("userName")
	passwd := this.GetString("password")
	//2数据处理
	if name == "" || passwd == ""{
		beego.Info("用户名或者密码不能为空")
		this.TplName = "register.html"
		return
	}
	//3插入数据库
	 //1)获取orm对象
	o:=orm.NewOrm()
	 //2)获取插入对象
	 user := models.User{}
	 //3)插入操作
	 user.UserName=name
	 user.Passwd=passwd
	 _,err:=o.Insert(&user)
	if err != nil {
		beego.Info("插入数据失败")

	}
	 //4返回登陆
	 this.TplName="login.html"
	 //状态码
	 //1xx 继续发送 2XX 请求成功 200 3xx 资源转移 302重定向 4xx 请求错误 404 5xx 服务器错误 500
	 this.Redirect("/",302)
	 //this.Ctx.WriteString("注册成功")



}

type LoginController struct {
	beego.Controller
}

func (this*LoginController)ShowLogin() {
	name := this.Ctx.GetCookie("userName")
	if name != "" {
		this.Data["userName"] = name
		this.Data["check"] = "checked"
	}
	this.TplName="login.html"
}

//1,拿到浏览器数据
//2,数据处理
//3,查找数据库
//4,返回视图
func (this*LoginController)HandleLogin() {
	//1 拿数据
	name := this.GetString("userName")
	passwd := this.GetString("password")
	beego.Info(name,passwd)
	//2,数据处理
	if name == ""||passwd=="" {
		beego.Info("name and password can't be empty!")
		this.TplName="login.html"
	}
	//3,查找数据
	//1)获取orm对象
	o:= orm.NewOrm()
	//2）获取查询对象
	user := models.User{}

	//3）查询
	user.UserName=name
	err := o.Read(&user,"UserName")
	if err!= nil {
		beego.Info("用户名失败")
		this.TplName="login.html"
		return
	}
	//4)判断密码是否一直
	if user.Passwd!=passwd {
		beego.Info("密码失败")
		this.TplName="login.html"
		return

	}
	//实现记住用户名
	check:=this.GetString("remember")//如果前端的checkbox选中了，则返回on否则返回空
	beego.Info(check)

	if check == "on" {
		this.Ctx.SetCookie("userName",name,time.Second*3600)
	}else{
		this.Ctx.SetCookie("userName","ss",-1)//删除cookie 值随便填，others参数填成负数
	}

	//
	this.SetSession("userName",name)

	//4,返回视图

	//this.Ctx.WriteString("登陆成功")
	this.Redirect("/Article/ShowArticle",302)

}