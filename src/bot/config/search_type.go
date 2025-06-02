package bot_config

import "strings"

type SearchType byte

const (
	YouTube SearchType = iota
	YouTubeMusic
	SoundCloud
)

const defaultSearchType = YouTube

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

func ParseSearchType(str string) SearchType {
	str = strings.TrimSpace(str)

	if val, ok := stringToSearchType[str]; ok {
		return val
	}
	return defaultSearchType
}
