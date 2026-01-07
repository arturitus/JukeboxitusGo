package bot

import (
	"context"
	"fmt"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (b *Bot) OnPlayerPause(player disgolink.Player, event lavalink.PlayerPauseEvent) {
	fmt.Printf("onPlayerPause: %v\n", event)
}

func (b *Bot) OnPlayerResume(player disgolink.Player, event lavalink.PlayerResumeEvent) {
	fmt.Printf("onPlayerResume: %v\n", event)
}

func (b *Bot) OnTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	fmt.Printf("onTrackStart: %v\n", event)
}

func (b *Bot) OnTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	fmt.Printf("onTrackEnd: %v\n", event)

	// MayStartNext is false if the track was stopped or replaced manually
	if !event.Reason.MayStartNext() {
		return
	}

	queue := b.Queues.Get(event.GuildID().String())
	var (
		nextTrack lavalink.Track
		ok        bool
	)

	switch queue.Type {
	case QueueTypeNormal:
		nextTrack, ok = queue.Next()

	case QueueTypeRepeatTrack:
		nextTrack = event.Track
		ok = true // Fix: Ensure we proceed to play

	case QueueTypeRepeatQueue:
		queue.Add(event.Track)
		nextTrack, ok = queue.Next()
		// ok will be true because we just added a track
	}

	// If no tracks are left in the queue and we aren't repeating
	if !ok {
		return
	}

	// Play the determined next track
	if err := player.Update(context.Background(), lavalink.WithTrack(nextTrack)); err != nil {
		log.Error("Failed to play next track: ", err)
	}
}

func (b *Bot) OnTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	fmt.Printf("onTrackException: %v\n", event)
}

func (b *Bot) OnTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	fmt.Printf("onTrackStuck: %v\n", event)
}

func (b *Bot) OnWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
	fmt.Printf("onWebSocketClosed: %v\n", event)
}
