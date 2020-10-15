package controller

import (
	"SezzleTest/config"
	"SezzleTest/middleware"
	"SezzleTest/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//CreateUsers ...
func CreateUsers(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	var userInput models.User
	err = json.Unmarshal(body, &userInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}

	if userInput.Name == "" || userInput.UserName == "" || userInput.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | name , username , password are mandatory"})
		return
	}
	userInput.Password = middleware.BcryptPassword(userInput.Password)
	timeStamp := time.Now().UnixNano()
	userInput.CartID = timeStamp
	userInput.ID = timeStamp
	userInput.CreatedAt = timeStamp / int64(time.Millisecond)
	token := middleware.GenerateToken(userInput.Password, userInput.Name, userInput.UserName, timeStamp, true)
	userInput.Token = token
	PrettyPrint(userInput)
	if err := config.DB.Create(userInput).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | please check the request body ", "error": err})
		return
	}
	var cart models.Cart
	cart.CreatedAt = timeStamp / int64(time.Millisecond)
	cart.ID = userInput.CartID
	cart.UserID = userInput.ID
	cart.ISPurchased = false
	if err := config.DB.Create(cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | please check the request body ", "error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Created Successfully .", "token": token})
	return
}
func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
func VerifyBcrypt(db_password, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(db_password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func UserLogin(c *gin.Context) {

	var Data struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}
	err = json.Unmarshal(body, &Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request | please check the request body"})
		return
	}

	if Data.UserName == "" || Data.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request |  username , password are mandatory"})
		return
	}
	var dbUser models.User
	if err := config.DB.Where("user_name = ?", Data.UserName).Find(&dbUser).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			// handle object not found
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | please check the request body ", "error": err})
		return
	}
	status := middleware.VerifyBcrypt(dbUser.Password, Data.Password)
	if status {
		token := middleware.GenerateToken(dbUser.Password, dbUser.Name, dbUser.UserName, dbUser.ID, true)
		fmt.Println(status)
		if err := config.DB.Model(&models.User{}).Where("user_name = ?", Data.UserName).Update("token", token).Error; err != nil {
			fmt.Println(err)
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User loggedin Successfully .", "token": token})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "incorrect password."})
	return
}

//UserList ...
func UserList(c *gin.Context) {
	var usersArr []models.User
	if err := config.DB.Raw("SELECT name , user_name , cart_id , created_at ,id FROM users ").Scan(&usersArr).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error | enable to store the token ", "error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user list fetched Successfully", "users": usersArr})
	return
}
