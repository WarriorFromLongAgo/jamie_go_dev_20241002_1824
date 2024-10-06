package anvil

import (
	"fmt"
	"your_project/config" // Adjust the import path as necessary
)

// GetAnvilURL constructs and returns the URL for the Anvil service
func GetAnvilURL(cfg *config.Configuration) string {
	return fmt.Sprintf("http://%s:%d", cfg.Anvil.Host, cfg.Anvil.Port)
}