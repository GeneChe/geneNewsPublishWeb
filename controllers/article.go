package controllers

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"testBee/models"
	"time"
)

type ArticleController struct {
	beego.Controller
}

// 展示文章列表页
func (this *ArticleController) ShowArticleList() {
	// 访问页面时, 先获取session, 如没有用户信息则直接返回到登录页
	// 注意: 空session值为 nil 而不是 ""
	userName := this.GetSession("userName")
	if userName == nil {	// 这里不要写成 userName == nil
		this.Redirect("/login",302)
		return
	}

	// 1.获取文章数据
	o := orm.NewOrm()
	qs := o.QueryTable("article")

	var articles []*models.Article
	//count, err := qs.All(&articles) // 满足查询条件的总条数
	//if err != nil {
	//	logs.Error("查询错误", err)
	//}
	//logs.Info("查询到: ", count)

	// 1.1查询总记录数页数
	//co, err := qs.Count()	// 表中的总数据个数. 此处的结果同上count
	//if err != nil {
	//	logs.Error("查询总记录数错误")
	//}

	pageSize := 2
	// math.Ceil()	向上取整
	// math.Floor() 向下取整
	// math.Round() 四舍五入
	// logs.Info(co / 2.0) 2.0与整型运算时, 默认当成整型

	// 1.1 获取查询条件
	page, err := this.GetInt("page", 1) // 没传page时默认page = 1
	if err != nil {
		logs.Error(err)
		page = 1
	}

	// 1.2 根据条件过滤数据 -- 下拉框选中的值
	typeName := this.GetString("select")

	// limit num, offset 返回的qs, 这里的offset < 0时, 当0处理
	// RelatedSel参数是 qs表对应的字段名
	// Filter相当于sql的where.
	// 第一个参数是指定 哪个表的哪个条件字段 `表名__字段名` orm中双下划线特殊用处就在这
	// 第二个参数是条件字段对应的值
	var unionQS orm.QuerySeter
	if typeName == "全部类型" || typeName == "" {
		unionQS = qs.RelatedSel("ArticleType")
	} else {
		unionQS = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName)
	}

	co, _ := unionQS.Count()
	tempCo, _ := unionQS.Limit(pageSize, (page-1) * pageSize).Count()
	// qs 加limit和不加limit返回的count值一样
	logs.Info("--->",co, tempCo)

	pageCount := math.Ceil(float64(co) / float64(pageSize))
	_, _ = unionQS.Limit(pageSize, (page-1) * pageSize).All(&articles)

	// 2.获取所有文章类型
	var types []*models.ArticleType
	_, _ = o.QueryTable("ArticleType").All(&types)
	totalTypes := []*models.ArticleType{{Id:0, TypeName:"全部类型"}}
	totalTypes = append(totalTypes, types...)

	// 3.传递数据
	this.Data["userName"] = userName
	this.Data["selected"] = typeName
	this.Data["types"] = totalTypes
	this.Data["page"] = page
	this.Data["pageCount"] = int(pageCount)
	this.Data["count"] = co
	this.Data["articles"] = articles
	this.Data["title"] = "文章列表"

	// layout设计
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

// 展示添加文章页
func (this *ArticleController) ShowAddArticle() {
	var types []*models.ArticleType
	_, _ = orm.NewOrm().QueryTable("ArticleType").All(&types)

	this.Data["title"] = "添加文章内容"
	this.Data["userName"] = this.GetSession("userName")
	this.Data["types"] = types
	this.Layout = "layout.html"
	this.TplName = "add.html"
}

// 处理添加文章
func (this *ArticleController) HandleAddArticle() {
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	if articleName == "" || content == "" {
		this.Data["errMsg"] = "标题或内容为空, 添加失败"
		this.TplName = "add.html"
		return
	}

	// <form> 中要上传图片时, 需要设置enctype="multipart/form-data"
	// 不然获取时报错 request Content-Type isn't multipart/form-data
	// 处理文件上传
	file, head, err := this.GetFile("uploadname")
	if err != nil {
		this.Data["errMsg"] = "图片上传错误"
		this.TplName = "add.html"
		logs.Error(err)
		return
	}
	defer file.Close()

	// 校验
	// 1.文件大小
	if head.Size > 5000000 {	// size是字节
		this.Data["errMsg"] = "图片过大"
		this.TplName = "add.html"
		return
	}
	// 2.文件类型 path.Ext()用来获取文件后缀
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errMsg"] = "文件类型错误"
		this.TplName = "add.html"
		return
	}
	// 3.防止重名 --	使用时间来当文件名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	// 存储, 没有`/`开头的路径表示: 当前go可执行文件所在目录, 即为项目根目录
	this.SaveToFile("uploadname", "static/img/"+fileName)

	// 处理文章类型
	/*
		// 方式一: 先查出类型对象, 再赋值给文章对象
		typeName := this.GetString("select")
		o := orm.NewOrm()
		var articleType models.ArticleType
		articleType.TypeName = typeName

		_ = o.Read(&articleType, "TypeName")
		article.ArticleType = &articleType
	*/

	typeName := this.GetString("select")
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	_ = o.Read(&articleType, "TypeName")

	/*
		// 方式二: 直接创建类型对象, 赋值一个id, 再赋值给文章对象
		typeId := this.GetString("select")
		var articleType models.ArticleType
		articleType.Id = typeId
		article.ArticleType = &articleType
	*/

	// 处理数据
	var article models.Article
	article.ArtiName = articleName
	article.Acontent = content
	// 在请求资源时, 如果路径没有 `/` 表示当前请求路径所在目录.	--- /article/static/img/
	// 而图片是存在 项目根目录下的 /static 中. 所以这里要加上 `/` 指明实际位置
	article.Aimg = "/static/img/"+fileName
	article.ArticleType = &articleType
	_, err = o.Insert(&article)
	if err != nil {
		logs.Error(err)
	}

	// 返回
	this.Redirect("/article/showArticleList", 302)
}

// 文章详情页
func (this *ArticleController) ShowArticleDetail() {
	articleId, err := this.GetInt("articleId")
	if err != nil {
		logs.Error("获取文章id错误", err)
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId

	// err = o.Read(&article)
	// Filter如果比较的是本表字段, 则不需要指定表名.
	// 查询一个结果使用one()
	err = o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id", articleId).One(&article)
	if err != nil {
		logs.Error("文章不存在", err)
		return
	}

	article.Acount++
	_, _ = o.Update(&article, "acount")

	// 多对多插入游览记录
	// 第一个参数 多对多的一个对象
	// 第二个参数 该对象中的多对多属性
	m2m := o.QueryM2M(&article, "Users")
	if m2m == nil {
		logs.Error("插入游览记录失败")
		return
	}

	userId := this.GetSession("userId")
	if userId == nil {
		this.Redirect("/login", 302)
		return
	}

	var user models.User
	user.Id = userId.(int)
	_, err = m2m.Add(&user)
	if err != nil {
		logs.Error("多对多操作失败")
	}

	// 多对多查询
	// 方式一: o.LoadRelated直接加载数据
	// _, _ = o.LoadRelated(&article, "Users")	// 用户重复
	// 方式二: 在用户表中去重
	var users []*models.User
	_, _ = o.QueryTable("User").Filter("Articles__Article__Id", articleId).Distinct().All(&users)

	this.Data["users"] = users
	this.Data["article"] = article
	this.Data["title"] = "文章详情"
	this.Data["userName"] = this.GetSession("userName")
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

// 显示编辑文章页
func (this *ArticleController) ShowUpdateArticle() {
	articleId, err := this.GetInt("articleId")
	if err != nil {
		logs.Error("更新文章失败, id错误")
		return
	}

	var article models.Article
	article.Id = articleId
	err = orm.NewOrm().Read(&article)
	if err != nil {
		logs.Error("文章%v不存在", articleId)
		return
	}
	logs.Error("%v的文章", articleId)

	this.Data["article"] = article
	this.Data["title"] = "更新文章内容"
	this.Data["userName"] = this.GetSession("userName")
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

// 处理编辑文章页
func (this *ArticleController) HandleUpdateArticle() {
	articleId, err := this.GetInt("articleId")
	if err != nil {
		logs.Error("文章id错误", err)
		return
	}

	articleName := this.GetString("articleName")
	content := this.GetString("content")
	filePath, err := UploadImg(&this.Controller, "uploadname", "static/img")
	logs.Error("上传图片", err)
	if articleName == "" || content == "" {
		logs.Error("标题或内容不能为空")
		return
	}

	var article models.Article
	article.Id = articleId
	article.ArtiName = articleName
	article.Acontent = content
	if filePath != "" {
		article.Aimg = filePath
		_, err = orm.NewOrm().Update(&article, "ArtiName", "Acontent", "Aimg")
		if err != nil {
			logs.Error(err)
		}
	} else {
		_, err = orm.NewOrm().Update(&article, "ArtiName", "Acontent")
		if err != nil {
			logs.Error(err)
		}
	}

	this.Redirect("/article/showArticleList", 302)
}

// 删除文章
func (this *ArticleController) DeleteArticle() {
	articleId, err := this.GetInt("articleId")
	if err != nil {
		logs.Error("获取id错误", err)
		return
	}

	var article models.Article
	article.Id = articleId
	_, err = orm.NewOrm().Delete(&article)
	if err != nil {
		logs.Error("删除操作错误", err)
		return
	}

	this.Redirect("/article/showArticleList", 302)
}

// 展示添加文章分类页
func (this *ArticleController) ShowAddType() {
	var types []*models.ArticleType
	_, _ = orm.NewOrm().QueryTable("ArticleType").All(&types)

	this.Data["types"] = types
	this.Data["title"] = "编辑文章类型"
	this.Data["userName"] = this.GetSession("userName")
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

// 处理添加文章分类
func (this *ArticleController) HandleAddType() {
	typeName := this.GetString("typeName")
	defer this.redirectToTypeList()
	if typeName == "" {
		logs.Error("分类名为空")
		return
	}

	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	err := o.Read(&articleType, "TypeName")
	if err == nil {
		logs.Info("分类名已存在", articleType)
		return
	}

	_, err = o.Insert(&articleType)
	if err != nil {
		logs.Error("插入分类错误", err)
		return
	}
}

func (this *ArticleController)redirectToTypeList() {
	this.Redirect("/article/addArticleType", 302)
}

// 封装图片上传函数
// desPath 是相对路径
func UploadImg(this *beego.Controller, inputKey, desPath string) (string, error) {
	file, head, err := this.GetFile(inputKey)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 大小
	if head.Size > 5000000 {
		return "", errors.New("文件超过5M")
	}

	// 类型
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		return "", errors.New("类型错误")
	}

	// 避免重名
	currentTimeStr := time.Now().Format("2006-01-02-15:04:05")
	savePath := path.Join(desPath, currentTimeStr+ext)
	err = this.SaveToFile(inputKey, savePath)
	if err != nil {
		return "", err
	}

	return "/"+savePath, nil
}