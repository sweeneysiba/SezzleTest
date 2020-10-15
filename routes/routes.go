package routes

import (
	"SezzleTest/controller"
	middleware "SezzleTest/middleware"

	"github.com/gin-gonic/gin"
)

//SetupRouter ... Configure routes
func SetupRouter() *gin.Engine {
	r := gin.Default()
	grp1 := r.Group("/user")
	{
		grp1.POST("create", func(c *gin.Context) {
			controller.CreateUsers(c)
		})
		grp1.POST("login", func(c *gin.Context) {
			controller.UserLogin(c)
		})
		grp1.GET("list", func(c *gin.Context) {
			controller.UserList(c)
		})
	}

	grp3 := r.Group("/cart")
	{
		grp3.POST("list", func(c *gin.Context) {
			controller.ListCart(c)
		})

	}
	grp2 := r.Group("/item")
	{
		grp2.GET("list", func(c *gin.Context) {
			controller.ListItems(c)
		})
		grp2.POST("create", func(c *gin.Context) {
			controller.CreateItem(c)
		})
	}
	grp4 := r.Group("/order")
	{
		grp4.GET("list", func(c *gin.Context) {
			controller.ListOrders(c)
		})
	}
	grp5 := r.Group("/cart")
	grp5.Use(middleware.IsAuthorizedApp())
	{
		grp5.POST("add", func(c *gin.Context) {
			controller.AddToCart(c)
		})

		grp5.GET("/:id/complete", func(c *gin.Context) {
			controller.CompleteCart(c)
		})
	}
	return r
}
