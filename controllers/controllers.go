package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Price float64            `json:"price,omitempty" bson:"price,omitempty"`
}

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName string             `json:"fullname,omitempty" bson:"fullname,omitempty"`
	Age      int64              `json:"age,omitempty" bson:"age,omitempty"`
	Phone    string             `json:"phone,omitempty" bson:"phone,omitempty"`
}

type CartReq struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userid" bson:"user,omitempty"`
	ProductID string             `json:"productid" bson:"product,omitempty"`
}

type CartRes struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User    *User              `json:"user" bson:"user,omitempty"`
	Product *Product           `json:"product" bson:"product,omitempty"`
}

var client *mongo.Client
var dbCollProd, dbCollUser, dbCollCart *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB...")

	dbCollProd = client.Database("shop").Collection("product")
	dbCollUser = client.Database("shop").Collection("user")
	dbCollCart = client.Database("shop").Collection("cart")

}

func CreateProduct(ctx *gin.Context) {

	var body, bodyRes Product

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to bind json", err)
		return
	}

	result, err := dbCollProd.InsertOne(ctx, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("Failed to create product ", err)
		return
	}

	err = dbCollProd.FindOne(ctx, bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&bodyRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to get product", err)
		return
	}

	ctx.JSON(http.StatusCreated, bodyRes)
}

func GetProductByID(ctx *gin.Context) {

	var body Product

	id := ctx.Param("id")

	err := dbCollProd.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to get product", err)
		return
	}

	ctx.JSON(http.StatusOK, body)
}

func UpdateProductByID(ctx *gin.Context) {

	var body, bodyRes Product

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to bind json", err)
		return
	}

	updateReq := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: body.Name},
			{Key: "price", Value: body.Price},
		},
		},
	}

	err = dbCollProd.FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: body.ID}}, updateReq).Decode(&bodyRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to update product", err)
		return
	}

	ctx.JSON(http.StatusOK, updateReq)
}

func DeleteProductByID(ctx *gin.Context) {
	id := ctx.Param("_id")

	resp, err := dbCollProd.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to delete product", err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func ListProduct(ctx *gin.Context) {
	var products []primitive.M

	cursor, err := dbCollProd.Find(ctx, bson.D{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to Get all  product", err)
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product primitive.M

		err = cursor.Decode(&product)
		if err != nil {
			log.Println(err)
			return
		}

		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {

		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

//User CRUD

func CreateUser(ctx *gin.Context) {

	var body, bodyRes User

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to bind json", err)
		return
	}

	result, err := dbCollUser.InsertOne(ctx, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("Failed to create user ", err)
		return
	}

	err = dbCollUser.FindOne(ctx, bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&bodyRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to get user", err)
		return
	}

	ctx.JSON(http.StatusCreated, bodyRes)
}

func GetUserByID(ctx *gin.Context) {

	id := ctx.Param("id")

	var body User
	err := dbCollUser.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to get user", err)
		return
	}

	ctx.JSON(http.StatusOK, body)
}

func UpdateUserID(ctx *gin.Context) {

	var body, bodyRes User

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to bind json", err)
		return
	}

	updateReq := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "fullname", Value: body.FullName},
			{Key: "age", Value: body.Age},
			{Key: "phone", Value: body.Phone},
		},
		},
	}

	err = dbCollUser.FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: body.ID}}, updateReq).Decode(&bodyRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to update user", err)
		return
	}

	ctx.JSON(http.StatusOK, updateReq)
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("_id")

	resp, err := dbCollUser.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to delete user", err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func ListUser(ctx *gin.Context) {
	var users []primitive.M

	cursor, err := dbCollUser.Find(ctx, bson.D{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to Get all  users", err)
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user primitive.M

		err = cursor.Decode(&user)
		if err != nil {
			log.Println(err)
			return
		}

		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {

		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// Cart CRUD
func CreateCart(ctx *gin.Context) {

	var body CartReq
	var bodyRes CartRes
	var userRes User
	var prodRes Product

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to bind json", err)
		return
	}

	_, err = dbCollCart.InsertOne(ctx, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("Failed to create cart ", err)
		return
	}

	// err = dbCollCart.FindOne(ctx, bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(&bodyRes)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	fmt.Println("failed to create cart", err)
	// 	return
	// }

	err = dbCollProd.FindOne(ctx, bson.D{{Key: "_id", Value: body.ProductID}}).Decode(&prodRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to get product", err)
		return
	}

	err = dbCollUser.FindOne(ctx, bson.D{{Key: "_id", Value: body.UserID}}).Decode(&userRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to get user", err)
		return
	}

	bodyRes.User = &userRes
	bodyRes.Product = &prodRes
	bodyRes.ID = body.ID

	ctx.JSON(http.StatusCreated, bodyRes)
}

func GetCartByID(ctx *gin.Context) {

	var body CartReq
	var bodyRes CartRes
	var userRes User
	var prodRes Product

	id := ctx.Param("id")

	err := dbCollCart.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to get cart", err)
		return
	}

	err = dbCollProd.FindOne(ctx, bson.D{{Key: "_id", Value: body.ProductID}}).Decode(&prodRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to get product", err)
		return
	}

	err = dbCollUser.FindOne(ctx, bson.D{{Key: "_id", Value: body.UserID}}).Decode(&userRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to get user", err)
		return
	}

	bodyRes.User = &userRes
	bodyRes.Product = &prodRes
	bodyRes.ID = body.ID

	ctx.JSON(http.StatusCreated, bodyRes)
}

func DeleteCart(ctx *gin.Context) {
	id := ctx.Param("_id")

	resp, err := dbCollCart.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to delete cart", err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func Listcart(ctx *gin.Context) {
	var carts []primitive.M

	cursor, err := dbCollCart.Find(ctx, bson.D{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("failed to Get all  carts", err)
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var cart primitive.M

		err = cursor.Decode(&cart)
		if err != nil {
			log.Println(err)
			return
		}

		carts = append(carts, cart)
	}

	if err := cursor.Err(); err != nil {

		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, carts)
}
