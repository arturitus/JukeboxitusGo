package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/json"
	"github.com/disgoorg/log"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "play",
		Description: "Plays a song",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "identifier",
				Description: "The song link or search query",
				Required:    true,
			},
		},
	},
	{
		Name:        "pause",
		Description: "Pauses the current song",
	},
	{
		Name:        "skip",
		Description: "Skips the current song",
	},
	{
		Name:        "now-playing",
		Description: "Shows the current playing song",
	},
	{
		Name:        "stop",
		Description: "Stops the current song and stops the player",
	},
	{
		Name:        "players",
		Description: "Shows all active players",
	},
	{
		Name:        "shuffle",
		Description: "Shuffles the current queue",
	},
	{
		Name:        "queue",
		Description: "Shows the current queue",
	},
	{
		Name:        "clear-queue",
		Description: "Clears the current queue",
	},
	{
		Name:        "queue-type",
		Description: "Sets the queue type",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "type",
				Description: "The queue type",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "default",
						Value: "default",
					},
					{
						Name:  "repeat-track",
						Value: "repeat-track",
					},
					{
						Name:  "repeat-queue",
						Value: "repeat-queue",
					},
				},
			},
		},
	},
	{
		Name:        "volume",
		Description: "Sets the player volume",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "level",
				Description: "Volume level (0-100)",
				Required:    true,
				MinValue:    json.Ptr(0.0),
				MaxValue:    100,
			},
		},
	},
	{
		Name:        "bass-boost",
		Description: "Toggles bass boost filter",
	},
	{
		Name:        "eight-d",
		Description: "Toggles 8-D audio filter",
	},
	{
		Name:        "lyrics",
		Description: "Get lyrics for the current song",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "artist",
				Description: "The name of the artist",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "title",
				Description: "The name of the song",
				Required:    false,
			},
		},
	},
}

func RegisterCommands(s *discordgo.Session) {
	// Replace GuildId with your actual test server ID if it's hardcoded
	cmds, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, GuildId, commands)
	if err != nil {
		log.Error("Failed to register commands: ", err)
	} else {
		log.Infof("Successfully registered %d commands for Guild: %s", len(cmds), GuildId)
	}
}
