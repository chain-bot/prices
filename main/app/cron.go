package app

import (
	"fmt"
	"github.com/robfig/cron"
)

func StartScrapperCron() {
	c := cron.New()
	_ = c.AddFunc("@every 1m", func() { fmt.Println("Every minute") })
	c.Start()
}
