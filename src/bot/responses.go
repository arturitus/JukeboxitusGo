package bot

import (
	"github.com/bwmarrin/discordgo"
)

// Constants for consistent styling
const (
	ColorDefault = 0x5865F2 // Blurple
	ColorSuccess = 0x57F287 // Green
	ColorWarning = 0xFEE75C // Yellow
	ColorError   = 0xED4245 // Red
)

const (
	IconPlay  = "▶️"
	IconPause = "⏸️"
	IconStop  = "⏹️"

	IconSkip    = "⏭️"
	IconShuffle = "🔀"
	IconRepeat  = "🔁"
	IconQueue   = "📜"
	IconSearch  = "🔍"
	IconSuccess = "✅"
	IconError   = "❌"
	IconEmpty   = "🏜️"

	IconVolume = "🔊"
	IconBass   = "🎚️"
	IconEightD = "🎧"
)

// SendResponse is your central "printing" function
func (b *Bot) SendResponse(i *discordgo.Interaction, title string, description string, color int) error {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Jukeboxitus Music",
		},
	}

	// Check if the interaction was already deferred/responded to
	// We do this by checking the Interaction's internal state if possible,
	// but a common pattern is to try an Edit if you know you deferred.
	// For a generic helper, we can try to respond, and if it fails because it's already acknowledged, we edit.

	err := b.Session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})

	if err != nil {
		// If we already deferred (like in Play), we must use Edit instead
		_, err = b.Session.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		})
	}

	return err
}

// SendComplexResponse handles cards with thumbnails and specific statuses
func (b *Bot) SendComplexResponse(i *discordgo.Interaction, title string, description string, thumbURL string, color int) error {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Jukeboxitus Music",
		},
	}

	// Add thumbnail only if a URL is provided
	if thumbURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: thumbURL,
		}
	}

	// Try to respond first (for commands like NowPlaying)
	err := b.Session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})

	// If already acknowledged (for deferred commands like Play), edit the response
	if err != nil {
		_, err = b.Session.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		})
	}

	return err
}
