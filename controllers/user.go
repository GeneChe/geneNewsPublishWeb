package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"testBee/models"
)

type UserController struct {
	beego.Controller
}

// 显示注册页
func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

// 处理注册数据
func (this *UserController) HandleRegister() {
	// 1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")

	// 2.校验数据
	if userName == "" || pwd == "" {
		this.Data["errMsg"] = "注册信息不完整, 请重新注册"
		this.TplName = "register.html"
		logs.Info("注册信息不完整, 请重新注册")
		return
	}

	// 3.处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.PassWord = pwd
	o.Insert(&user)

	// 4.返回页面
	// this.Ctx.WriteString("注册成功")
	this.Redirect("/login", 302)	// 向游览器响应redirect, 返回一个新url, 并让其重新发起这个url的请求
	// this.TplName = "login.html" 			// 直接返回页面, 游览器地址栏地址不变化
}

// 显示登陆页
func (this *UserController) ShowLogin() {
	// cookie中是否有用户名
	userName := this.Ctx.GetCookie("userName")
	if userName == "" {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}
	this.TplName = "login.html"
}

// 处理登录数据
func (this *UserController) HandleLogin() {
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	if userName == "" || pwd == "" {
		this.Data["errMsg"] = "登录信息不完整, 请重新登录"
		this.TplName = "login.html"
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName

	err := o.Read(&user, "name")
	if err != nil {
		this.Data["errMsg"] = "用户不存在"
		this.TplName = "login.html"
		return
	}

	if pwd != user.PassWord {
		 this.Data["errMsg"] = "密码错误"
		 this.TplName = "login.html"
		return
	}

	// 获取记住状态
	remember := this.GetString("remember")
	// this.Ctx.SetCookie()
	// 第一个参数: key	string
	// 第二个参数: value	string
	// 第三个参数: 时间, 单位秒.	设为负数时表示立即失效
	// cookie不能存中文, 但可将中文通过base64加密成字符串再存. 取得时候再通过base64解密即可完成存储
	if remember == "on" {
		this.Ctx.SetCookie("userName", userName, 100)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}

	// 登录成功设置session(这是session中的临时session)
	// 注意: 在设置session之前要在 app.conf 中设置`sessionon = true` 不然报错
	this.SetSession("userName", userName)
	this.SetSession("userId", user.Id)

	//this.Ctx.WriteString("登录成功")
	this.Redirect("/article/showArticleList", 302)
}

// 退出登录
func (this *UserController) ShowLogout() {
	// 退出登录时, 删除session
	this.DelSession("userName")
	this.DelSession("userId")
	this.Redirect("/login", 302)
}