package main

import (
	"fmt"
	"log"
	"net/http"
	f "online-shop/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()    

	//Product methods
	router.POST("/product", f.CreateProduct) 
	router.GET("/product/:id", f.GetProductByID)  
	router.PUT("/product", f.UpdateProductByID) 
	router.DELETE("/product/:id", f.DeleteCart)
	router.GET("/product/all", f.ListProduct)


	//User methods
	router.POST("/users", f.CreateUser)
	router.GET("/users/:id", f.GetUserByID)
	router.PUT("/users", f.UpdateUserID)
	router.DELETE("/users/:id", f.DeleteUser)
	router.GET("/users/all", f.ListUser)


	//Cart methods
	router.POST("/cart", f.CreateCart)
    router.GET("/cart/:id", f.GetCartByID)
    router.DELETE("/cart/:id", f.DeleteCart)
    router.GET("/cart/all", f.Listcart)


	


	fmt.Println("Listening :8080...")

	err := http.ListenAndServe(":8080", router) 
    if err != nil {
        log.Println("Error starting server:", err)
    }
}