package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dsn  string
	port uint
)

var DB *gorm.DB

func init() {
	flag.StringVar(&dsn, "dsn", "root:root@tcp(127.0.0.1:3306)/pdp", "mysql数据库连接")
	flag.UintVar(&port, "port", 8066, "服务端口号")
	flag.Parse()

	dsn = fmt.Sprintf("%v?charset=utf8mb4&parseTime=True&loc=Local", dsn)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(50)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(50)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Minute)
	if err := DB.AutoMigrate(&Questionnaire{}); err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()

	r.Use(cors.Default())
	r.POST("/", FormPost)
	r.Run(fmt.Sprintf(":%v", port))
}

func FormPost(c *gin.Context) {
	var q Questionnaire
	if err := c.ShouldBind(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	q.SetResult()
	if err := DB.Create(&q).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, q)
}

type Questionnaire struct {
	ID         int       `gorm:"primaryKey;autoIncrement;column:id" db:"id" json:"id" form:"id"`
	Name       string    `gorm:"column:name;not null;size:256" db:"name" json:"name" form:"name" binding:"required"`
	Result     string    `gorm:"column:result;not null;size:256" db:"result" json:"result" form:"result"`
	CreatedAt  time.Time `gorm:"column:created_at" db:"created_at" json:"created_at" form:"created_at"`
	Question1  uint8     `gorm:"column:question_1" db:"question_1" json:"question_1" form:"question_1"`
	Question2  uint8     `gorm:"column:question_2" db:"question_2" json:"question_2" form:"question_2"`
	Question3  uint8     `gorm:"column:question_3" db:"question_3" json:"question_3" form:"question_3"`
	Question4  uint8     `gorm:"column:question_4" db:"question_4" json:"question_4" form:"question_4"`
	Question5  uint8     `gorm:"column:question_5" db:"question_5" json:"question_5" form:"question_5"`
	Question6  uint8     `gorm:"column:question_6" db:"question_6" json:"question_6" form:"question_6"`
	Question7  uint8     `gorm:"column:question_7" db:"question_7" json:"question_7" form:"question_7"`
	Question8  uint8     `gorm:"column:question_8" db:"question_8" json:"question_8" form:"question_8"`
	Question9  uint8     `gorm:"column:question_9" db:"question_9" json:"question_9" form:"question_9"`
	Question10 uint8     `gorm:"column:question_10" db:"question_10" json:"question_10" form:"question_10"`
	Question11 uint8     `gorm:"column:question_11" db:"question_11" json:"question_11" form:"question_11"`
	Question12 uint8     `gorm:"column:question_12" db:"question_12" json:"question_12" form:"question_12"`
	Question13 uint8     `gorm:"column:question_13" db:"question_13" json:"question_13" form:"question_13"`
	Question14 uint8     `gorm:"column:question_14" db:"question_14" json:"question_14" form:"question_14"`
	Question15 uint8     `gorm:"column:question_15" db:"question_15" json:"question_15" form:"question_15"`
	Question16 uint8     `gorm:"column:question_16" db:"question_16" json:"question_16" form:"question_16"`
	Question17 uint8     `gorm:"column:question_17" db:"question_17" json:"question_17" form:"question_17"`
	Question18 uint8     `gorm:"column:question_18" db:"question_18" json:"question_18" form:"question_18"`
	Question19 uint8     `gorm:"column:question_19" db:"question_19" json:"question_19" form:"question_19"`
	Question20 uint8     `gorm:"column:question_20" db:"question_20" json:"question_20" form:"question_20"`
	Question21 uint8     `gorm:"column:question_21" db:"question_21" json:"question_21" form:"question_21"`
	Question22 uint8     `gorm:"column:question_22" db:"question_22" json:"question_22" form:"question_22"`
	Question23 uint8     `gorm:"column:question_23" db:"question_23" json:"question_23" form:"question_23"`
	Question24 uint8     `gorm:"column:question_24" db:"question_24" json:"question_24" form:"question_24"`
	Question25 uint8     `gorm:"column:question_25" db:"question_25" json:"question_25" form:"question_25"`
	Question26 uint8     `gorm:"column:question_26" db:"question_26" json:"question_26" form:"question_26"`
	Question27 uint8     `gorm:"column:question_27" db:"question_27" json:"question_27" form:"question_27"`
	Question28 uint8     `gorm:"column:question_28" db:"question_28" json:"question_28" form:"question_28"`
	Question29 uint8     `gorm:"column:question_29" db:"question_29" json:"question_29" form:"question_29"`
	Question30 uint8     `gorm:"column:question_30" db:"question_30" json:"question_30" form:"question_30"`
}

// 把第5、10、14、18、24、30题的分加起来就是你的“老虎”分数
// 把第3、6、13、20、22、29题的分加起来就是你的“孔雀”分数
// 把第2、8、15、17、25、28题的分加起来就是你的“考拉”分数
// 把第1、7、11、16、21、26题的分加起来就是你的“猫头鹰”分数
// 把第4、9、12、19、23、27题的分加起来就是你的“变色龙”分数
func (q *Questionnaire) SetResult() {
	// 老虎
	tiger := q.Question5 + q.Question10 + q.Question14 + q.Question18 + q.Question24 + q.Question30
	// 孔雀
	peacock := q.Question3 + q.Question6 + q.Question13 + q.Question20 + q.Question22 + q.Question29
	// 考拉
	kangaroo := q.Question2 + q.Question8 + q.Question15 + q.Question17 + q.Question25 + q.Question28
	// 猫头鹰
	owl := q.Question1 + q.Question7 + q.Question11 + q.Question16 + q.Question21 + q.Question26
	// 变色龙
	chameleon := q.Question4 + q.Question9 + q.Question12 + q.Question19 + q.Question23 + q.Question27

	if tiger >= peacock && tiger >= kangaroo && tiger >= owl && tiger >= chameleon {
		q.Result = "老虎"
	} else if peacock >= tiger && peacock >= kangaroo && peacock >= owl && peacock >= chameleon {
		q.Result = "孔雀"
	} else if kangaroo >= tiger && kangaroo >= peacock && kangaroo >= owl && kangaroo >= chameleon {
		q.Result = "考拉"
	} else if owl >= tiger && owl >= peacock && owl >= kangaroo && owl >= chameleon {
		q.Result = "猫头鹰"
	} else if chameleon >= tiger && chameleon >= peacock && chameleon >= kangaroo && chameleon >= owl {
		q.Result = "变色龙"
	} else {
		q.Result = "未知"
	}
}
