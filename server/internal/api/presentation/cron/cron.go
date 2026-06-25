package cron

import (
	"context"

	"go.uber.org/fx"
	"starliner.app/internal/api/application"
)

type Cron struct {
	scheduler               *Scheduler
	organizationApplication *application.OrganizationApplication
}

func NewCron(scheduler *Scheduler, organizationApplication *application.OrganizationApplication) *Cron {
	return &Cron{scheduler: scheduler, organizationApplication: organizationApplication}
}

func RegisterCron(lc fx.Lifecycle, cr *Cron) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { return cr.StartCron() },
		OnStop:  func(ctx context.Context) error { return cr.StopCron() },
	})
}

func (c *Cron) StartCron() error {
	err := c.scheduler.Schedule("0 3 * * *", func() error {
		return c.organizationApplication.CleanUpOldInvites(context.Background())
	})
	if err != nil {
		return err
	}
	c.scheduler.Start()
	return nil
}

func (c *Cron) StopCron() error {
	return c.scheduler.Stop()
}
