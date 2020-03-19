package routers

import (
	"classsOne/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	//"github.com/astaxie/beego/plugins/cors"
)

func init() {
		//beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//	//AllowAllOrigins:  true,
		//	AllowOrigins:      []string{"https://192.168.1.104"},
		//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//	AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		//	ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		//	AllowCredentials: true,
		//}))

    //beego.Router("/", &controllers.MainController{})
    beego.InsertFilter("/Article/*",beego.BeforeRouter,FilterFunc)//插入一个路由过滤器，放在router之前
    beego.Router("/register",&controllers.RegController{},"get:ShowReg;post:HandleReg")
	beego.Router("/",&controllers.LoginController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/Article/ShowArticle",&controllers.ArticleController{},"get:ShowArticleList;post:HandleSelect")
	beego.Router("/Article/AddArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/Article/ArticleContent",&controllers.ArticleController{},"get:ShowArticleContent")
	beego.Router("/Article/DeleteArticle",&controllers.ArticleController{},"get:HandleDelete")
	beego.Router("/Article/UpdateArticle",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
	//添加类型
	beego.Router("/Article/AddArticleType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
	//退出登陆

	beego.Router("/Article/Logout",&controllers.ArticleController{},"get:Logout")
}

var FilterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302,"/")
	}
}