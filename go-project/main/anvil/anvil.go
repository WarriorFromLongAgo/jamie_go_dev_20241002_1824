package anvil

import (
	"fmt"

	"go-project/main/config"
)

// GetAnvilURL constructs and returns the URL for the Anvil service
func GetAnvilURL(cfg *config.Configuration) string {
	return fmt.Sprintf("http://%s:%d", cfg.Anvil.Host, cfg.Anvil.Port)
}
