package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupDatabase() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True",
		DatabaseSetting.User,
		DatabaseSetting.Password,
		DatabaseSetting.Host,
		DatabaseSetting.Name)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	db.AutoMigrate(&PriceData{})
	db.AutoMigrate(&SellData{})
	db.AutoMigrate(&Product{})
}

type PriceData struct {
	ProductID uint      `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	Price     uint      `json:"price"`
}

type SellData struct {
	ProductID uint      `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID         uint        `json:"id"`
	Name       string      `json:"name"`
	StartPrice uint        `json:"start_price"`
	SellData   []SellData  `json:"sell_data" gorm:"constraint:OnDelete:CASCADE;"`
	PriceData  []PriceData `json:"price_data" gorm:"constraint:OnDelete:CASCADE;"`
}

func GetProducts(db *gorm.DB, from time.Time, to time.Time) ([]Product, error) {
	var products []Product
	err := db.Model(&Product{}).
		Preload("PriceData", func(db *gorm.DB) *gorm.DB {
			return db.Where("created_at BETWEEN ? AND ?", from, to).Order("created_at ASC")
		}).
		Preload("SellData", func(db *gorm.DB) *gorm.DB {
			return db.Where("created_at BETWEEN ? AND ?", from, to).Order("created_at ASC")
		}).
		Find(&products).Error

	return products, err
}

func SellProduct(db *gorm.DB, id uint) error {
	sellData := SellData{ProductID: id}
	err := db.Create(&sellData).Error
	return err
}

func SetNewPrice(db *gorm.DB, id uint, price uint) error {
	priceData := PriceData{ProductID: id, Price: price}
	err := db.Create(&priceData).Error
	return err
}

func CloseDB() {
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}
