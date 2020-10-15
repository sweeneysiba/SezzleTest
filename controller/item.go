package controller

import (
	"SezzleTest/config"
	"SezzleTest/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//CreateItem ...
func CreateItem(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	var ItemInput models.Item
	err = json.Unmarshal(body, &ItemInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}

	if ItemInput.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | name  are mandatory"})
		return
	}

	timeStamp := time.Now().UnixNano()
	ItemInput.ID = timeStamp
	ItemInput.CreatedAt = timeStamp / int64(time.Millisecond)
	PrettyPrint(ItemInput)
	if err := config.DB.Create(ItemInput).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | please check the request body ", "error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item Created Successfully ."})
	return
}

//ListItems ...
func ListItems(c *gin.Context) {
	var ItemArr []models.Item
	if err := config.DB.Raw("SELECT name , created_at ,id FROM items ").Scan(&ItemArr).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "items list fetched Successfully", "items": ItemArr})
	return
}
