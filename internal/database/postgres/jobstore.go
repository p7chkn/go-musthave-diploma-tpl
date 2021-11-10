package postgres

import (
	"context"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/utils"
	"time"
)

func (db *PostgreDataBase) AddJob(ctx context.Context, job models.JobStoreRow) error {
	sqlCreateJob := `INSERT INTO jobstore (type, next_time_execute, parameters, count) VALUES ($1, $2, $3, $4)`
	_, err := db.conn.ExecContext(ctx, sqlCreateJob, job.Type, job.NextTimeExecute, job.Parameters, job.Count)
	return err
}

func (db *PostgreDataBase) GetJobToExecute(ctx context.Context, maxCount int) ([]models.JobStoreRow, error) {
	var result []models.JobStoreRow
	now := time.Now()
	sqlGetJobToExecute := `SELECT id, type, next_time_execute, parameters, count, executed FROM jobstore 
						   WHERE count < $1 AND next_time_execute <= $2 AND executed = FALSE`
	rows, err := db.conn.QueryContext(ctx, sqlGetJobToExecute, maxCount, now)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var jsRow models.JobStoreRow
		if err = rows.Scan(&jsRow.ID, &jsRow.Type, &jsRow.NextTimeExecute, &jsRow.Parameters,
			&jsRow.Count, &jsRow.Executed); err != nil {
			return result, nil
		}
		result = append(result, jsRow)

	}
	return result, nil
}

func (db *PostgreDataBase) ExecuteJob(ctx context.Context, jobID string) error {
	sqlExecuteJob := `UPDATE jobstore SET executed = TRUE WHERE id = $1`
	_, err := db.conn.ExecContext(ctx, sqlExecuteJob, jobID)
	return err
}

func (db *PostgreDataBase) IncreaseCounter(ctx context.Context, jobID string, count int) error {
	timeExecute := time.Now().Add(utils.CalculateAdditionTime(count))
	sqlIncreaseCounter := `UPDATE jobstore SET count = count + 1, next_time_execute = $1 WHERE id = $2`
	_, err := db.conn.ExecContext(ctx, sqlIncreaseCounter, timeExecute, jobID)
	return err
}
