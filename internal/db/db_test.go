package db_test

import (
	"database/sql"
	"testing"
	"time"

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

	date := time.Now().Truncate(time.Second)

	// create test for funancial_logs
	t.Run("create financial log", func(t *testing.T) {
		// create financial log
		financialLogToInsert := sqlc.CreateFinancialLogParams{
			Amount:      100,
			Currency:    "BRL",
			Investiment: "test",
			CreatedAt:   sql.NullTime{Time: date, Valid: true},
		}

		insertedFinancialLog, err := ds.Queries.CreateFinancialLog(ds.Ctx, financialLogToInsert)

		assert.NoError(t, err)
		assert.NotNil(t, insertedFinancialLog)

		// assert CreateFinancialLogParams
		assert.Equal(t, financialLogToInsert.Amount, insertedFinancialLog.Amount)
		assert.Equal(t, financialLogToInsert.Currency, insertedFinancialLog.Currency)
		assert.Equal(t, financialLogToInsert.Investiment, insertedFinancialLog.Investiment)
		assert.Equal(t, insertedFinancialLog.ID, int64(1))
		assert.Equal(t, insertedFinancialLog.CreatedAt.Time.UTC(), date.UTC())
	})

	t.Run("fetch financial log", func(t *testing.T) {
		// fetch inserted financial log
		financialLog, err := ds.Queries.GetFinancialLog(ds.Ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, financialLog)

		assert.Equal(t, financialLog.Amount, int64(100))
		assert.Equal(t, financialLog.Currency, "BRL")
		assert.Equal(t, financialLog.Investiment, "test")
		assert.Equal(t, financialLog.ID, int64(1))
	})

	t.Run("update financial log", func(t *testing.T) {
		// update financial log
		newAmount := int64(200)
		financialLogUpdatePayload := sqlc.UpdateFinancialLogParams{
			ID:          1,
			Amount:      newAmount,
			Currency:    "BRL",
			Investiment: "test",
		}

		financialLog, err := ds.Queries.UpdateFinancialLog(ds.Ctx, financialLogUpdatePayload)
		assert.NoError(t, err)

		assert.NotNil(t, financialLog)
		assert.Equal(t, financialLog.Amount, newAmount)
	})

	t.Run("list financial logs", func(t *testing.T) {
		// listing financial logs
		financialLogs, err := ds.Queries.ListFinancialLogs(ds.Ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, financialLogs)

		// assert ListFinancialLogs
		assert.Len(t, financialLogs, 1)
		assert.Equal(t, financialLogs[0].ID, int64(1))
		assert.Equal(t, financialLogs[0].Amount, int64(200))
		assert.Equal(t, financialLogs[0].Currency, "BRL")
		assert.Equal(t, financialLogs[0].Investiment, "test")
	})

	t.Run("delete financial log", func(t *testing.T) {
		// delete financial log
		financialLog, err := ds.Queries.DeleteFinancialLog(ds.Ctx, 1)
		assert.NoError(t, err)

		_, err = ds.Queries.GetFinancialLog(ds.Ctx, financialLog.ID)
		assert.Error(t, err)
	})

}
