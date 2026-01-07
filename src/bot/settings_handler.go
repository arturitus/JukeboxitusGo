package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
	// lyrics "github.com/rhnvrm/lyric-api-go"
)

// Volume handles setting the player volume
func (b *Bot) Volume(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.SendResponse(event.Interaction, "Setting Error",
			fmt.Sprintf("%s No active player found.", IconError), ColorError)
	}

	volume := int(data.Options[0].IntValue())
	if err := player.Update(context.Background(), lavalink.WithVolume(volume)); err != nil {
		return b.SendResponse(event.Interaction, "Setting Error",
			fmt.Sprintf("%s Could not set volume: `%s`", IconError, err), ColorError)
	}

	return b.SendResponse(event.Interaction, "Volume Updated",
		fmt.Sprintf("%s Volume set to **%d%%**", IconVolume, volume), ColorSuccess)
}

// BassBoost toggles a heavy bass equalizer
func (b *Bot) BassBoost(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.SendResponse(event.Interaction, "Setting Error",
			fmt.Sprintf("%s No active player found.", IconError), ColorError)
	}

	var filters lavalink.Filters
	var statusMsg string
	var color int

	// 1. Better Toggle Logic: Check if the first band has boost
	isBoosted := false
	currentEq := player.Filters().Equalizer
	if currentEq != nil && len(currentEq) > 0 && currentEq[0] > 0 {
		isBoosted = true
	}

	if isBoosted {
		// Resetting: Create an EQ with all 0.0 values
		// filters.Equalizer = &lavalink.Equalizer{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		filters.Equalizer = nil
		statusMsg = fmt.Sprintf("%s Bass Boost: **OFF**", IconBass)
		color = ColorDefault
	} else {
		// Applying: Boost the first 3 bands
		// Applying: Boost the first 3 bands significantly
		// Band 0: 25Hz, Band 1: 40Hz, Band 2: 63Hz
		// filters.Equalizer = &lavalink.Equalizer{1.0, 0.8, 0.6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		statusMsg = fmt.Sprintf("%s Bass Boost: **ON**", IconBass)
		color = ColorSuccess
	}

	if err := player.Update(context.Background(), lavalink.WithFilters(filters)); err != nil {
		return b.SendResponse(event.Interaction, "Filter Error",
			fmt.Sprintf("%s Failed to apply filters: `%s`", IconError, err), ColorError)
	}

	return b.SendResponse(event.Interaction, "Filter Updated", statusMsg, color)
}

// EightD toggles the Rotation filter for an 8D effect
func (b *Bot) EightD(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.SendResponse(event.Interaction, "Setting Error",
			fmt.Sprintf("%s No active player found.", IconError), ColorError)
	}

	var filters lavalink.Filters

	// Check current state safely.
	// If the 3rd party lib is failing to unmarshal the response,
	// we might need to track the toggle state ourselves in a map.
	isCurrentlyOn := player.Filters().Rotation != nil

	if isCurrentlyOn {
		filters.Rotation = nil
		b.SendResponse(event.Interaction, "Filter Updated", fmt.Sprintf("%s 8-D Audio: **OFF**", IconEightD), ColorDefault)
	} else {
		// We MUST use an int here because the library says so.
		// We'll use 1 (1Hz).
		filters.Rotation = &lavalink.Rotation{
			RotationHz: 1,
		}
		b.SendResponse(event.Interaction, "Filter Updated", fmt.Sprintf("%s 8-D Audio: **ON**", IconEightD), ColorSuccess)
	}

	// Send the update to Lavalink
	return player.Update(context.Background(), lavalink.WithFilters(filters))
}

func (b *Bot) Lyrics(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	var artist, title string

	// 1. Extract optional arguments from the interaction
	argMap := make(map[string]string)
	for _, opt := range data.Options {
		argMap[opt.Name] = opt.StringValue()
	}

	artist = cleanLyricQuery(argMap["artist"])
	title = cleanLyricQuery(argMap["title"])

	// 2. Fallback: If both are empty, use the current player track
	if artist == "" && title == "" {
		player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
		if player == nil || player.Track() == nil {
			return b.SendResponse(event.Interaction, "Lyrics Error",
				fmt.Sprintf("%s No song is currently playing and no search terms provided.", IconError), ColorError)
		}

		fullTitle := player.Track().Info.Title
		parts := strings.Split(fullTitle, "-")
		if len(parts) > 1 {
			artist = cleanLyricQuery(parts[0])
			title = cleanLyricQuery(parts[1])
		} else {
			title = cleanLyricQuery(fullTitle)
		}
	}

	// 3. Defer response (Genius search + Scraping takes time)
	b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// 4. Fetch Lyrics using your Genius function
	lyricText, err := b.getGeniusLyrics(artist, title)

	// Fallback: If search failed with artist+title, try searching with just title
	if (err != nil || lyricText == "") && artist != "" {
		lyricText, err = b.getGeniusLyrics("", title)
	}

	if err != nil || lyricText == "" {
		return b.SendResponse(event.Interaction, "Lyrics Not Found",
			fmt.Sprintf("%s Could not find lyrics for **%s %s**", IconSearch, artist, title), ColorWarning)
	}

	// 5. Truncate for Discord limits (4096 characters)
	if len(lyricText) > 4000 {
		lyricText = lyricText[:3997] + "..."
	}

	// 6. Final response
	displayHeader := fmt.Sprintf("Lyrics: %s", title)
	if artist != "" {
		displayHeader = fmt.Sprintf("Lyrics: %s - %s", artist, title)
	}

	return b.SendResponse(event.Interaction, displayHeader, lyricText, ColorDefault)
}

type GeniusResponse struct {
	Response struct {
		Hits []struct {
			Result struct {
				URL           string `json:"url"`
				Title         string `json:"title"`
				PrimaryArtist struct {
					Name string `json:"name"`
				} `json:"primary_artist"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

func (b *Bot) getGeniusLyrics(artist, title string) (string, error) {

	// Merge artist and title for a clean search query
	query := fmt.Sprintf("%s %s", artist, title)
	searchURL := fmt.Sprintf("https://api.genius.com/search?q=%s", url.QueryEscape(query))

	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("Authorization", "Bearer "+b.GeniusToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var gResp GeniusResponse
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		return "", err
	}

	if len(gResp.Response.Hits) == 0 {
		return "", fmt.Errorf("no results found for %s", query)
	}

	// Genius returns "hits" ordered by relevance. The first one is almost always correct.
	songURL := gResp.Response.Hits[0].Result.URL

	return scrapeLyrics(songURL)
}

func scrapeLyrics(geniusURL string) (string, error) {
	// We need a User-Agent so Genius doesn't block the request
	req, _ := http.NewRequest("GET", geniusURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("genius returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var lyricsBuilder strings.Builder

	// Genius uses the 'data-lyrics-container' attribute for the actual lyric text
	doc.Find("div[data-lyrics-container='true']").Each(func(i int, s *goquery.Selection) {
		// Replace <br> tags with actual newlines
		s.Find("br").ReplaceWithHtml("\n")

		// Get the text and append it
		lyricsBuilder.WriteString(s.Text() + "\n")
	})

	finalLyrics := strings.TrimSpace(lyricsBuilder.String())

	if finalLyrics == "" {
		return "", fmt.Errorf("could not parse lyrics from page")
	}

	return finalLyrics, nil
}

func cleanLyricQuery(input string) string {
	// Remove common YouTube suffixes
	re := regexp.MustCompile(`(?i)\(?(official|music|video|lyric|hd|audio|visualizer)\)?`)
	input = re.ReplaceAllString(input, "")

	// Remove brackets/parentheses and their contents if they are left over
	re = regexp.MustCompile(`\[.*?\]|\(.*?\)|feat\..*`)
	input = re.ReplaceAllString(input, "")

	return strings.TrimSpace(input)
}
