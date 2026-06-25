package cron

import (
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	scheduler gocron.Scheduler
}

func NewScheduler() (*Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &Scheduler{scheduler: scheduler}, nil
}

func (c *Scheduler) Schedule(crontab string, task func() error) error {
	_, err := c.scheduler.NewJob(gocron.CronJob(crontab, false), gocron.NewTask(task))
	return err
}

func (c *Scheduler) Start() {
	c.scheduler.Start()
}

func (c *Scheduler) Stop() error {
	return c.scheduler.Shutdown()
}
