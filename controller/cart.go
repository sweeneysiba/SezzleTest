package controller

import (
	"SezzleTest/config"
	"SezzleTest/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ListCart ...
func ListCart(c *gin.Context) {
	var CartArray []models.Cart
	if err := config.DB.Raw("SELECT cart_id , created_at ,id , user_id FROM cart ").Scan(&CartArray).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item list fetched Successfully", "items": CartArray})
	return
}

//AddToCart ...
func AddToCart(c *gin.Context) {
	getId := c.MustGet("user_id")
	var ItemArr []models.Item
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	err = json.Unmarshal(body, &ItemArr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	if err := config.DB.Model(&models.Cart{}).Where(" user_id =? ", getId).Update("is_purchased", true).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
		return
	}
	for _, item := range ItemArr {
		if status := isItemExist(item.ID); status {
			var cartItem models.CartItem
			cartItem.ItemID = item.ID
			cartIdDer, cartId := getCart(getId.(string))
			fmt.Println("cartId    :  ", cartId.ID)
			if cartIdDer {
				cartItem.CartID = cartId.ID
				PrettyPrint(cartItem)
				if err := config.DB.Create(cartItem).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | unable to add items in cart ", "error": err})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | unable to add items in cart as cart id not found ", "error": err})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "items added to cart successfully"})
	return
}

//CompleteCart ...
func CompleteCart(c *gin.Context) {
	cartId := c.Param("id")
	userId := c.MustGet("user_id")
	timeStamp := time.Now().UnixNano()

	fmt.Println(" cartId ", cartId, " userId", userId)
	status, userInput := getCart(cartId)
	if status && strconv.Itoa(int(userInput.UserID)) == userId {
		if err := config.DB.Model(&models.Cart{}).Where("id = ? and user_id =? ", cartId, userId).Update("is_purchased", false).Error; err != nil {
			fmt.Println(err)
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
			return
		}
		// if err := config.DB.Where("cart_id = ?", cartId).Delete(&models.CartItem{}).Error; err != nil {
		// 	fmt.Println(err)
		// 	if gorm.IsRecordNotFoundError(err) {

		// 	}
		// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable remove items form cart ", "error": err})
		// 	return
		// }
		var cart models.Cart
		cart.CreatedAt = timeStamp / int64(time.Millisecond)
		cart.ID = timeStamp
		cart.UserID = userInput.UserID
		cart.ISPurchased = true
		
		if err := config.DB.Create(cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | please check the request body ", "error": err})
			return
		}
		if err := config.DB.Model(&models.User{}).Where("id = ?", userId).Update("cart_id", timeStamp).Error; err != nil {
			fmt.Println(err)
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
			return
		}
		var orders models.Orders
		orders.ID = timeStamp
		orders.UserID = userInput.UserID
		orders.CartID = userInput.ID
		orders.CreatedAt = timeStamp / int64(time.Millisecond)
		if err := config.DB.Create(orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | unable to place the order ", "error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "cart checkout completed"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | unable to place the order"})
	return
}

func isItemExist(itemID int64) bool {
	var Item models.Item
	if err := config.DB.Raw("SELECT name , created_at ,id FROM items WHERE id= ?", itemID).Scan(&Item).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			return false
		}
		return false
	}
	return true
}
func getCart(cartId string) (bool, models.Cart) {
	var carts models.Cart
	fmt.Println(cartId)
	if err := config.DB.Raw("SELECT * FROM carts WHERE user_id= ?", cartId).Scan(&carts).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			return false, carts
		}
		return false, carts
	}
	return true, carts
}
