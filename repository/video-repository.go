package repository

import (
	"awesomeProject/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type VideoRepository interface {
	Save(video entity.Video)
	Update(video entity.Video)
	Delete(video entity.Video)
	FindAll() []entity.Video
}

type database struct {
	connection *gorm.DB
}

func NewVideoRepository() VideoRepository {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=postgres password=root dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic("fail to connect db")
	}
	db.AutoMigrate(&entity.Video{}, &entity.Person{})
	return &database{
		connection: db,
	}
}

func (db *database) Save(video entity.Video) {
	db.connection.Create(&video)
}
func (db *database) Update(video entity.Video) {
	db.connection.Save(&video)
}
func (db *database) Delete(video entity.Video) {
	db.connection.Delete(&video)
}
func (db *database) FindAll() []entity.Video {
	var videos []entity.Video
	db.connection.Set("gorm:auto_preload", true).Find(&videos)
	return videos
}
