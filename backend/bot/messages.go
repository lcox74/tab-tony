package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var accessMessage = func(access AccessRecord) discordgo.MessageEmbed {

	// Format the Services
	services := []string{}
	if access.IsNewsScope() {
		services = append(services, "- News")
	}
	if access.IsZerotierScope() {
		services = append(services, "- ZeroTier")
	}
	servicesString := strings.Join(services, "\n")

	return discordgo.MessageEmbed{
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
}	

var zerotierJoinMessage = func(nwid string) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title: "TAB ZerotTier",
		Description: fmt.Sprintf("You have been authorised on the TAB internal network. Please join the network on your authorised machine by doing the following in your Terminal:\n```bash\nsudo zerotier-cli join %s\n```\nThis will allow you to be on a SD-WAN connection with other people connected to the TAB network.", nwid),
	}
}

