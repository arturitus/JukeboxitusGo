//go:build !windows
// +build !windows

package build

import (
	"embed"
	"errors"
	"fmt"
)

// embeddedConfigFS is an empty embed.FS for non-Windows builds,
// as config.yaml is not intended to be embedded on these platforms.
var embeddedConfigFS embed.FS

const ConfigFile = "config.yaml"

// GetEmbeddedConfig returns an error for non-Windows builds,
// indicating that config.yaml was not embedded.
func GetEmbeddedConfig() (embed.FS, error) {
	fmt.Println("INFO: Building for Non-Windows. 'config.yaml' is NOT embedded.")
	return embeddedConfigFS, errors.New("config.yaml is not embedded on non-Windows platforms")
}

func BuildType() string {
	return "not windows"
}
