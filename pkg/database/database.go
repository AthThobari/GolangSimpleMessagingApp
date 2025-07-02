package database

import (
	"gorm.io/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *gorm.DB

var MongoDB *mongo.Collection