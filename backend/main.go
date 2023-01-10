package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lcox74/discord-bot/actions"
	"github.com/lcox74/discord-bot/bot"
	"github.com/lcox74/discord-bot/zerotier"
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
	router.GET("/zerotier/:auth", getZeroTierAuth)
	router.GET("/zerotier/:auth/network", getZeroTierNetwork)

	router.Run("0.0.0.0:3000")

}

func postNews(c *gin.Context) {
	var newsPost actions.NewsAction

	if err := c.BindJSON(&newsPost); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check Auth
	if !discordBot.GetUserFromKey(bot.SCOPE_NEWS, newsPost.AccessToken) {
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

	if !discordBot.GetUserFromKey(bot.SCOPE_NEWS, access_key) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"access": access_key})
}

func getZeroTierAuth(c *gin.Context) {
	access_key := c.Param("auth")
	if access_key == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !discordBot.GetUserFromKey(bot.SCOPE_ZEROTIER, access_key) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"access": access_key})
}

type ztNode struct {
	Name     string `json:"name"`
	IP       string `json:"ip"`
	IsActive bool   `json:"active"`
	Image    string `json:"image"`
}

func getZeroTierNetwork(c *gin.Context) {
	access_key := c.Param("auth")
	if access_key == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !discordBot.GetUserFromKey(bot.SCOPE_ZEROTIER, access_key) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get ZT Network Id
	nwid := os.Getenv("ZEROTIER_GENERAL_NET_ID")

	// Get ZT Network Members
	members, err := zerotier.GetNetworkMembers(nwid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := make(map[string][]ztNode)

	for _, member := range members {
		// Get Node Member Data
		d, err := discordBot.Db.GetMember(nwid, member.MemberId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Get User Access to fetch Username
		acc, err := discordBot.Db.GetUserAccess(d.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		node := ztNode{
			Name:     d.MemberName,
			IP:       member.Config.AssignedIps[0],
			IsActive: member.IsOnline,
			Image:    acc.UserAvatarUrl,
		}

		// Create User Array (If Not Exists)
		if _, ok := response[acc.UserName]; !ok {
			response[acc.UserName] = make([]ztNode, 0)
		}
		response[acc.UserName] = append(response[acc.UserName], node)
	}
	c.IndentedJSON(http.StatusOK, response)
}
