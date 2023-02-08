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
		SkipDefaultTransaction: true,
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

/*
// 存储json格式数据

type Info struct {
	Status string `json:"status"`
	Addr   string `json:"addr"`
	Age    int    `json:"age"`
}

// Scan 从数据库中读取
func (i *Info) Scan(v interface{}) error {
	bytes, ok := v.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal json value: ", v))
	}
	return json.Unmarshal(bytes, i)
}

// Value 存储数据
func (i Info) Value() (driver.Value, error) {
	return json.Marshal(i)
}

type AuthModel struct {
	ID   uint
	Name string
	Info Info `gorm:"type:string"`
}

func main() {
	InitDB()
	DB.AutoMigrate(&AuthModel{})
	//DB.Debug().Create(&AuthModel{
	//	Name: "测试007",
	//	Info: Info{
	//		Status: "我姮好",
	//		Addr:   "1.1.1.1",
	//		Age:    10,
	//	},
	//})

	var u AuthModel
	DB.Debug().Take(&u)
	fmt.Println(u)
}
*/

/*
type Array []string

// Scan 从数据库中读取
func (i *Array) Scan(v interface{}) error {
	bytes, ok := v.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("解析失败: %v  %T", v, v))
	}
	*i = strings.Split(string(bytes), "|")
	return nil
}

// Value 存储数据
func (i Array) Value() (driver.Value, error) {
	return strings.Join(i, "|"), nil
}

type HostModel struct {
	ID    uint
	IP    string
	Ports Array `gorm:"type:string"`
}

func main() {
	InitDB()
	DB.AutoMigrate(&HostModel{})
	DB.Create(&HostModel{
		IP:    "1.1.1.1",
		Ports: []string{"1", "2", "3"},
	})
}
*/

// 枚举
const (
	Running Status = 1
	Except  Status = 2
	OffLine Status = 3
)

type Status int

func (s Status) MarshalJSON() ([]byte, error) {
	var str string
	switch s {
	case Running:
		str = "Running"
	case Except:
		str = "Except"
	case OffLine:
		str = "OffLine"
	}
	return json.Marshal(str)
}

type Host struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"`
}

func main() {
	h := Host{1, "我们我们", Running}
	bytes, _ := json.Marshal(h)
	fmt.Println(string(bytes))
	InitDB()
	DB.AutoMigrate(&Host{})
	DB.Debug().Create(&Host{Name: "99999", Status: Except})
	// INSERT INTO `tb_host` (`name`,`status`) VALUES ('99999','2')
	var host Host
	DB.Take(&host)
	marshal, _ := json.Marshal(&host)
	fmt.Println(string(marshal))
	// {"id":1,"name":"99999","status":"Except"}
}
