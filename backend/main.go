package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lcox74/discord-bot/actions"
	"github.com/lcox74/discord-bot/bot"
	cors "github.com/rs/cors/wrapper/gin"
)

var token = os.Getenv("DISCORD_TOKEN")

var discordBot *bot.Bot

func main() {
	var err error

	fmt.Println(token)
	
	// Create a new Discord session using the provided bot token.
	discordBot, err = bot.CreateBot(token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	defer discordBot.Stop()

	go discordBot.Run()

	router := gin.Default()

	router.Use(cors.Default())

    router.POST("/news", postNews)
    router.GET("/news/:auth", getNewsAuth)

    router.Run("0.0.0.0:3000")

}


func postNews(c *gin.Context) {
	var newsPost actions.NewsAction

	if err := c.BindJSON(&newsPost); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check Auth
	if !discordBot.GetUser(bot.SCOPE_NEWS, newsPost.AccessToken) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	
	// Fetch User Data
	accessRecord, err := discordBot.Db.ValidateAccessKey(newsPost.AccessToken)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Add User Data
	newsPost.User = accessRecord.UserName
	discordBot.AddAction(newsPost)

	fmt.Println(newsPost)
	c.IndentedJSON(http.StatusCreated, newsPost)
}

func getNewsAuth(c *gin.Context) {
	access_key := c.Param("auth")
	if access_key == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !discordBot.GetUser(bot.SCOPE_NEWS, access_key) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"access": access_key})
}