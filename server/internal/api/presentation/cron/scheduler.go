package cron

import (
	"context"

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

func (c *Scheduler) Schedule(ctx context.Context, task func() error) error {
	_, err := c.scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(3, 0, 0))),
		gocron.NewTask(task))
	if err != nil {
		return err
	}
	return nil
}

func (c *Scheduler) Start() {
	c.scheduler.Start()
}

func (c *Scheduler) Stop() error {
	return c.scheduler.Shutdown()
}
