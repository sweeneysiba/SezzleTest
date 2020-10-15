//main.go
package main

import (
	"SezzleTest/config"
	"SezzleTest/models"
	Routes "SezzleTest/routes"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var err error

func Sum(x int, y int) int {
	return x + y
}

func main() {
	config.DB, err = gorm.Open("mysql", config.DbURL(config.BuildDBConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}

	defer config.DB.Close()
	config.DB.AutoMigrate(&models.User{}, &models.Item{}, &models.CartItem{}, &models.Cart{}, &models.Orders{})
	config.DB.Model(&models.CartItem{}).AddForeignKey("cart_id", "carts(id)", "RESTRICT", "CASCADE")
	config.DB.Model(&models.CartItem{}).AddForeignKey("item_id", "items(id)", "RESTRICT", "CASCADE")
	config.DB.Model(&models.Orders{}).AddForeignKey("cart_id", "carts(id)", "RESTRICT", "CASCADE")
	config.DB.Model(&models.Orders{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "CASCADE")
	// config.DB.Model(&models.RatePlan{}).AddForeignKey("hotel_id", "hotels(hotel_id)", "RESTRICT", "RESTRICT")
	r := Routes.SetupRouter()
	//running
	r.Run(":8080")
}
