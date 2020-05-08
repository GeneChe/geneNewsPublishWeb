package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"testBee/models"
)

type MainController struct {
	beego.Controller
}

// 游览器发出的都是get请求
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["Author"] = "gene"
	c.TplName = "index.tpl"
}

func (c *MainController) Post() {
	c.Data["Author"] = "hannibal"
	c.TplName = "index.tpl"
}

func (c *MainController) ShowGet() {
	// 操作表过程
	// 1. 获取orm对象
	o := orm.NewOrm()

	// 2. 执行操作函数 增删改查
	// 插入操作
	/*
	// 2.1 获取插入对象
	var user models.User
	// 2.2 给插入对象赋值
	user.Name = "gene"
	user.PassWord = "123456"
	// 2.3 使用orm对象插入操作
	count, err := o.Insert(&user)
	if err != nil {
		logs.Info("o.Insert err:", err)
		return
	}
	logs.Info(count)
	*/

	// 查询操作
	/*
	// 查询结果保存到的对象
	var user models.User
	user.Id = 1	// 查询条件

	// 根据id字段来查询
	// err := o.Read(&user, "id")
	err := o.Read(&user)	// 如果查询id字段, 可以不写
	if err != nil {
		logs.Info("o.Read err:", err)
		return
	}
	logs.Info(user)
	 */

	// 更新操作
	/*
	var user models.User
	user.Id = 1
	err := o.Read(&user)
	if err != nil {
		logs.Info("更新数据不存在", err)
		return
	}
	user.Name = "hannibal"
	count, err := o.Update(&user)
	if err != nil {
		logs.Error("更新失败")
		return
	}
	logs.Info(count)
	 */

	// 删除操作
	var user models.User
	// user.Id = 2
	// count, err := o.Delete(&user)

	user.Name = "gene"
	count, err := o.Delete(&user, "name")
	if err != nil {
		logs.Info("删除失败!")
		return
	}
	logs.Info(count)

	// 3. 返回结果
	c.Data["Author"] = "哈哈"
	c.TplName = "index.tpl"
}