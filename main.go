package main

import (
	"github.com/astaxie/beego"
	_ "testBee/models"
	_ "testBee/routers" // `_` 作用是执行routers包里面的init函数
)

func main() {
	// 关联视图函数
	beego.AddFuncMap("prePage", ShowPrePage)
	beego.AddFuncMap("nextPage", ShowNextPage)
	beego.AddFuncMap("showSelected", ShowSelectedType)

	beego.Run()
}

/* 视图函数: 处理视图中简单业务逻辑
 1.创建后台函数
 2.在视图中定义函数
	{{.page | prePage}}  `|` 之前的是函数参数, 这种形式只能传一个参数
	{{nextPage .page .pageCount}} 这种形式可以传多个参数
 3.在beego.Run()之前使用addFuncMap关联起来
*/

// 定义一个后台函数
func ShowPrePage(page int) int {
	if page == 1 {
		return page
	}
	return page - 1
}

func ShowNextPage(page, pageCount int) int {
	if page == pageCount {
		return page
	}
	return page + 1
}

func ShowSelectedType(selectedName, optionName string) string {
	if selectedName == optionName {
		return "selected"
	}
	return ""
}