package controller

import (
	"SezzleTest/config"
	"SezzleTest/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ListOrders ...
func ListOrders(c *gin.Context) {
	var usersArr []models.Orders
	if err := config.DB.Raw("SELECT cart_id , created_at ,id , user_id FROM orders ").Scan(&usersArr).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "orders not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to fetch orders", "error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order list fetched Successfully", "orders": usersArr})
	return

}
