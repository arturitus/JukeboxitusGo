package bot

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (b *Bot) Shuffle(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := b.Queues.Get(event.GuildID)

	// Error: No tracks to shuffle
	if queue == nil || len(queue.Tracks) == 0 {
		return b.SendResponse(
			event.Interaction,
			"Queue Error",
			fmt.Sprintf("%s There is no active queue to shuffle right now.", IconEmpty),
			ColorError,
		)
	}

	queue.Shuffle()

	// Success
	return b.SendResponse(
		event.Interaction,
		"Queue Shuffled",
		fmt.Sprintf("%s Successfully shuffled **%d** tracks!", IconShuffle, len(queue.Tracks)),
		ColorSuccess,
	)
}

func (b *Bot) Skip(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	// 1. Get the player
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.SendResponse(event.Interaction, "Playback Error",
			fmt.Sprintf("%s No active player found.", IconError), ColorError)
	}

	// 2. Get the queue
	queue := b.Queues.Get(event.GuildID)
	if queue == nil {
		return b.SendResponse(event.Interaction, "Queue Error",
			fmt.Sprintf("%s No queue found for this server.", IconEmpty), ColorError)
	}

	// 3. Try to get the next track
	nextTrack, ok := queue.Next()
	if !ok {
		return b.SendResponse(event.Interaction, "End of Queue",
			fmt.Sprintf("%s No more tracks to skip to.", IconEmpty), ColorWarning)
	}

	// 4. Update the player with the new track
	err := player.Update(context.Background(), lavalink.WithTrack(nextTrack))
	if err != nil {
		return b.SendResponse(event.Interaction, "Playback Error",
			fmt.Sprintf("%s Error while playing the next track: `%s`", IconError, err), ColorError)
	}

	// 5. Success Card
	return b.SendResponse(
		event.Interaction,
		"Track Skipped",
		fmt.Sprintf("%s Skipped to: **[`%s`](<%s>)**", IconSkip, nextTrack.Info.Title, *nextTrack.Info.URI),
		ColorSuccess,
	)
}

func (b *Bot) QueueType(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := b.Queues.Get(event.GuildID)

	if queue == nil {
		return b.SendResponse(event.Interaction, "Configuration Error",
			fmt.Sprintf("%s No active player found.", IconError), ColorError)
	}

	// Capture the value from Discord
	val := data.Options[0].StringValue()
	newType := QueueType(val)

	// Check if the type is actually valid by checking our String() output
	if newType.String() == "unknown" {
		return b.SendResponse(event.Interaction, "Configuration Error",
			fmt.Sprintf("%s Invalid queue mode selected.", IconError), ColorError)
	}

	queue.Type = newType

	return b.SendResponse(
		event.Interaction,
		"Queue Mode Updated",
		fmt.Sprintf("%s Queue mode has been set to: **%s**", IconRepeat, queue.Type.String()),
		ColorSuccess,
	)
}

func (b *Bot) ClearQueue(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := b.Queues.Get(event.GuildID)

	// 1. Error: No queue exists for this guild
	if queue == nil {
		return b.SendResponse(
			event.Interaction,
			"Queue Error",
			fmt.Sprintf("%s No active player or queue found.", IconError),
			ColorError,
		)
	}

	// 2. Logic: Clear the tracks
	// We can check the count before clearing to give a more detailed message
	count := len(queue.Tracks)
	queue.Clear()

	// 3. Success Card
	return b.SendResponse(
		event.Interaction,
		"Queue Cleared",
		fmt.Sprintf("%s Successfully removed **%d** tracks from the queue.", IconSuccess, count),
		ColorSuccess,
	)
}

func (b *Bot) Queue(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := b.Queues.Get(event.GuildID)

	// 1. Error: No queue found
	if queue == nil {
		return b.SendResponse(event.Interaction, "Queue Status",
			fmt.Sprintf("%s No active player found.", IconError), ColorError)
	}

	// 2. Case: Empty Queue
	if len(queue.Tracks) == 0 {
		return b.SendResponse(event.Interaction, "Queue Status",
			fmt.Sprintf("%s The queue is currently empty.", IconEmpty), ColorDefault)
	}

	// 3. Logic: Build the track list string
	var tracks string
	for i, track := range queue.Tracks {
		// Stop adding if we approach the embed description limit (4096)
		line := fmt.Sprintf("**%d.** [`%s`](<%s>)\n", i+1, track.Info.Title, *track.Info.URI)
		if len(tracks)+len(line) > 4000 {
			tracks += "...and more"
			break
		}
		tracks += line
	}

	// 4. Success Card
	return b.SendResponse(
		event.Interaction,
		fmt.Sprintf("%s Current Queue (%s)", IconQueue, queue.Type.String()),
		tracks,
		ColorDefault,
	)
}
