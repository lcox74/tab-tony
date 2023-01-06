package actions

import (
	"fmt"
	"os"

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

// Target Channel for News
var newsChannelID = os.Getenv("NEWS_CHANNEL_ID")

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

	_, err := s.ChannelMessageSendEmbed(newsChannelID, &message)

	return err
}

func joinTags(tags []string) string {
	var result string
	for _, tag := range tags {
		result += fmt.Sprintf("`%s` ", tag)
	}
	return result
}
