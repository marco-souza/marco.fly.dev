// service to manage cron jobs
package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

type expressionMap map[int]string

type cronService struct {
	scheduler             *cron.Cron
	registeredExpressions expressionMap
}

type Cron struct {
	ID         int
	Expression string
}

func new() *cronService {
	scheduler := cron.New()
	exprMap := expressionMap{}

	return &cronService{
		scheduler,
		exprMap,
	}
}

func (c *cronService) Start() error {
	log.Println("starting scheduler")
	c.scheduler.Start()
	return registerPersistedJobs()
}

func (c *cronService) Stop() {
	log.Println("stopping scheduler")
	c.scheduler.Stop()
}

func (c *cronService) Add(cronExpr string, handler func()) error {
	entryID, err := c.scheduler.AddFunc(cronExpr, handler)
	if err != nil {
		return err
	}

	c.registeredExpressions[int(entryID)] = cronExpr
	log.Println("cron job created: ", entryID)
	return nil
}

func (c *cronService) List() []Cron {
	entries := c.scheduler.Entries()

	crons := []Cron{}
	for _, entry := range entries {
		a := c.scheduler.Entry(entry.ID)
		log.Println(a)

		crons = append(crons, Cron{
			ID:         int(entry.ID),
			Expression: c.registeredExpressions[int(entry.ID)],
		})
	}

	return crons
}

func (c *cronService) Del(id int) {
	c.scheduler.Remove(cron.EntryID(id))
	delete(c.registeredExpressions, id)
}

var CronService = new()
