package database

import (
	"fmt"
	"log"

	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.GetEnv("DB_USER", "root"),
		env.GetEnv("DB_PASSWORD", "password"),
		env.GetEnv("DB_HOST", "localhost"),
		env.GetEnv("DB_PORT", "3306"),
		env.GetEnv("DB_NAME", "impl_messaging"),
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if DB == nil {
		panic("Database connection is not initialized")
	}
	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println(DB)
	err = DB.AutoMigrate(&models.User{}, &models.UserSession{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	fmt.Println("successfully migrate database!")
	
}
