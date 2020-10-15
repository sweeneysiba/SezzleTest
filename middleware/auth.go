package middleware

import (
	"SezzleTest/config"
	"SezzleTest/models"
	"fmt"
	"log"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var mySigningKey = []byte("1qaz2wsx3edc4rfv5tgb6yhn7ujm8ik,9ol.0p;/")

//IsAuthorizedApp ...
func IsAuthorizedApp() gin.HandlerFunc {
	return func(c *gin.Context) {

		authtoken := c.Request.Header.Get("Authorization")
		if authtoken != "" {

			token, err := jwt.Parse(authtoken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				fmt.Println(err)
				c.JSON(401, gin.H{
					"status":  "unauthorized",
					"message": "Failed Authentication",
				})
				c.AbortWithStatus(401)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fmt.Println(claims["id"])
				sameToken := isCheckToken(claims["user_id"].(string), authtoken)
				if sameToken {
					c.Set("user_id", claims["user_id"].(string))
					c.Next()
				}
			} else {
				c.JSON(401, gin.H{
					"status":  "unauthorized",
					"message": "Failed Authentication",
				})
				c.AbortWithStatus(401)
			}
			return
		} else {
			c.JSON(401, gin.H{
				"status":  "unauthorized",
				"message": "Failed Authentication",
			})
			c.AbortWithStatus(401)
			return
		}
	}
}

func GenerateToken(password, name, username string, userid int64, expiry bool) string {
	fmt.Println("user id ", userid)
	token := jwt.New(jwt.SigningMethodHS256)
	if expiry {
		token.Claims = jwt.MapClaims{
			"name":     name,
			"password": password,
			"username": username,
			"user_id":  strconv.Itoa(int(userid)),
		}
	} else {
		token.Claims = jwt.MapClaims{
			"name":     name,
			"password": password,
			"username": username,
			"user_id":  strconv.Itoa(int(userid)),
		}
	}
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}

func isCheckToken(userId, token string) bool {
	var usersArr models.User
	if err := config.DB.Raw("SELECT token FROM users where id=? ", userId).Scan(&usersArr).Error; err != nil {
		fmt.Println(err)
		if gorm.IsRecordNotFoundError(err) {
			fmt.Println("no records found")
			return false
		}
		fmt.Println("error ", err)
		return false
	}
	if usersArr.Token == token {
		fmt.Println("token match")
		return true
	}
	fmt.Println("token mismatch")
	return false
}
