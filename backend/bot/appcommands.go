package bot

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/lcox74/discord-bot/zerotier"
)

// ==================
// 		Commands
// ==================

var ApplicationCmdDelcares = []*discordgo.ApplicationCommand{
	{
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
	},
	{
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
				Required:    true,
			},
		},
	},
}


// ==================
// 		Handlers
// ==================

func RequestCmd (bot *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) {

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
		message := accessMessage(access)
		bot.sendEmbededDM(user, &message)

		// Send the response
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("I've just updated your access, please check your DMs <@%s>", user.ID),
			},
		})
	}
}

func JoinNetCmd (bot *Bot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("Join Net")

	// Get the user
	var user *discordgo.User = getuser(i)
	if user == nil {
		return
	}

	// Check access
	if !bot.GetUserFromId(SCOPE_ZEROTIER, user.ID) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("You don't have access to the Zerotier network, please request access first <@%s>", user.ID),
			},
		})
		return
	}

	
	fmt.Printf("User: %s\n", user.Username)
	var addr, nickname string

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options
	
	for _, option := range options {
		switch option.Name {
		case "address":
			fmt.Printf("addr: %s\n", option.Value)
			addr = option.Value.(string)
		case "nickname":
			fmt.Printf("nickname: %s\n", option.Value)				
			nickname = option.Value.(string)
		}
	}

	nwid := os.Getenv("ZEROTIER_GENERAL_NET_ID")

	// Create a ZeroTier Network Member
	err := zerotier.AuthoriseMember(nwid, addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Update the member name
	bot.Db.AddMember(nwid, addr, nickname, user.ID)

	// Send the access
	message := zerotierJoinMessage(nwid)
	bot.sendEmbededDM(user, &message)

	// Send the response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("I've just added device %s to the network, please check your DMs <@%s>", nickname, user.ID),
		},
	})
}