package appInitialize

import "pinnacle-primary-be/internal/app/ping"

func init() {
	apps = append(apps, &ping.Ping{Name: "ping module"})
}
