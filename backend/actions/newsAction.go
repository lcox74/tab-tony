package actions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type NewsAction struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	SourceURL   string   `json:"url"`
	ImageURL    string   `json:"image_url"`
	Tags        []string `json:"tags"`
	AccessToken string   `json:"access"`
	User        string   `json:"-"`
}

func (action NewsAction) Execute(s *discordgo.Session) error {

	message := discordgo.MessageEmbed{
		Title:       action.Title,
		Description: action.Content,
		URL:         action.SourceURL,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   joinTags(action.Tags),
				Value:  "\u200B",
				Inline: true,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: action.ImageURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Added by @%s", action.User),
		},
	}

	if action.ImageURL == "" {
		message.Image = nil
	}
	if len(action.Tags) == 0 {
		message.Fields = nil
	}

	for _, server := range s.State.Guilds {
		channels, _ := s.GuildChannels(server.ID)
		for _, c := range channels {
			if c.Name == "tech-news" {
				_, err := s.ChannelMessageSendEmbed(c.ID, &message)
				if err != nil {
					fmt.Println(err)
				}
			}	
		}
	}


	return nil
}

func joinTags(tags []string) string {
	var result string
	for _, tag := range tags {
		result += fmt.Sprintf("`%s` ", tag)
	}
	return result
}
