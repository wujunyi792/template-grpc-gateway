package corn

import (
	"github.com/robfig/cron"
	"pinnacle-primary-be/pkg/logger"
)

func init() {
	c := cron.New()
	err := c.AddFunc("0 0/10 * * * *", func() {})
	if err != nil {
		logger.NameSpace("cron").Fatal(err)
	}
	c.Start()
	logger.NameSpace("cron").Info("corn routerInitialize SUCCESS ")
}
