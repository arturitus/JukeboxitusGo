package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"
	"gopkg.in/yaml.v3"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgolink/v3/disgolink"
)

//go:embed config.yaml
var configFile embed.FS

var (
	urlPattern    = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)
	GuildId       = os.Getenv("GUILD_ID")
)

type SearchType byte

const (
	YouTube SearchType = iota
	YouTubeMusic
	SoundCloud
)

var stringToSearchType = map[string]SearchType{
	"youTube":      YouTube,
	"youTubeMusic": YouTubeMusic,
	"soundCloud":   SoundCloud,
}

var searchTypeToString = map[SearchType]string{
	YouTube:      "youTube",
	YouTubeMusic: "youTubeMusic",
	SoundCloud:   "soundCloud",
}

func (s SearchType) String() string {
	if str, ok := searchTypeToString[s]; ok {
		return str
	}
	return searchTypeToString[YouTube]
}

const defaultSearchType = YouTube

func ParseSearchType(str string) SearchType {
	str = strings.TrimSpace(str)

	if val, ok := stringToSearchType[str]; ok {
		return val
	}
	return defaultSearchType
}

type LavalinkConfig struct {
	Name       string `yaml:"Name"`
	Hostname   string `yaml:"Hostname"`
	Port       int    `yaml:"Port"`
	Password   string `yaml:"Password"`
	Secured    bool   `yaml:"Secured"`
	SearchType string `yaml:"SearchType"`
}

type Config struct {
	Token    string         `yaml:"Token"`
	Lavalink LavalinkConfig `yaml:"Lavalink"`
}

type Bot struct {
	Session    *discordgo.Session
	Lavalink   disgolink.Client
	Handlers   map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error
	Queues     *QueueManager
	SearchType SearchType
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelInfo)
	log.Info("starting discordgo example...")
	log.Info("discordgo version: ", discordgo.VERSION)
	log.Info("disgolink version: ", disgolink.Version)

	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	token, tokenFromEnv := getEnv("TOKEN", config.Token)
	name, nameFromEnv := getEnv("NAME", config.Lavalink.Name)
	hostName, hostNameFromEnv := getEnv("HOSTNAME", config.Lavalink.Hostname)
	portStr, portFromEnv := getEnv("PORT", strconv.Itoa(config.Lavalink.Port))
	searchTypeStr, searchTypeFromEnv := getEnv("SEARCH_TYPE", config.Lavalink.SearchType)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("invalid PORT value: %v", err)
	}
	password, passwordFromEnv := getEnv("PASSWORD", config.Lavalink.Password)

	securedStr, securedFromEnv := getEnv("SECURED", strconv.FormatBool(config.Lavalink.Secured))
	secured, _ := strconv.ParseBool(securedStr)

	fmt.Printf("Token (%s): %q\n", checkSource(tokenFromEnv), token)
	fmt.Printf("Lavalink:\n")
	fmt.Printf("	Name (%s): %q\n", checkSource(nameFromEnv), name)
	fmt.Printf("	Hostname (%s): %q\n", checkSource(hostNameFromEnv), hostName)
	fmt.Printf("	Port (%s): %d\n", checkSource(portFromEnv), port)
	fmt.Printf("	Password (%s): %q\n", checkSource(passwordFromEnv), password)
	fmt.Printf("	Secured (%s): %v\n", checkSource(securedFromEnv), secured)
	fmt.Printf("	SerachType (%s): %v\n", checkSource(searchTypeFromEnv), searchTypeStr)

	b := &Bot{
		Queues: &QueueManager{
			queues: make(map[string]*Queue),
		},
		SearchType: ParseSearchType(searchTypeStr),
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	b.Session = session

	session.State.TrackVoice = true
	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentsGuildVoiceStates

	session.AddHandler(b.onApplicationCommand)
	session.AddHandler(b.onVoiceServerUpdate)
	session.AddHandler(b.onVoiceStateUpdate)

	if err = session.Open(); err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	registerCommands(session)

	b.Lavalink = disgolink.New(snowflake.MustParse(session.State.User.ID),
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
	)
	b.Handlers = map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error{
		"play":        b.play,
		"pause":       b.pause,
		"now-playing": b.nowPlaying,
		"stop":        b.stop,
		"skip":        b.skip,
		"queue":       b.queue,
		"clear-queue": b.clearQueue,
		"queue-type":  b.queueType,
		"shuffle":     b.shuffle,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	address := fmt.Sprintf("%s:%d", hostName, port)
	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     name,
		Address:  address,
		Password: password,
		Secure:   secured,
	})
	if err != nil {
		log.Fatal(err)
	}
	version, err := node.Version(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("node version: %s", version)

	log.Info("DiscordGo example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func (b *Bot) onApplicationCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {
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

func (b *Bot) onVoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
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

func (b *Bot) onVoiceServerUpdate(session *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	time.Sleep(500 * time.Millisecond)
	b.Lavalink.OnVoiceServerUpdate(context.Background(), snowflake.MustParse(event.GuildID), event.Token, event.Endpoint)
}

func loadConfig(filePath string) (*Config, error) {
	data, err := configFile.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getEnv(key, fallback string) (string, bool) {
	if value, exists := os.LookupEnv(key); exists {
		return value, true
	}
	return fallback, false
}

func checkSource(fromEnv bool) string {
	if fromEnv {
		return "env"
	}
	return "config"
}
