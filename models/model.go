package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// models中放的是 表的设计

// 2. 创建类
type User struct {
	Id int				// id	如果没指定主键, orm自动将int类型的id字段当成主键, 并有自增特性
	Name string			// name
	PassWord string		// 生成的表字段为pass_word, 自动将驼峰转成下划线
	// 在orm中 `__` 有特殊含义
	// Pass_Word  string //生成字段为 pass__word, 不要这么定义字段名

	Articles []*Article		`orm:"reverse(many)"`
}

type Article struct {
	Id int 				`orm:"pk;auto"`
	ArtiName string 	`orm:"size(20)"`
	Atime	 time.Time	`orm:"auto_now"`
	Acount	 int		`orm:"default(0);null"`	// orm默认就是not null
	Acontent string		`orm:"size(500)"`
	Aimg 	 string		`orm:"size(100)"`

	// 一对多, 在多的那个表中 加`orm:"rel(fk)"`	fk表示外键
	ArticleType *ArticleType	`orm:"rel(fk)"`	// 正向设置
	// 多对多, 在两表中加上 `orm:"rel(m2m)"` many to many		`orm:"reverse(many)"`
	Users []*User		`orm:"rel(m2m)"`
}

type ArticleType struct {
	Id int
	TypeName string		`orm:"size(20)"`

	// 一对多, 在一的那个表中 加`orm:"reverse(many)"`
	Articles []*Article		`orm:"reverse(many)"`	// 反向设置
}

func init() {
	/*
		* 方式一: 用sql语句操作数据库
		// 设置链接数据库参数
		// 第一个参数: 数据库驱动
		// 第二个参数: 连接数据库字符串	---	用户名:密码@[链接方式](地址:端口)/数据库名?charset=utf8
		conn, err := sql.Open("mysql", "root:zj2fighting@(192.168.1.7:3306)/zj?charset=utf8")
		if err != nil {
			logs.Warning("sql.Open err:", err)
			logs.Error("sql.Open err:", err)
			return
		}
		defer conn.Close()

		// 创建表
		_, err = conn.Exec("create table zjTest(name char(20) , password char(20))")
		if err != nil {
			logs.Error("conn.Exec err:", err)
			return
		}

		// 增删改查
		//conn.Exec("insert into zjTest (name, password) values (?,?)", "zj", "gene")

		rows, _ := conn.Query("select name from zjTest")
		var name string
		for rows.Next() {
			rows.Scan(&name)
			logs.Warning(name)
		}
	*/

	// 方式二: orm方式操作数据库
	// 1. 注册数据库
	// 第一个参数数据库别名 default 必须有
	orm.RegisterDataBase("default", "mysql", "root:zj2fighting@(192.168.1.8:3306)/zj?charset=utf8&loc=Local")
	// 3. 注册表关联类
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	// 4. 生成表
	// 第一个参数别名default
	// 第二个参数是否强制刷新表
	//      设置成true后, 将原表drop, 再生成新表, 原表数据丢失
	// 		设置成false后, 如果表已经存在, 此句无效果(不会更新表字段). 如果表不存在会创建表
	// 第三个参数过程是否可见(生成表语句的提示过程)
	orm.RunSyncdb("default", false, true)

	// 5. 操作表一般放在Controller中
}
