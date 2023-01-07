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
	
	// Create a new Discord session using the provided bot token.
	discordBot, err = bot.CreateBot(token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Open the websocket and begin listening.
	discordBot.Start()
	go discordBot.Run()

	router := gin.Default()

	router.Use(cors.Default())

    router.POST("/news", postNews)
    router.GET("/news/:auth", getNewsAuth)

    router.Run("0.0.0.0:3000")

	discordBot.Stop()
}


func postNews(c *gin.Context) {
	var newsPost actions.NewsAction

	if err := c.BindJSON(&newsPost); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check Auth
	newsPost.User = discordBot.GetUser(newsPost.AccessToken)
	if newsPost.User == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	discordBot.AddAction(newsPost)

	fmt.Println(newsPost)
	c.IndentedJSON(http.StatusCreated, newsPost)
}

func getNewsAuth(c *gin.Context) {
	user := c.Param("auth")
	if user == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if discordBot.GetUser(user) == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"access": user})
}