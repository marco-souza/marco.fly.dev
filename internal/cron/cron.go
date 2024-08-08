// service to manage cron jobs
package cron

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/marco-souza/marco.fly.dev/internal/db"
	"github.com/marco-souza/marco.fly.dev/internal/db/sqlc"
	"github.com/marco-souza/marco.fly.dev/internal/di"
	"github.com/marco-souza/marco.fly.dev/internal/lua"
	"github.com/robfig/cron/v3"
)

type runningCronJobs map[int]cron.EntryID

type Cron struct {
	sqlc.Cron
	status string
}

var logger = slog.With("service", "cron")

type TaskScheduleService struct {
	scheduler   *cron.Cron
	runningJobs runningCronJobs
	db          *db.DatabaseService
}

func New() *TaskScheduleService {
	dbs := di.MustInject(db.DatabaseService{})
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(fmt.Errorf("error creating cron scheduler: %w", err))
	}

	scheduler := cron.New(cron.WithLocation(location))

	return &TaskScheduleService{
		runningJobs: runningCronJobs{},
		scheduler:   scheduler,
		db:          dbs,
	}
}

func (tss *TaskScheduleService) Start() error {
	tss.scheduler.Start()

	if err := tss.registerPersistedJobs(); err != nil {
		tss.scheduler.Stop()
		return err
	}

	if err := tss.registerLocalJobs("scripts"); err != nil {
		tss.scheduler.Stop()
		return err
	}

	return nil
}

func (tss *TaskScheduleService) Stop() error {
	tss.scheduler.Stop()
	return nil
}

func (tss *TaskScheduleService) AddScript(name, cronExpr, script string) error {
	dbs := di.MustInject(db.DatabaseService{})
	job, err := dbs.Queries.CreateCronJob(dbs.Ctx, sqlc.CreateCronJobParams{
		Name:       name,
		Expression: cronExpr,
		Script:     script,
	})
	if err != nil {
		return err
	}

	// register cron job
	scriptHandler := func() {
		logger.Info("executing cron job", "name", job.Name)

		if _, err := lua.Run(job.Script); err != nil {
			logger.Error("error executing cron job", "name", job.Name, "err", err)
		}
	}

	// if it was not possible
	if err := tss.register(int(job.ID), job.Expression, scriptHandler); err != nil {
		tss.Del(int(job.ID))
		return err
	}

	return nil
}

func (tss *TaskScheduleService) List() []Cron {
	crons := []Cron{}

	entries, err := tss.db.Queries.ListCronJobs(tss.db.Ctx)
	if err != nil {
		logger.Warn("error loading persisted cron jobs", "err", err)
		return crons
	}

	for _, entry := range entries {
		status := "not running"
		if _, ok := tss.runningJobs[int(entry.ID)]; ok {
			status = "running"
		}

		crons = append(crons, Cron{
			entry,
			status,
		})
	}

	return crons
}

func (tss *TaskScheduleService) Del(id int) {
	// remove from db
	err := tss.db.Queries.DeleteCronJob(tss.db.Ctx, int64(id))
	if err != nil {
		logger.Warn("error deleting cron job", "err", err)
		return
	}

	// if its running, stop it
	if entryID, ok := tss.runningJobs[id]; ok {
		delete(tss.runningJobs, id)
		tss.scheduler.Remove(entryID)
	}
}
