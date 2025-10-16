package database

import (
	"context"
	"fmt"
	"log"

	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	log.Println("successfully migrate database!")

}

func SetupMongoDB() {
	uri := env.GetEnv("MONGODB_URI", "")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Uji koneksi (penting untuk pastikan sukses konek)
	if err := client.Ping(context.TODO(), nil); err != nil {
		panic(fmt.Errorf("failed to ping MongoDB: %w", err))
	}

	coll := client.Database("message").Collection("message_history")
	MongoDB = coll

	log.Println("successfully connected to MongoDB")
}
