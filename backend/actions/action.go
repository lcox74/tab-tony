package actions

import "github.com/bwmarrin/discordgo"

type Action interface {
	Execute(s *discordgo.Session) error
}