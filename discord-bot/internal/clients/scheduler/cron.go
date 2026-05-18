package scheduler

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/go-co-op/gocron/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/workday"
)

func InitCronJob(client *bot.Client) gocron.Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic("Failed to start cron scheduler!")
	}

	workday.InitJobsCache()

	j, err := scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(6, 0, 0),
				gocron.NewAtTime(8, 0, 0),
				gocron.NewAtTime(12, 0, 0),
			),
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
