package bot_config

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
