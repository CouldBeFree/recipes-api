// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//	Schemes: http
//  Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Mohamed Labouardy <mohamed@labouardy.com> https://labouardy.com
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta

package main

import (
	"context"
	"fmt"
	"log"

	handlers "github.com/CouldBeFree/recipes-api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// var recipes []Recipe
var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

// func UpdateRecipeHandler(c *gin.Context) {
// 	id := c.Param("id")
// 	var recipe Recipe
// 	if err := c.ShouldBindJSON(&recipe); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error()})
// 		return
// 	}
// 	objectId, _ := primitive.ObjectIDFromHex(id)
// 	_, err = collection.UpdateOne(ctx, bson.M{
// 		"_id": objectId,
// 	}, bson.D{{"$set", bson.D{
// 		{"name", recipe.Name},
// 		{"instructions", recipe.Instructions},
// 		{"ingredients", recipe.Ingredients},
// 		{"tags", recipe.Tags},
// 	}}})
// 	if err != nil {
// 		fmt.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
// }

// func DeleteRecipeHandler(c *gin.Context) {
// 	id := c.Param("id")
// 	index := -1
// 	for i := 0; i < len(recipes); i++ {
// 		if recipes[i].ID == id {
// 			index = i
// 		}
// 	}
// 	if index == -1 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Recipe not found"})
// 		return
// 	}
// 	recipes = append(recipes[:index], recipes[index+1:]...)
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Recipe has been deleted"})
// }

// func SearchRecipesHandler(c *gin.Context) {
// 	tag := c.Query("tag")
// 	listOfRecipes := make([]Recipe, 0)
// 	for i := 0; i < len(recipes); i++ {
// 		found := false
// 		for _, t := range recipes[i].Tags {
// 			if strings.EqualFold(t, tag) {
// 				found = true
// 			}
// 		}
// 		if found {
// 			listOfRecipes = append(listOfRecipes, recipes[i])
// 		}
// 	}
// 	c.JSON(http.StatusOK, listOfRecipes)
// }

var recipesHandler *handlers.RecipesHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping()
	fmt.Println(status)
	collection := client.Database("demo").Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	//router.GET("/recipes/search", SearchRecipesHandler)
	router.Run(":5000")
}
