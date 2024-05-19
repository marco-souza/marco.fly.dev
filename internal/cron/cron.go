// service to manage cron jobs
package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

type CronService struct {
	scheduler *cron.Cron
}

func New() *CronService {
	scheduler := cron.New()
	return &CronService{
		scheduler,
	}
}

func (c *CronService) Start() {
	log.Println("starting scheduler")
	c.scheduler.Start()
}

func (c *CronService) Stop() {
	log.Println("stopping scheduler")
	c.scheduler.Stop()
}

func (c *CronService) Add(cronExpr string, handler func()) error {
	entryID, err := c.scheduler.AddFunc(cronExpr, handler)
	if err != nil {
		return err
	}

	log.Println("cron job created: ", entryID)
	return nil
}

type Cron struct {
	id         int
	expression string // TODO: get the persisted expression
}

func (c *CronService) List() []Cron {
	entries := c.scheduler.Entries()

	crons := []Cron{}
	for _, entry := range entries {
		a := c.scheduler.Entry(entry.ID)
		log.Println(a)

		crons = append(crons, Cron{
			id: int(entry.ID),
		})
	}

	return crons
}

func (c *CronService) Del(id int) {
	c.scheduler.Remove(cron.EntryID(id))
}
