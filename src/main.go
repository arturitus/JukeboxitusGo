package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"
	"gopkg.in/yaml.v3"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgolink/v3/disgolink"

	"jukeboxitus/src/bot"
	bot_config "jukeboxitus/src/bot/config"
	"jukeboxitus/src/build"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelInfo)
	log.Info("starting discordgo example...")
	log.Info("discordgo version: ", discordgo.VERSION)
	log.Info("disgolink version: ", disgolink.Version)

	config, err := loadConfig(build.ConfigFile)
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

	b := &bot.Bot{
		Queues: &bot.QueueManager{
			Queues: make(map[string]*bot.Queue),
		},
		SearchType: bot_config.ParseSearchType(searchTypeStr),
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	b.Session = session

	session.State.TrackVoice = true
	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentsGuildVoiceStates

	session.AddHandler(b.OnApplicationCommand)
	session.AddHandler(b.OnVoiceServerUpdate)
	session.AddHandler(b.OnVoiceStateUpdate)

	if err = session.Open(); err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	bot.RegisterCommands(session)

	b.Lavalink = disgolink.New(snowflake.MustParse(session.State.User.ID),
		disgolink.WithListenerFunc(b.OnPlayerPause),
		disgolink.WithListenerFunc(b.OnPlayerResume),
		disgolink.WithListenerFunc(b.OnTrackStart),
		disgolink.WithListenerFunc(b.OnTrackEnd),
		disgolink.WithListenerFunc(b.OnTrackException),
		disgolink.WithListenerFunc(b.OnTrackStuck),
		disgolink.WithListenerFunc(b.OnWebSocketClosed),
	)
	b.Handlers = map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error{
		"play":        b.Play,
		"pause":       b.Pause,
		"now-playing": b.NowPlaying,
		"stop":        b.Stop,
		"skip":        b.Skip,
		"queue":       b.Queue,
		"clear-queue": b.ClearQueue,
		"queue-type":  b.QueueType,
		"shuffle":     b.Shuffle,
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

func loadConfig(filePath string) (*bot_config.Config, error) {
	embeded, err := build.GetEmbeddedConfig()
	if err != nil {
		return nil, err
	}

	data, err := embeded.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config bot_config.Config
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
