// This file contains the scheduler package which is responsible for scheduling the tasks to be executed.
package cron

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/marco-souza/marco.fly.dev/internal/lua"
)

// Setup scheduler by initializing cronjobs, registering lua scripts persisted
func (tss *TaskScheduleService) registerPersistedJobs() error {
	logger.Info("loading persisted cron jobs", "ctx", tss.db.Ctx, "q", tss.db.Queries)
	crons, err := tss.db.Queries.ListCronJobs(tss.db.Ctx)
	if err != nil {
		logger.Error("error loading persisted cron jobs", "err", err)
		return err
	}

	if len(crons) == 0 {
		logger.Info("no job found")
		return nil
	}

	logger.Info("setup persisted cron jobs: ", "length", len(crons))

	for _, c := range crons {
		logger := logger.With("id", c.ID)

		cronHandler := func() {
			logger.Info("executing cron job", "name", c.Name)

			if _, err := lua.Run(c.Script); err != nil {
				logger.Error("error executing cron job", "name", c.Name, "err", err)
			}
		}

		if err := tss.register(int(c.ID), c.Expression, cronHandler); err != nil {
			logger.Warn("error adding cron job", "name", c.Name, "err", err)
			continue
		}
	}

	return nil
}

func (tss *TaskScheduleService) register(id int, cronExpr string, handler func()) error {
	entryID, err := tss.scheduler.AddFunc(cronExpr, handler)
	if err != nil {
		return err
	}

	tss.runningJobs[id] = entryID
	logger.Info("cron job registered", "id", entryID)

	return nil
}

func (tss *TaskScheduleService) registerLocalJobs(scriptFolder string) error {
	logger.Info("loading local cron jobs")

	localCronJobs, err := os.ReadDir(scriptFolder)
	if err != nil {
		return err
	}

	fileCounter := 0
	for _, f := range localCronJobs {
		// ignore any file that doesn't end with .lua
		if f.IsDir() || filepath.Ext(f.Name()) != ".lua" || f.Name()[:1] == "_" {
			logger.Debug("ignoring file", "name", f.Name())
			continue
		}

		name := f.Name()
		rawFile, err := os.ReadFile(filepath.Join(scriptFolder, name))
		if err != nil {
			logger.Warn("error reading cron job", "name", name, "err", err)
			continue
		}

		script := string(rawFile)

		firstLine := strings.Split(string(script), "\n")[0]
		cronExpr := strings.TrimSpace(firstLine)[len("--cron: "):] // ignore '--cron: '

		fileCounter++
		logger.Info("registering cronjob", "expression", cronExpr, "name", name)

		baseInt := 10000 // offset to avoid conflict with persisted jobs
		localID := baseInt + fileCounter

		if err := tss.register(localID, cronExpr, func() {
			logger.Info("executing cron job", "name", name)
			lua.Run(script) // ignore error
		}); err != nil {
			baseInt--
			logger.Warn("error registering cron job", "name", name, "err", err)
			continue
		}
	}

	return nil
}
