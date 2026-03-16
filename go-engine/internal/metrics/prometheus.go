package metrics

import (
	"fmt"
	"time"
)

func Render(startedAt time.Time) string {
	uptimeSeconds := time.Since(startedAt).Seconds()

	return fmt.Sprintf(
		"# HELP go_engine_up Indicates whether go-engine is running.\n# TYPE go_engine_up gauge\ngo_engine_up 1\n# HELP go_engine_uptime_seconds Uptime in seconds.\n# TYPE go_engine_uptime_seconds counter\ngo_engine_uptime_seconds %.0f\n",
		uptimeSeconds,
	)
}
