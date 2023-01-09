package bot

import (
	"fmt"
	"strings"
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

var commandHandlers = map[string]func (bot *Bot, s *discordgo.Session, i *discordgo.InteractionCreate){
	"request": func (bot *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) {

		// Access options in the order provided by the user.
		options := i.ApplicationCommandData().Options

		// Go through the services
		for _, option := range options {
			var scope Scope
			switch option.Value {
			case "news":
				scope = SCOPE_NEWS
			case "zerotier":
				scope = SCOPE_ZEROTIER
			}

			// Get the user
			var user *discordgo.User = getuser(i)
			if user == nil {
				return
			}

			// Add the access
			access, err := bot.Db.AddAccess(user.ID, user.Username, scope)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Send the access
			bot.SendAccessDM(user, access)

			// Send the response
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("I've just updated your access, please check your DMs <@%s>", user.ID),
				},
			})
		}
	},
	"join-net": func (bot *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) {
		fmt.Println("Join Net")

		// Get the user
		var user *discordgo.User = getuser(i)
		if user == nil {
			return
		}
		fmt.Printf("User: %s\n", user.Username)

		// Access options in the order provided by the user.
		options := i.ApplicationCommandData().Options

		for _, option := range options {
			switch option.Name {
			case "address":
				fmt.Printf("addr: %s\n", option.Value)
			case "nickname":
				fmt.Printf("nickname: %s\n", option.Value)				
			}
		}
	},
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


	// Register an App command for the `request` command
	discord.ApplicationCommandCreate(discord.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "request",
		Description: "Request something access for a service",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "service",
				Description: "The service you want to request access for",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "News",
						Value: "news",
					},
					{
						Name:  "Zerotier",
						Value: "zerotier",
					},
				},
			},
		},
	})

	// Register an App command for the `join-net` command
	discord.ApplicationCommandCreate(discord.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "join-net",
		Description: "Add a device to the TAB Zerotier network",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "address",
				Description: "ZeroTier address eg. `34994c713f`",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "nickname",
				Description: "Nickname for the device eg. `Tony's Desktop`",
				Required:    false,
			},
		},
	})


	

	return bot, nil
}

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

func (bot *Bot) GetUser(scope Scope, access string) bool {

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

// Send Access DM
func (bot *Bot) SendAccessDM(user *discordgo.User, access AccessRecord) error {
	// Get the user's DM channel
	channel, err := bot.discord.UserChannelCreate(user.ID)
	if err != nil {
		return err
	}

	// Format the Services
	services := []string{}
	if access.IsNewsScope() {
		services = append(services, "- News")
	}
	if access.IsZerotierScope() {
		services = append(services, "- ZeroTier")
	}
	servicesString := strings.Join(services, "\n")

	// Create the message
	message := discordgo.MessageEmbed{
		Title: "TAB Access",
		Description: "You have requested access to the TAB Discord Server. The following is your access key which is needed for accessing the TAB services which are also listed.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Key",
				Value: fmt.Sprintf("```\n%s\n```", access.AccessKey),
			},
			{
				Name: "Services",
				Value: fmt.Sprintf("```\n%s\n```", servicesString),
			},
		},
	}

	// Send the message
	_, err = bot.discord.ChannelMessageSendEmbed(channel.ID, &message)
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