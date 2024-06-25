package main

import (
	"fmt"

	"github.com/marco-souza/marco.fly.dev/internal/db"
	"github.com/marco-souza/marco.fly.dev/internal/db/sqlc"
)

func main() {
	if err := db.Init("./sqlite.db"); err != nil {
		fmt.Println("error creating cronjob: ", db.Ctx, err)
		panic(err)
	}
	defer db.Close()

	client := db.Queries

	// create cronjob
	cron := sqlc.CreateCronJobParams{
		Name:       "test",
		Expression: "* * * * *",
		Script:     "echo test",
	}

	insertedCron, err := client.CreateCronJob(db.Ctx, cron)
	if err != nil {
		fmt.Println("error creating cronjob: ", db.Ctx, err)
		panic(err)
	}

	fmt.Println("insertedCron", insertedCron)

	// fetch inserted cron
	fetchedCron, err := client.GetCronJob(db.Ctx, insertedCron.ID)
	if err != nil {
		fmt.Println("error fetching cronjob: ", db.Ctx, err)
		panic(err)
	}

	fmt.Println("fetchedCron", fetchedCron)

	// update cron
	cronUpdatePayload := sqlc.UpdateCronJobParams{
		ID:   fetchedCron.ID,
		Name: "test2",
	}

	if err = client.UpdateCronJob(db.Ctx, cronUpdatePayload); err != nil {
		fmt.Println("error fetching cronjob: ", db.Ctx, err)
		panic(err)
	}

	// listing crons
	crons, err := client.ListCronJobs(db.Ctx)
	if err != nil {
		fmt.Println("error listing cronjobs: ", db.Ctx, err)
		panic(err)
	}

	fmt.Println("ListCronJobs", crons)
}
