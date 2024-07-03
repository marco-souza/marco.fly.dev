// This file contains the scheduler package which is responsible for scheduling the tasks to be executed.
package scheduler

import (
	"log"

	"github.com/marco-souza/marco.fly.dev/internal/cron"
	"github.com/marco-souza/marco.fly.dev/internal/db"
	"github.com/marco-souza/marco.fly.dev/internal/lua"
)

// Setup scheduler by initializing cronjobs, registering lua scripts persisted
func Setup() error {
	log.Println("setting up scheduler")
	cron.CronService.Start()

	log.Println("loading persisted cron jobs")
	// load persisted cron jobs
	crons, err := db.Queries.ListCronJobs(db.Ctx)
	if err != nil {
		log.Println("error loading persisted cron jobs: ", err)
		return err
	}

	if len(crons) == 0 {
		log.Println("no job found")
		return nil
	}

	log.Println("setup persisted cron jobs: ", len(crons))
	for _, c := range crons {
		// TODO: running id
		err := cron.CronService.Add(c.Expression, func() {
			log.Println("executing cron job: ", c.Name)

			output, err := lua.Runtime.Run(c.Script)
			if err != nil {
				log.Println("error running lua script: ", err)
				return
			}

			log.Println("lua script output: ", output)
		})

		if err != nil {
			log.Println("error adding cron job: ", err)
			return err
		}
	}

	return nil
}

// Close scheduler by stopping the cron CronService
func Close() error {
	log.Println("closing scheduler")
	cron.CronService.Stop()
	return nil
}
