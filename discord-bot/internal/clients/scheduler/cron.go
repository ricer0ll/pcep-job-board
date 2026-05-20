package scheduler

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/go-co-op/gocron/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/workday"
)

func InitCronJob(client *bot.Client) gocron.Scheduler {
	location, _ := time.LoadLocation("America/Los_Angeles")
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(location))
	if err != nil {
		panic("Failed to start cron scheduler!")
	}

	workday.InitJobsCache()

	j, err := scheduler.NewJob(
		gocron.CronJob(
			"0 6-12 * * 1-5",
			false,
		),
		gocron.NewTask(workday.GetNewJobPostings, client),
	)

	if err != nil {
		panic("Failed to create job for scheduler")
	}

	scheduler.Start()
	slog.Info(fmt.Sprintf("Started cron job with id %s", j.ID()))
	return scheduler
}
