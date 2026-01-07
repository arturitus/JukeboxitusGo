package bot

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/bwmarrin/discordgo"
// 	"github.com/disgoorg/snowflake/v2"

// 	"github.com/disgoorg/disgolink/v3/disgolink"
// 	"github.com/disgoorg/disgolink/v3/lavalink"

// 	bot_config "jukeboxitus/src/bot/config"
// )

// func (b *Bot) Shuffle(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
// 	queue := b.Queues.Get(event.GuildID)

// 	// Error: No tracks to shuffle
// 	if queue == nil || len(queue.Tracks) == 0 {
// 		return b.SendResponse(
// 			event.Interaction,
// 			"Queue Error",
// 			fmt.Sprintf("%s There is no active queue to shuffle right now.", IconEmpty),
// 			ColorError,
// 		)
// 	}

// 	queue.Shuffle()

// 	// Success
// 	return b.SendResponse(
// 		event.Interaction,
// 		"Queue Shuffled",
// 		fmt.Sprintf("%s Successfully shuffled **%d** tracks!", IconShuffle, len(queue.Tracks)),
// 		ColorSuccess,
// 	)
// }

// func (b *Bot) Skip(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
// 	// 1. Get the player
// 	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
// 	if player == nil {
// 		return b.SendResponse(event.Interaction, "Playback Error",
// 			fmt.Sprintf("%s No active player found.", IconError), ColorError)
// 	}

// 	// 2. Get the queue
// 	queue := b.Queues.Get(event.GuildID)
// 	if queue == nil {
// 		return b.SendResponse(event.Interaction, "Queue Error",
// 			fmt.Sprintf("%s No queue found for this server.", IconEmpty), ColorError)
// 	}

// 	// 3. Try to get the next track
// 	nextTrack, ok := queue.Next()
// 	if !ok {
// 		return b.SendResponse(event.Interaction, "End of Queue",
// 			fmt.Sprintf("%s No more tracks to skip to.", IconEmpty), ColorWarning)
// 	}

// 	// 4. Update the player with the new track
// 	err := player.Update(context.Background(), lavalink.WithTrack(nextTrack))
// 	if err != nil {
// 		return b.SendResponse(event.Interaction, "Playback Error",
// 			fmt.Sprintf("%s Error while playing the next track: `%s`", IconError, err), ColorError)
// 	}

// 	// 5. Success Card
// 	return b.SendResponse(
// 		event.Interaction,
// 		"Track Skipped",
// 		fmt.Sprintf("%s Skipped to: **[`%s`](<%s>)**", IconSkip, nextTrack.Info.Title, *nextTrack.Info.URI),
// 		ColorSuccess,
// 	)
// }

// func (b *Bot) QueueType(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
// 	queue := b.Queues.Get(event.GuildID)

// 	if queue == nil {
// 		return b.SendResponse(event.Interaction, "Configuration Error",
// 			fmt.Sprintf("%s No active player found.", IconError), ColorError)
// 	}

// 	// Capture the value from Discord
// 	val := data.Options[0].StringValue()
// 	newType := QueueType(val)

// 	// Check if the type is actually valid by checking our String() output
// 	if newType.String() == "unknown" {
// 		return b.SendResponse(event.Interaction, "Configuration Error",
// 			fmt.Sprintf("%s Invalid queue mode selected.", IconError), ColorError)
// 	}

// 	queue.Type = newType

// 	return b.SendResponse(
// 		event.Interaction,
// 		"Queue Mode Updated",
// 		fmt.Sprintf("%s Queue mode has been set to: **%s**", IconRepeat, queue.Type.String()),
// 		ColorSuccess,
// 	)
// }

// func (b *Bot) ClearQueue(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
// 	queue := b.Queues.Get(event.GuildID)

// 	// 1. Error: No queue exists for this guild
// 	if queue == nil {
// 		return b.SendResponse(
// 			event.Interaction,
// 			"Queue Error",
// 			fmt.Sprintf("%s No active player or queue found.", IconError),
// 			ColorError,
// 		)
// 	}

// 	// 2. Logic: Clear the tracks
// 	// We can check the count before clearing to give a more detailed message
// 	count := len(queue.Tracks)
// 	queue.Clear()

// 	// 3. Success Card
// 	return b.SendResponse(
// 		event.Interaction,
// 		"Queue Cleared",
// 		fmt.Sprintf("%s Successfully removed **%d** tracks from the queue.", IconSuccess, count),
// 		ColorSuccess,
// 	)
// }

// func (b *Bot) Queue(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
// 	queue := b.Queues.Get(event.GuildID)

// 	// 1. Error: No queue found
// 	if queue == nil {
// 		return b.SendResponse(event.Interaction, "Queue Status",
// 			fmt.Sprintf("%s No active player found.", IconError), ColorError)
// 	}

// 	// 2. Case: Empty Queue
// 	if len(queue.Tracks) == 0 {
// 		return b.SendResponse(event.Interaction, "Queue Status",
// 			fmt.Sprintf("%s The queue is currently empty.", IconEmpty), ColorDefault)
// 	}

// 	// 3. Logic: Build the track list string
// 	var tracks string
// 	for i, track := range queue.Tracks {
// 		// Stop adding if we approach the embed description limit (4096)
// 		line := fmt.Sprintf("**%d.** [`%s`](<%s>)\n", i+1, track.Info.Title, *track.Info.URI)
// 		if len(tracks)+len(line) > 4000 {
// 			tracks += "...and more"
// 			break
// 		}
// 		tracks += line
// 	}

// 	// 4. Success Card
// 	return b.SendResponse(
// 		event.Interaction,
// 		fmt.Sprintf("%s Current Queue (%s)", IconQueue, queue.Type.String()),
// 		tracks,
// 		ColorDefault,
// 	)
// }

// func (b *Bot) Pause(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
// 	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
// 	if player == nil {
// 		return b.SendResponse(event.Interaction, "Playback Error",
// 			fmt.Sprintf("%s No active player found.", IconError), ColorError)
// 	}

// 	// Toggle the current paused state
// 	willPause := !player.Paused()
// 	err := player.Update(context.Background(), lavalink.WithPaused(willPause))

// 	if err != nil {
// 		return b.SendResponse(event.Interaction, "Playback Error",
// 			fmt.Sprintf("%s Error while updating player: `%s`", IconError, err), ColorError)
// 	}

// 	// Prepare UI based on the new state
// 	status := "Resumed"
// 	icon := IconPlay
// 	color := ColorSuccess

// 	if willPause {
// 		status = "Paused"
// 		icon = IconPause
// 		color = ColorWarning // Yellow is great for "Paused/Wait"
// 	}

// 	return b.SendResponse(
// 		event.Interaction,
// 		"Player Status",
// 		fmt.Sprintf("%s Player is now **%s**", icon, status),
// 		color,
// 	)
// }

// func (b *Bot) Stop(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
// 	// 1. Check if player exists
// 	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
// 	if player == nil {
// 		return b.SendResponse(event.Interaction, "Playback Error",
// 			fmt.Sprintf("%s No active player found.", IconError), ColorError)
// 	}

// 	// 2. Disconnect from voice channel
// 	// Passing an empty string for ChannelID tells Discord to leave.
// 	if err := b.Session.ChannelVoiceJoinManual(event.GuildID, "", false, false); err != nil {
// 		return b.SendResponse(event.Interaction, "Connection Error",
// 			fmt.Sprintf("%s Error while disconnecting: `%s`", IconError, err), ColorError)
// 	}

// 	// 3. Success Card
// 	return b.SendResponse(
// 		event.Interaction,
// 		"Disconnected",
// 		fmt.Sprintf("%s The player has been stopped and I have left the voice channel.", IconStop),
// 		ColorError, // Red is standard for "Stop/Disconnect"
// 	)
// }

// func (b *Bot) NowPlaying(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
// 	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
// 	if player == nil {
// 		return b.SendResponse(event.Interaction, "Player Status", fmt.Sprintf("%s No player found.", IconError), ColorError)
// 	}

// 	track := player.Track()
// 	if track == nil {
// 		return b.SendResponse(event.Interaction, "Player Status", fmt.Sprintf("%s Nothing playing.", IconEmpty), ColorDefault)
// 	}

// 	description := fmt.Sprintf("%s **Currently Playing**\n[`%s`](<%s>)\n\n`%s / %s`",
// 		IconPlay, track.Info.Title, *track.Info.URI, formatPosition(player.Position()), formatPosition(track.Info.Length))

// 	// Using the new complex function to show the artwork!
// 	return b.SendComplexResponse(
// 		event.Interaction,
// 		"Now Playing",
// 		description,
// 		*track.Info.ArtworkURL,
// 		ColorDefault,
// 	)
// }

// func (b *Bot) Play(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
// 	identifier := data.Options[0].StringValue()

// 	// 1. Handle Search Types
// 	if !urlPattern.MatchString(identifier) && !searchPattern.MatchString(identifier) {
// 		switch b.SearchType {
// 		case bot_config.YouTube:
// 			identifier = lavalink.SearchTypeYouTube.Apply(identifier)
// 		case bot_config.YouTubeMusic:
// 			identifier = lavalink.SearchTypeYouTubeMusic.Apply(identifier)
// 		case bot_config.SoundCloud:
// 			identifier = lavalink.SearchTypeSoundCloud.Apply(identifier)
// 		}
// 	}

// 	// 2. Voice State Check
// 	voiceState, err := b.Session.State.VoiceState(event.GuildID, event.Member.User.ID)
// 	if err != nil {
// 		return b.SendResponse(event.Interaction, "Connection Error",
// 			fmt.Sprintf("%s You must be in a voice channel to play music!", IconError), ColorError)
// 	}

// 	// 3. Defer Response (Giving Lavalink time to work)
// 	if err := b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
// 	}); err != nil {
// 		return err
// 	}

// 	player := b.Lavalink.Player(snowflake.MustParse(event.GuildID))
// 	queue := b.Queues.Get(event.GuildID)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	var toPlay *lavalink.Track
// 	b.Lavalink.BestNode().LoadTracksHandler(ctx, identifier, disgolink.NewResultHandler(
// 		// --- SINGLE TRACK LOADED ---
// 		func(track lavalink.Track) {
// 			b.SendComplexResponse(event.Interaction, "Track Added",
// 				fmt.Sprintf("%s Added [`%s`](<%s>) to queue.", IconPlay, track.Info.Title, *track.Info.URI),
// 				*track.Info.ArtworkURL, ColorSuccess)

// 			if player.Track() == nil {
// 				toPlay = &track
// 			} else {
// 				queue.Add(track)
// 			}
// 		},
// 		// --- PLAYLIST LOADED ---
// 		func(playlist lavalink.Playlist) {
// 			b.SendComplexResponse(event.Interaction, "Playlist Added",
// 				fmt.Sprintf("%s Loaded **%d** tracks from playlist: `%s`", IconQueue, len(playlist.Tracks), playlist.Info.Name),
// 				*playlist.Tracks[0].Info.ArtworkURL, ColorSuccess)

// 			if player.Track() == nil {
// 				toPlay = &playlist.Tracks[0]
// 				queue.Add(playlist.Tracks[1:]...)
// 			} else {
// 				queue.Add(playlist.Tracks...)
// 			}
// 		},
// 		// --- SEARCH RESULT LOADED ---
// 		func(tracks []lavalink.Track) {
// 			b.SendComplexResponse(event.Interaction, "Search Result",
// 				fmt.Sprintf("%s Playing search result: [`%s`](<%s>)", IconSearch, tracks[0].Info.Title, *tracks[0].Info.URI),
// 				*tracks[0].Info.ArtworkURL, ColorSuccess)

// 			if player.Track() == nil {
// 				toPlay = &tracks[0]
// 			} else {
// 				queue.Add(tracks[0])
// 			}
// 		},
// 		// --- NOTHING FOUND ---
// 		func() {
// 			b.SendResponse(event.Interaction, "No Results",
// 				fmt.Sprintf("%s Nothing found for: `%s`", IconEmpty, identifier), ColorDefault)
// 		},
// 		// --- ERROR ---
// 		func(err error) {
// 			b.SendResponse(event.Interaction, "Search Error",
// 				fmt.Sprintf("%s Error: `%s`", IconError, err), ColorError)
// 		},
// 	))

// 	if toPlay == nil {
// 		return nil
// 	}

// 	// Join and Play
// 	if err := b.Session.ChannelVoiceJoinManual(event.GuildID, voiceState.ChannelID, false, false); err != nil {
// 		return err
// 	}

// 	return player.Update(context.Background(), lavalink.WithTrack(*toPlay))
// }

// func formatPosition(position lavalink.Duration) string {
// 	if position == 0 {
// 		return "0:00"
// 	}
// 	return fmt.Sprintf("%d:%02d", position.Minutes(), position.SecondsPart())
// }
