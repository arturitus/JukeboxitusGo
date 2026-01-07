package bot

import (
	"context"
	"os"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"

	bot_config "jukeboxitus/src/bot/config"
)

type Bot struct {
	Session     *discordgo.Session
	Lavalink    disgolink.Client
	Handlers    map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error
	Queues      *QueueManager
	SearchType  bot_config.SearchType
	GeniusToken string
}

var (
	urlPattern    = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)
	GuildId       = os.Getenv("GUILD_ID")
)

func (b *Bot) OnApplicationCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {
	data := event.ApplicationCommandData()

	handler, ok := b.Handlers[data.Name]
	if !ok {
		log.Info("unknown command: ", data.Name)
		return
	}
	if err := handler(event, data); err != nil {
		log.Error("error handling command: ", err)
	}
}

func (b *Bot) OnVoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	if event.UserID != session.State.User.ID {
		return
	}

	var channelID *snowflake.ID
	if event.ChannelID != "" {
		id := snowflake.MustParse(event.ChannelID)
		channelID = &id
	}
	time.Sleep(500 * time.Millisecond)
	b.Lavalink.OnVoiceStateUpdate(context.Background(), snowflake.MustParse(event.GuildID), channelID, event.SessionID)
	if event.ChannelID == "" {
		b.Queues.Delete(event.GuildID)
	}
}

func (b *Bot) OnVoiceServerUpdate(session *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	time.Sleep(500 * time.Millisecond)
	b.Lavalink.OnVoiceServerUpdate(context.Background(), snowflake.MustParse(event.GuildID), event.Token, event.Endpoint)
}
