package bot

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lcox74/discord-bot/actions"
)

type actionQueue struct {
	messages []actions.Action
	lock     sync.RWMutex
}

type Bot struct {
	discord *discordgo.Session
	isReady bool

	queue  actionQueue
	Db BotDatabase
}

// Handler functions for the Application Commands
var commandHandlers = map[string]func (bot *Bot, s *discordgo.Session, i *discordgo.InteractionCreate){
	"request": RequestCmd,
	"join-net": JoinNetCmd,
}


func CreateBot(token string) (*Bot, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		discord: discord,
		isReady: false,
	}

	// Open database
	bot.Db, err = OpenAccessDb()
	if err != nil {
		return nil, err
	}

	// Regiser a handler for the ready event
	discord.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		// Set the playing status.
		s.UpdateGameStatus(0, "Idle Simulator")

		// Set bot to ready
		bot.isReady = true

		// Log ready
		fmt.Printf("Bot is ready! (User: %s)\n", event.User.Username)
	})

	// Register a handler for handling Application Commands
	discord.AddHandler(func (s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(bot, s, i)
		}
	})

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Register all the application commands
	for _, cmd := range ApplicationCmdDelcares {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", cmd)
		if err != nil {
			fmt.Println(err)
		}
	}

	return bot, nil
}

// Add an action to the queue
func (bot *Bot) AddAction(action actions.Action) {
	bot.queue.lock.Lock()
	defer bot.queue.lock.Unlock()

	bot.queue.messages = append(bot.queue.messages, action)
}


func (bot *Bot) Stop() error {
	// Cleanly close down the Discord session.
	bot.discord.Close()
	bot.Db.Close()
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

func (bot *Bot) GetUserFromKey(scope Scope, access string) bool {

	record, err := bot.Db.ValidateAccessKey(access)
	if err != nil {
		fmt.Println(err)
		return false
	}

	switch scope {
	case SCOPE_NEWS:
		return record.IsNewsScope()
	case SCOPE_ZEROTIER:
		return record.IsZerotierScope()
	}

	return false
}

func (bot *Bot) GetUserFromId(scope Scope, userId string) bool {

	record, err := bot.Db.GetUserAccess(userId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	switch scope {
	case SCOPE_NEWS:
		return record.IsNewsScope()
	case SCOPE_ZEROTIER:
		return record.IsZerotierScope()
	}

	return false
}

// ======================
//    Helper Functions
// ======================

func (bot *Bot) sendEmbededDM(user *discordgo.User, embed *discordgo.MessageEmbed) error {
	// Get the user's DM channel
	channel, err := bot.discord.UserChannelCreate(user.ID)
	if err != nil {
		return err
	}

	// Send the message
	_, err = bot.discord.ChannelMessageSendEmbed(channel.ID, embed)
	if err != nil {
		return err
	}

	return nil
}

func getuser(i *discordgo.InteractionCreate) *discordgo.User {
	if i.Member == nil {
		return i.User
	} 
	return i.Member.User		
}