package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
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

type User struct {
	ID       uint
	Name     string
	Age      int
	UserInfo *UserInfo
}

type UserInfo struct {
	UserID uint // 外键
	ID     uint
	Addr   string
	Like   string
}
type Tag struct {
	ID   uint
	Name string
	//Articles []Article `gorm:"many2many:article_tags;"` // 用于反向引用
}
type Article struct {
	ID    uint
	Title string
	Tags  []Tag `gorm:"many2many:article_tags;"`
}

type ArticleTag struct {
	ArticleID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
	CreateAt  time.Time
}

type ArticleModel struct {
	ID    uint
	Title string
	Tags  []TagModel `gorm:"many2many:article_tags;joinForeignKey:ArticleID;JoinReferences:TagID"`
}

type TagModel struct {
	ID       uint
	Name     string
	Articles []ArticleModel `gorm:"many2many:article_tags;joinForeignKey:TagID;JoinReferences:ArticleID"` // 用于反向引用
}

type ArticleTagModel struct {
	ArticleID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
	CreateAt  time.Time
}

func main() {
	InitDB()
	//DB.AutoMigrate(&User{}, &UserInfo{})
	//DB.Create(&User{
	//	Name: "你不知道",
	//	Age:  8,
	//	UserInfo: &UserInfo{
	//		Addr: "天堂",
	//		Like: "地狱",
	//	},
	//})
	// 根据用户查询用户详细
	//var u User
	//DB.Preload("UserInfo").Take(&u)
	//s, _ := json.Marshal(u)
	//fmt.Println(string(s))
	// 根据用户详情查询用户    需要相互关联才行
	//var userInfo UserInfo
	//DB.Take(&userInfo)
	//fmt.Println(userInfo)
	// 删除 先查出来再删除
	//var u User
	//DB.Debug().Take(&u) //  SELECT * FROM `tb_user` LIMIT 1
	//DB.Debug().Select("UserInfo").Delete(&u)
	// DELETE FROM `tb_user_info` WHERE `tb_user_info`.`user_id` = 2
	// DELETE FROM `tb_user` WHERE `tb_user`.`id` = 2
	//DB.AutoMigrate(&Tag{}, &Article{})
	// 多对多的添加
	//DB.Debug().Create(&Article{
	//	Title: "golang学习",
	//	Tags: []Tag{
	//		{
	//			Name: "go",
	//		},
	//		{
	//			Name: "goo",
	//		},
	//	},
	//})
	////	执行sql如下
	//// INSERT INTO `tb_tag` (`name`) VALUES ('go'),('goo') ON DUPLICATE KEY UPDATE `id`=`id`
	//// INSERT INTO `tb_article_tags` (`article_id`,`tag_id`) VALUES (1,1),(1,2) ON DUPLICATE KEY UPDATE `article_id`=`article_id`
	//// INSERT INTO `tb_article` (`title`) VALUES ('golang学习')

	//var tag Tag
	//DB.Take(&tag, "name = ?", "goo")
	//tags := []Tag{tag, Tag{Name: "xxxx"}}
	//DB.Create(&Article{Title: "python基础", Tags: tags})

	//	查询文章，同时显示标签
	//var a Article
	//DB.Debug().Preload("Tags").Take(&a)
	// SELECT * FROM `tb_article_tags` WHERE `tb_article_tags`.`article_id` = 1
	// SELECT * FROM `tb_tag` WHERE `tb_tag`.`id` IN (1,2)
	// SELECT * FROM `tb_article` LIMIT 1
	//fmt.Println(a.Tags)
	//for _, tag := range a.Tags {
	//	fmt.Println(tag.Name)
	//}
	//marshal, _ := json.Marshal(&a.Tags)
	//fmt.Println(string(marshal))

	// 多对多的更新
	// 先删除原有的标签
	//var article Article
	//DB.Preload("Tags").Take(&article, 1)
	//DB.Model(&article).Association("Tags").Delete(article.Tags)

	// 在添加新的标签
	//var tag Tag
	//DB.Take(&tag, "1")
	//DB.Model(&article).Association("Tags").Append(&tag)

	// 直接替换标签
	//var article Article
	//DB.Preload("Tags").Take(&article, 1)
	//var tag Tag
	//DB.Take(&tag, "3")
	//DB.Model(&article).Association("Tags").Replace(&tag)

	//// 设置Article的Tag表为ArticleTag
	//DB.SetupJoinTable(&ArticleModel{}, "Tags", &ArticleTagModel{})
	//// 如果Tag要反向引用Articles，也需要加上
	//DB.SetupJoinTable(&TagModel{}, "Articles", &ArticleTagModel{})
	//if err := DB.AutoMigrate(&ArticleModel{}, &TagModel{}, &ArticleTagModel{}); err != nil {
	//	log.Panic(err)
	//}
}
