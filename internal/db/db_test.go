package db_test

import (
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/db"
	"github.com/marco-souza/marco.fly.dev/internal/db/sqlc"
	"github.com/marco-souza/marco.fly.dev/internal/di"
	"github.com/stretchr/testify/assert"
)

func TestDbClient(t *testing.T) {
	var id int64
	cronToInsert := sqlc.CreateCronJobParams{
		Name:       "test",
		Expression: "* * * * *",
		Script:     "echo test",
	}

	di.Injectable(config.Config{SqliteUrl: ":memory:"})

	ds := db.New()
	err := ds.Start()
	assert.NoError(t, err)

	defer ds.Stop()

	t.Run("create cronjob", func(t *testing.T) {
		// create cronjob
		insertedCron, err := ds.Queries.CreateCronJob(ds.Ctx, cronToInsert)

		assert.NoError(t, err)
		assert.NotNil(t, insertedCron)

		// assert CreateCronJobParams
		assert.Equal(t, cronToInsert.Name, insertedCron.Name)
		assert.Equal(t, cronToInsert.Expression, insertedCron.Expression)
		assert.Equal(t, cronToInsert.Script, insertedCron.Script)

		id = insertedCron.ID
	})

	t.Run("fetch cronjob", func(t *testing.T) {
		// fetch inserted cron
		cron, err := ds.Queries.GetCronJob(ds.Ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, cron)

		assert.Equal(t, cronToInsert.Name, cron.Name)
		assert.Equal(t, cronToInsert.Expression, cron.Expression)
		assert.Equal(t, cronToInsert.Script, cron.Script)
	})

	t.Run("update cronjob", func(t *testing.T) {
		// update cron
		newName := "new name"
		cronUpdatePayload := sqlc.UpdateCronJobParams{
			ID:         id,
			Name:       newName,
			Script:     cronToInsert.Script,
			Expression: cronToInsert.Expression,
		}

		err = ds.Queries.UpdateCronJob(ds.Ctx, cronUpdatePayload)
		assert.NoError(t, err)

		// fetch inserted cron
		cron, err := ds.Queries.GetCronJob(ds.Ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, cron)

		assert.Equal(t, cronUpdatePayload.Name, cron.Name)
		assert.Equal(t, cronUpdatePayload.Script, cron.Script)
	})

	t.Run("list cronjobs", func(t *testing.T) {
		// listing crons
		crons, err := ds.Queries.ListCronJobs(ds.Ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, crons)

		// assert ListCronJobs
		assert.Len(t, crons, 1)
		assert.Equal(t, crons[0].ID, id)
		assert.Equal(t, crons[0].Expression, cronToInsert.Expression)
		assert.Equal(t, crons[0].Script, cronToInsert.Script)
	})
}
