// service to manage cron jobs
package cron

import (
	"log"

	"github.com/marco-souza/marco.fly.dev/internal/db"
	"github.com/marco-souza/marco.fly.dev/internal/db/sqlc"
	"github.com/marco-souza/marco.fly.dev/internal/lua"
	"github.com/robfig/cron/v3"
)

type runningCronJobs map[int]cron.EntryID

type Cron struct {
	sqlc.Cron
	status string
}

var (
	scheduler   = cron.New()
	runningJobs = runningCronJobs{}
)

func Start() error {
	log.Println("starting scheduler")
	scheduler.Start()
	return registerPersistedJobs()
}

func Stop() {
	log.Println("stopping scheduler")
	scheduler.Stop()
}

func AddScript(name, cronExpr, script string) error {
	job, err := db.Queries.CreateCronJob(db.Ctx, sqlc.CreateCronJobParams{
		Name:       name,
		Expression: cronExpr,
		Script:     script,
	})
	if err != nil {
		return err
	}

	// register cron job
	scriptHandler := func() {
		log.Printf("executing cron job: %s (%e)\n", job.Name, err)
		if _, err := lua.Run(job.Script); err != nil {
			log.Printf("error executing cron job: %s (%e)\n", job.Name, err)
		}
	}

	// if it was not possible
	if err := register(int(job.ID), job.Expression, scriptHandler); err != nil {
		Del(int(job.ID))
		return err
	}

	return nil
}

func List() []Cron {
	crons := []Cron{}

	entries, err := db.Queries.ListCronJobs(db.Ctx)
	if err != nil {
		log.Println("error loading persisted cron jobs: ", err)
		return crons
	}

	for _, entry := range entries {
		status := "not running"
		if _, ok := runningJobs[int(entry.ID)]; ok {
			status = "running"
		}

		crons = append(crons, Cron{
			entry,
			status,
		})
	}

	return crons
}

func Del(id int) {
	// remove from db
	err := db.Queries.DeleteCronJob(db.Ctx, int64(id))
	if err != nil {
		log.Println("error deleting cron job: ", err)
		return
	}

	// if its running, stop it
	if entryID, ok := runningJobs[id]; ok {
		delete(runningJobs, id)
		scheduler.Remove(entryID)
	}
}
