package app

import (
	"fmt"
	"github.com/robfig/cron"
)

func StartScrapperCron() (*cron.Cron, error) {
	c := cron.New()
	err := c.AddFunc("@every 1m", func() { fmt.Println("Every minute") })
	if err != nil {
		return nil, err
	}
	c.Start()
	return c, nil
}
