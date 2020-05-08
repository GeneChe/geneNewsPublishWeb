package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"testBee/controllers"
)

func init() {
	// 添加过滤器
	beego.InsertFilter("/article/*", beego.BeforeExec, filter)

	// router
	/*
		// 给请求指定自定义方法, 一个请求指定一个方法
		beego.Router("/login", &controllers.LoginController{}, "get:ShowLogin;post:PostFunc")
		// 给多个请求指定一个方法
		beego.Router("/index", &controllers.IndexController{}, "get,post:HandleFunc")
		// 给所有请求指定一个方法
		beego.Router("/index", &controllers.IndexController{}, "*:HandleFunc")
		// 当两种指定方法冲突时, 越明确的优先级越高
		beego.Router("/index", &controllers.IndexController{}, "*:HandleFunc;post:PostFunc")
	*/
	// beego.Router("/", &controllers.MainController{}) 相同path先注册的优先匹配
	beego.Router("/", &controllers.MainController{}, "get:ShowGet")
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/logout", &controllers.UserController{}, "get:ShowLogout")

	// 文章列表页
	beego.Router("/article/showArticleList", &controllers.ArticleController{}, "get:ShowArticleList")
	// 添加文章页
	beego.Router("/article/addArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
	// 文章详情页
	beego.Router("/article/showArticleDetail", &controllers.ArticleController{}, "get:ShowArticleDetail")
	// 编辑文章页
	beego.Router("/article/updateArticle", &controllers.ArticleController{}, "get:ShowUpdateArticle;post:HandleUpdateArticle")
	// 删除文章
	beego.Router("/article/deleteArticle", &controllers.ArticleController{}, "get:DeleteArticle")
	// 添加文章分类
	beego.Router("/article/addArticleType", &controllers.ArticleController{}, "get:ShowAddType;post:HandleAddType")
}

// 这里不能使用自动推导
var filter = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}