package appInitialize

import (
	"pinnacle-primary-be/internal/app/healthz"
)

func init() {
	apps = append(apps, &healthz.Healthz{Name: "healthz module"})
}
