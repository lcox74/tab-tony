package bot

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lcox74/discord-bot/actions"
	"github.com/lcox74/discord-bot/requests"
)

type actionQueue struct {
	messages []actions.Action
	lock     sync.RWMutex
}

type Bot struct {
	discord *discordgo.Session
	isReady bool

	queue     actionQueue
	newsUsers map[string]string
}

var generalChannelID = os.Getenv("GENERAL_CHANNEL_ID")

func CreateBot(token string) (*Bot, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		discord:   discord,
		isReady:   false,
		newsUsers: make(map[string]string),
	}

	// Get original Access
	accessFile, err := ioutil.ReadFile("access.json")
	if err == nil {

		var access map[string]string
		err = json.Unmarshal(accessFile, &access)
		if err != nil {
			return nil, err
		}

		bot.newsUsers = access
	}

	// Regiser a handler for the ready event
	discord.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		// Set the playing status.
		s.UpdateGameStatus(0, "Idle Simulator")
		s.ChannelMessageSend(generalChannelID, "Look out look out here's Trevor the Trout!")

		// Set bot to ready
		bot.isReady = true

		// Log ready
		fmt.Printf("Bot is ready! (User: %s)\n", event.User.Username)
	})

	discord.AddHandler(bot.onMessage)

	return bot, nil
}

func (bot *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "/request news" {
		fmt.Printf("Requesting news for %s\n", m.Author.Username)
		password, ok := bot.newsUsers[m.Author.Username]
		if !ok {
			password = randomPassword() // TODO: Check Collisions
			bot.newsUsers[m.Author.Username] = password

			jsonByte, _ := json.Marshal(bot.newsUsers)
			err := ioutil.WriteFile("access.json", jsonByte, 0644)
			if err != nil {
				fmt.Println("Error writing to access.json:", err)
			}
		}
		requests.NewsRequest(s, m, password)
	}
}

func (bot *Bot) AddAction(action actions.Action) {
	bot.queue.lock.Lock()
	defer bot.queue.lock.Unlock()

	bot.queue.messages = append(bot.queue.messages, action)
}

func (bot *Bot) Start() error {
	// Open the websocket and begin listening.
	err := bot.discord.Open()
	if err != nil {
		return err
	}
	return nil
}

func (bot *Bot) Stop() error {
	// Cleanly close down the Discord session.
	bot.discord.Close()
	return nil
}

// Run the bot, make sure to run as a goroutine
func (bot *Bot) Run() {
	for {
		if bot.isReady {
			bot.queue.lock.Lock()
			for _, action := range bot.queue.messages {
				action.Execute(bot.discord)
			}
			bot.queue.messages = []actions.Action{}
			bot.queue.lock.Unlock()
		} else {
			fmt.Println("Bot is not ready yet, will try again in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
}

func (bot *Bot) GetUser(access string) string {

	for id, password := range bot.newsUsers {
		if password == access {
			return id
		}
	}
	return ""
}

func randomPassword() string {
	b := make([]byte, 32)
	rand.Read(b)

	// Calculate the SHA1 hash of the message
	h := sha1.New()
	h.Write(b)
	hash := h.Sum(nil)

	// Encode the hash to base64
	return base64.StdEncoding.EncodeToString(hash)

}
