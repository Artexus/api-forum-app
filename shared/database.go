package shared

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

func Init() {
	//PostgreSQL
	text := fmt.Sprintf("host=%v user=%v port=%v dbname=%v sslmode=disable password=%v",
		os.Getenv("DB_HOST_NAME"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	db, err = gorm.Open("postgres", text)
	if err != nil {
		log.Println(text)
		log.Fatal(fmt.Errorf("[ERROR] connect to database: %s", err.Error()))
		return
	}
	db.LogMode(true)
}

func GetDb() *gorm.DB {
	return db
}

func CloseDb() {
	db.Close()
}
