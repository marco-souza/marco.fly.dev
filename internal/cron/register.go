// This file contains the scheduler package which is responsible for scheduling the tasks to be executed.
package cron

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/marco-souza/marco.fly.dev/internal/db"
	"github.com/marco-souza/marco.fly.dev/internal/lua"
)

var (
	logPrefix = "cronjob: [main]: "
	logger    = log.New(log.Writer(), logPrefix, log.Flags())
)

// Setup scheduler by initializing cronjobs, registering lua scripts persisted
func registerPersistedJobs() error {
	logger.Println("loading persisted cron jobs")
	crons, err := db.Queries.ListCronJobs(db.Ctx)
	if err != nil {
		logger.Println("error loading persisted cron jobs: ", err)
		return err
	}

	if len(crons) == 0 {
		logger.Println("no job found")
		return nil
	}

	logger.Println("setup persisted cron jobs: ", len(crons))
	for _, c := range crons {
		jobPrefix := fmt.Sprintf("cronjob: [%d]: ", c.ID)
		logger.SetPrefix(jobPrefix)

		cronHandler := func() {
			logger.Printf("executing cron job: %s\n", c.Name)

			if _, err := lua.Run(c.Script); err != nil {
				logger.Printf("error executing cron job: %s (%e)\n", c.Name, err)
			}
		}

		if err := register(int(c.ID), c.Expression, cronHandler); err != nil {
			logger.Printf("error adding cron job: %s (%e)\n", c.Name, err)
			return err
		}
	}
	logger.SetPrefix(logPrefix)
	logger.Println("setup  local cron jobs")

	registerLocalScripts("scripts")

	return nil
}

func register(id int, cronExpr string, handler func()) error {
	entryID, err := scheduler.AddFunc(cronExpr, handler)
	if err != nil {
		return err
	}

	runningJobs[id] = entryID
	logger.Println("cron job registered: ", entryID)

	return nil
}

func registerLocalScripts(scriptFolder string) {
	logger.Println("loading local cron jobs")

	localCronJobs, err := os.ReadDir(scriptFolder)
	if err != nil {
		logger.Println("error loading local cron jobs: ", err)
		return
	}

	fileCounter := 0
	for _, f := range localCronJobs {
		// ignore any file that doesn't end with .lua
		if f.IsDir() || filepath.Ext(f.Name()) != ".lua" || f.Name()[:1] == "_" {
			logger.Println("ignoring file: ", f.Name())
			continue
		}

		name := f.Name()
		rawFile, err := os.ReadFile(filepath.Join(scriptFolder, name))
		if err != nil {
			logger.Printf("error reading cron job: %s (%e)\n", name, err)
			continue
		}

		script := string(rawFile)

		firstLine := strings.Split(string(script), "\n")[0]
		cronExpr := strings.TrimSpace(firstLine)[len("--cron: "):] // ignore '--cron: '

		fileCounter++
		logger.Printf("registering cronjob: cron:%s name:%s", cronExpr, name)

		baseInt := 10000 // offset to avoid conflict with persisted jobs
		localID := baseInt + fileCounter
		register(localID, cronExpr, func() {
			logger.Printf("executing cron job: %s\n", name)
			lua.Run(script) // ignore error
		})
	}
}
