//go:build windows

package build

import (
	"embed"
	"fmt"
)

//go:embed config.yaml
var embeddedConfigFS embed.FS

const ConfigFile = "config.yaml"

// GetEmbeddedConfig returns the embedded config file system for Windows.
// This function will only be available when building for Windows.
func GetEmbeddedConfig() (embed.FS, error) {
	fmt.Println("INFO: Building for Windows. 'config.yaml' is embedded.")
	return embeddedConfigFS, nil
}
