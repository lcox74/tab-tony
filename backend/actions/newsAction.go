package actions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// NewsAction represents an action for posting news to a Discord channel
type NewsAction struct {
	Title       string   `json:"title"`     // title of the news article
	Content     string   `json:"content"`   // content of the news article
	SourceURL   string   `json:"url"`       // source URL of the news article
	ImageURL    string   `json:"image_url"` // image URL for the news article [Optional]
	Tags        []string `json:"tags"`      // tags for the news article [Optional]
	AccessToken string   `json:"access"`    // access token for the user who added the news article
	User        string   `json:"-"`         // Discord username of the user who added the news article
}

func (action NewsAction) IsValid() bool {
	return action.Title != "" && action.Content != "" && action.SourceURL != "" && action.AccessToken != ""
}

// Execute sends the news article as an embedded message to the "tech-news" channel on all servers that the bot is a member of
func (action NewsAction) Execute(s *discordgo.Session) error {
	if !action.IsValid() {
		return fmt.Errorf("invalid news action")
	}
	
	// Create an embedded message with the news article information
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

	// Setting Defaults if optionals arent provided
	if action.ImageURL == "" { message.Image = nil }
	if len(action.Tags) == 0 { message.Fields = nil }

	// Iterate through all servers that the bot is a member of
	for _, server := range s.State.Guilds {
		// Iterate through the list of channels in the server
		channels, _ := s.GuildChannels(server.ID)
		for _, c := range channels {
			// If the channel is named "tech-news", send the embedded message to the channel
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

// joinTags joins a slice of tags into a single string
func joinTags(tags []string) string {
	var result string
	for _, tag := range tags {
		result += fmt.Sprintf("`%s` ", tag)
	}
	return result
}
