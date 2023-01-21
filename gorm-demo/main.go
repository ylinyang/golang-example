package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var DB *gorm.DB

func InitDB() {
	host := "127.0.0.1"
	port := 3306
	username := "root"
	password := "123456"
	dbname := "db_test"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tb_", //表名前缀
			SingularTable: true,  // 单数表名
			//NoLowerCase:   false, //打开小写转换
		},
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln("连接数据库失败,", err)
		return
	}
	DB = db
	log.Println("连接数据成功")
}

type Student struct {
	ID     uint   `gorm:"size:3;primaryKey"`
	Name   string `gorm:"type:varchar(12);comment:用户名"`
	Age    int    `gorm:"size:4"`
	Gender bool
	Email  *string `gorm:"size:32"`
}

//func (s *Student) BeforeCreate(tx *gorm.DB) (err error) {
//	s.Age = 99
//	return nil
//}

// User 用户表
type User struct {
	Id       uint
	Name     string
	Articles []Article `gorm:"foreignKey:UserId"`
}

// Article 文章列表
type Article struct {
	Id     uint
	Title  string
	UserId uint
	User   User `gorm:"foreignKey:UserId"`
}

func main() {
	InitDB()
	//DB.AutoMigrate(&Article{}, &User{})
	// 单纯创建
	//DB.Debug().Create(&User{
	//	Name: "张三",
	//	Articles: []Article{
	//		{
	//			Title: "golang",
	//		},
	//		{
	//			Title: "python",
	//		},
	//	},
	//})

	//	 创建关联已有表内容
	//DB.Debug().Create(&Article{
	//	Title: "测试001",
	//	User: User{
	//		Name: "张三",
	//	},
	//})
	//DB.Debug().Create(&Article{
	//	Title:  "测试002",
	//	UserId: 2,
	//})
	// 给已有用户绑定文章
	//var User User
	//DB.Take(&User, 1)
	//var article Article
	//DB.Take(&article, 6)
	//DB.Model(&User).Association("Articles").Append(&article)
	////	 预加载
	//var article Article
	//DB.Preload("User").Take(&article)
	//fmt.Println(article) // {1 golang 1 {1 张三 []}}
	// 查询某个用户下面有多少个文章
	//var user User
	//// 查询全部用户的
	//DB.Preload("Articles").Take(&user)
	//// 查询文章id小于2的
	//DB.Preload("Articles", "id < ?", 2).Take(&user)
	//
	//// 等价于 如下
	//DB.Preload("Articles", func(db *gorm.DB) *gorm.DB {
	//	return db.Where("id < ? ", 2)
	//}).Take(&user)
	//fmt.Println(user)

	// 清楚外键关联，并不会真正删除数据
	// 删除用户，删除相关文章的用户ID
	//var user User
	//DB.Preload("Articles").Take(&user, 2)
	//fmt.Println(user)
	//DB.Model(&user).Association("Articles").Delete(&user.Articles)
	//	DB.Delete(&user)

	// 通过select用户数据和文章数据都会被删除
	var user User
	DB.Preload("Articles").Take(&user, 1)
	DB.Select("Articles").Delete(&user)
}

// 插入单条、多条记录
func insert() {
	//email := "x@qq.com"
	s := Student{
		Name: "0010-test",
		//Age:    32,
		Gender: true,
		//Email: &email,
	}
	if DB.Create(&s).Error != nil {
		log.Println("创建成功")
	}

	var studentList []Student
	for i := 0; i < 3; i++ {
		student := Student{
			Name:  fmt.Sprintf("100%d", i),
			Age:   i,
			Email: nil,
		}
		studentList = append(studentList, student)
	}
	if DB.Create(&studentList).Error != nil {
		return
	}
}

// 查询单条、多条数据
func find() {
	// 随机查询一条
	//var s Student
	//DB.Debug().Take(&s)      // SELECT * FROM `tb_student` LIMIT 1
	//fmt.Println(s, *s.Email) // {1 008-test 80 false 0x14000220ae0} 8@qq.com 最后一个存储的是指针

	// 按照主键查询
	var priKey Student
	DB.Debug().Take(&priKey, 6) // SELECT * FROM `tb_student` WHERE `tb_student`.`id` = 6 LIMIT 1
	fmt.Println(priKey)

	var priKeyNotFound Student
	if DB.Debug().Take(&priKeyNotFound, 100).Error == gorm.ErrRecordNotFound {
		fmt.Println("主键不存在")
	}

	// 按照其他字段进行查询，但是也只查询一条
	var otherKey Student
	DB.Debug().Take(&otherKey, "Email = ?", "8@qq.com") // SELECT * FROM `tb_student` WHERE name = '008-test' LIMIT 1
	fmt.Println(otherKey)

	// 查询第一条
	var s Student
	DB.Debug().First(&s) // SELECT * FROM `tb_student` ORDER BY `tb_student`.`id` LIMIT 1
	fmt.Println(s)

	// 查询最后一条
	var l Student
	DB.Debug().Last(&l) //  SELECT * FROM `tb_student` ORDER BY `tb_student`.`id` DESC LIMIT 1
	fmt.Println(l)

	// 获取查询的记录数
	fmt.Println(DB.Find(&otherKey).RowsAffected) // 1

	fmt.Println("-----------------")
	//查询多条记录 不跟条件默认查询全部
	var studentList []Student
	DB.Debug().Find(&studentList, "Email = ?", "8@qq.com") // SELECT * FROM `tb_student` WHERE Email = '8@qq.com'
	for _, v := range studentList {
		fmt.Println(v)
	}
	//  {1 008-test 80 false 0x140002210d0}
	//  {6 0098-test 11 false 0x140002210f0}

	//	由于Email为指针类型，通过序列化，转化可以直接查看
	for _, v := range studentList {
		data, _ := json.Marshal(&v)
		fmt.Println(string(data))
	}
	// {"ID":1,"Name":"008-test","Age":80,"Gender":false,"Email":"8@qq.com"}
	// {"ID":6,"Name":"0098-test","Age":11,"Gender":false,"Email":"8@qq.com"}

	// 根据主键列表去查询
	var studentListByPriKey []Student
	DB.Debug().Find(&studentListByPriKey, []int{1, 3, 5}) // SELECT * FROM `tb_student` WHERE `tb_student`.`id` IN (1,3,5)
	fmt.Println(studentListByPriKey)                      // [{1 008-test 80 false 0x14000221350} {3 0010-test 0 true <nil>} {5 1000 0 false <nil>}]

	// 根据其他条件去查询
	var studentListByOther []Student
	DB.Debug().Find(&studentListByOther, "name in (?)", []string{"0098-test", "0099-test"}) // SELECT * FROM `tb_student` WHERE name in ('0098-test','0099-test')
	fmt.Println(studentListByOther)                                                         // [{6 0098-test 11 false 0x14000285360} {7 0099-test 32 true 0x14000285380}]
}

// 更新数据
func update() {
	// save 用于保存所有字段，即使是零值也会保存
	var studentOne Student
	DB.Debug().Take(&studentOne, "name = ?", "007-test") // SELECT * FROM `tb_student` WHERE name = '007-test' LIMIT 1
	fmt.Println(&studentOne)
	studentOne.Age = 20
	DB.Debug().Select("age").Save(&studentOne) // UPDATE `tb_student` SET `name`='20259-test' WHERE `id` = 7
	//DB.Debug().Save(&studentOne) // UPDATE `tb_student` SET `name`='20259-test',`age`=32,`gender`=true,`email`='99@qq.com' WHERE `id` = 7

	//var studentList []Student
	//DB.Debug().Find(&studentList, []int{1, 2}).Update("gender", true) // UPDATE `tb_student` SET `gender`=true WHERE `tb_student`.`id` IN (1,2) AND `id` IN (1,2)

	//DB.Debug().Find(&studentList, []int{2, 3}).Updates(Student{
	//	Age:    100,
	//	Gender: true,
	//}) // UPDATE `tb_student` SET `age`=100,`gender`=true WHERE `tb_student`.`id` IN (2,3) AND `id` IN (2,3)

	//DB.Debug().Find(&studentList, []int{3, 4}).Updates(map[string]any{
	//	"name": "007-f",
	//}) //  UPDATE `tb_student` SET `name`='007-f' WHERE `tb_student`.`id` IN (3,4) AND `id` IN (3,4)
}

// 删除数据
func deleteData() {
	var student Student
	DB.Debug().Delete(&student, 3)           // DELETE FROM `tb_student` WHERE `tb_student`.`id` = 3
	DB.Debug().Delete(&student, []int{2, 5}) // DELETE FROM `tb_student` WHERE `tb_student`.`id` IN (2,5)
}

// 条件查询
func where() {
	//	 查询用户名是张三的
	//var studentList []Student
	//DB.Debug().Where("name = ? ", "张三").Find(&studentList) // SELECT * FROM `tb_student` WHERE name = '张三'
	//DB.Debug().Find(&studentList, "name = ?", "张三")
	// 查询用户名不是张三的
	//var studentList []Student
	//DB.Debug().Not("name = ? ", "张三").Find(&studentList) // SELECT * FROM `tb_student` WHERE NOT name = '张三'

	//var studentList []Student
	//DB.Debug().Where("not name = ?", "张三").Find(&studentList)

	// 查询用户名包含张三、李四的
	//var studentList []Student
	//DB.Debug().Where("name in ?", []string{"张三", "李四"}).Find(&studentList) // SELECT * FROM `tb_student` WHERE name in ('张三','李四')
	//DB.Find(&studentList, "name in (?)", []string{"张三", "李四"})

	// 模糊匹配 用户名带杨的
	//var studentList []Student
	//DB.Debug().Where("name like ?", "杨%").Find(&studentList) // SELECT * FROM `tb_student` WHERE name like '杨%'
	//	DB.Debug().Where("name like ?", "杨_").Find(&studentList) // SELECT * FROM `tb_student` WHERE name like '杨_'

	// 查找age大于20，且邮箱是qq的
	// var studentList []Student
	//DB.Debug().Where("age > ? and email like ?", 22, "%@qq.com").Find(&studentList) // SELECT * FROM `tb_student` WHERE age > 22 and email like '%@qq.com'
	//DB.Where("age > ?", 22).Where("email like ?", "%@qq.com").Find(&studentList)

	// 查询gender 为false，且邮箱是qq的
	//var studentList []Student
	// 	DB.Debug().Where("gender = ? or email like ?", false, "%@qq.com").Find(&studentList) // SELECT * FROM `tb_student` WHERE gender = false or email like '%@qq.com'
	//DB.Where("gender = ?", false).Or("email like ?", "%qq@.com").Find(&studentList)

	// 使用结构体查询
	//var studentList []Student
	//DB.Debug().Where(&Student{Name: "张三"}).Find(&studentList) // SELECT * FROM `tb_student` WHERE `tb_student`.`name` = '张三'
	//DB.Debug().Where(map[string]any{
	//	"name": "张三",
	//	"age":  "99",
	//}).Find(&studentList)

	//var studentList []Student
	//DB.Debug().Select("name").Find(&studentList) // SELECT `name` FROM `tb_student`
	//fmt.Println(studentList) // [{0 张三 0 false <nil>} {0 李四 0 false <nil>}]

	// scan 语法
	//var studentList []Student
	//type User struct {
	//	Name string
	//	Age  int
	//}
	//var userList []User
	//DB.Debug().Select("name", "age").Limit(2).Find(&studentList).Scan(&userList)
	//// SELECT `name`,`age` FROM `tb_student` LIMIT 2
	//fmt.Println(studentList) //[{0 张三 20 false <nil>} {0 李四 22 false <nil>}]
	//fmt.Println(userList)    // [{张三 20} {李四 22}]
	//type User struct {
	//	Title string `gorm:"column:name"`
	//	Age   int
	//}
	//var userList []User
	//DB.Model(Student{}).Select("name", "age").Scan(&userList)
	//fmt.Println(userList)

	// 按照年龄排序 desc 降序 asc 升序
	//var studentList []Student
	//DB.Debug().Order("age desc").Find(&studentList)

	//var studentList []Student
	// 每页两条数据  查询第一页
	//DB.Debug().Limit(2).Offset(0).Find(&studentList) // SELECT * FROM `tb_student` LIMIT 2
	// 每页两条数据 查询第二页  offset 为 (页数 - 1) * 查询条数
	//DB.Debug().Limit(2).Offset(2).Find(&studentList) //  SELECT * FROM `tb_student` LIMIT 2 OFFSET 2
	//fmt.Println(studentList) // [{5 杨四 99 true 0x14000220ac0} {4 杨七一 24 false 0x14000220ae0} {3 老五 23 false 0x14000220b00}]

	// 按照年龄去重  Distinct去除重复字段
	//var ageList []int
	//DB.Debug().Model(Student{}).Select("age").Distinct("age").Scan(&ageList) // SELECT DISTINCT `age` FROM `tb_student`
	//DB.Debug().Model(Student{}).Select("Distinct age").Scan(&ageList)
	//fmt.Println(ageList)

	// 将 gender 统计分组
	//type Group struct {
	//	Count    int
	//	Gender   string
	//	NameList string
	//}
	//var groupList []Group
	//DB.Model(Student{}).Select("group_concat(name) as name_list", "count(id) as count", "gender").Group("gender").Scan(&groupList)
	//fmt.Println(groupList) // [{4 0 张三,李四,老五,杨七一} {1 1 杨四}]

	//type Name struct {
	//	Count int `gorm:"column:count(id)"`
	//	Name  string
	//}
	//var nameList []Name
	//DB.Model(Student{}).Select("count(id)", "name").Group("name").Scan(&nameList)
	//fmt.Println(nameList)

	// 执行原生sql
	//var studentList []Student
	//DB.Debug().Raw("select * from tb_student where name =?", "张三").Find(&studentList) // select * from tb_student where name ='张三'

	//	查询大于平均年龄的用户
	//var studentList []Student
	//DB.Debug().Where("age > (?)", DB.Model(Student{}).Select("avg(age)")).Find(&studentList) // SELECT * FROM `tb_student` WHERE age > (SELECT avg(age) FROM `tb_student`)

	// 命名参数
	//var student Student
	//DB.Debug().Take(&student, "name = @name and age = @age", sql.Named("name", "杨七一"), sql.Named("age", "24")) // SELECT * FROM `tb_student` WHERE name = '杨七一' and age = '24' LIMIT 1

	//var student Student
	//DB.Debug().Take(&student, map[string]any{
	//	"name": "张三",
	//	"age":  20,
	//})
	//fmt.Println(student) // {2 李四 22 false 0x1400021eb30}

	//var res []map[string]any
	//DB.Debug().Model(&Student{}).Where("age = 20").Find(&res)
	//fmt.Println(res) // [map[age:20 email:0x14000220ad0 gender:false id:1 name:张三] map[age:20 email:0x14000220b00 gender:true id:5 name:杨四]]

	//	查询年龄大于20的
	var res []map[string]any
	DB.Debug().Model(Student{}).Scopes(age).Find(&res) // SELECT * FROM `tb_student` WHERE age > 20
}

func age(db *gorm.DB) *gorm.DB {
	return db.Where("age > ?", 20)
}

/*
CREATE TABLE `tb_students` (
	`name` varchar(12) NOT NULL COMMENT '用户名',
	`uuid` varchar(191) UNIQUE COMMENT '主键',
	`_semail` varchar(32),
	`_sy_addr` varchar(16),
	`_sgender` boolean DE true,
	PRIMARY KEY (`uuid`)
)

*/
