package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Message struct { 
	Text string `json:"text"`
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to find .env file")
	}

	r := gin.New()
	r.Use(CORSMiddleware())

	r.POST("/", func(c *gin.Context) {

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		var message Message
		err = json.Unmarshal(body, &message)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		fmt.Println(message.Text)
		
		// var jsonData Message

		// if err := c.ShouldBindJSON(&jsonData); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }

		// jsonDataBytes, err := json.MarshalIndent(jsonData, "", "  ")
		// if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON data"})
		// 		return
		// }
		
		// json := string(jsonDataBytes)
		// fmt.Println(jsonDataBytes)
		// fmt.Println(json)

		// c.JSON(http.StatusOK, gin.H{
		// 	"message": "hello world",
		// })

	})

	r.Run()

}