package logs

import (
	"context"

	"github.com/jackc/pgx/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rickCrz7/Inventory-API/utils"
	log "github.com/sirupsen/logrus"
)

type Dao struct{}

func NewDao() *Dao {
	return &Dao{}
}

func (d *Dao) GetLogs(ctx context.Context, tx pgx.Tx, device_id string) ([]*utils.DeviceLog, error) {
	log.Printf("Fetching logs for Device ID: %s", device_id)

	rows, err := tx.Query(ctx, `
		SELECT id, device_id, log_type, note, created_at, created_by FROM device_logs WHERE device_id = $1
	`, device_id)
	if err != nil {
		log.Errorf("Error fetching logs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var logs []*utils.DeviceLog
	for rows.Next() {
		var logEntry utils.DeviceLog
		if err := rows.Scan(&logEntry.ID, &logEntry.DeviceID, &logEntry.LogType, &logEntry.Note, &logEntry.CreatedAt, &logEntry.CreatedBy); err != nil {
			log.Errorf("Error scanning log entry: %v", err)
			continue
		}
		logs = append(logs, &logEntry)
	}
	if err := rows.Err(); err != nil {
		log.Errorf("Error iterating over log rows: %v", err)
		return nil, err
	}
	return logs, nil
}

func (d *Dao) CreateLog(ctx context.Context, tx pgx.Tx, logEntry *utils.DeviceLog) error {
	log.Printf("Creating log for Device ID: %s", logEntry.DeviceID)

	if logEntry.ID == "" {
		var err error
		logEntry.ID, err = gonanoid.New()
		if err != nil {
			log.Errorf("Error generating log ID: %v", err)
			return err
		}
	}

	_, err := tx.Exec(ctx, `
		INSERT INTO device_logs (id, device_id, log_type, note, created_at, created_by) VALUES ($1, $2, $3, $4, $5, $6)
	`, logEntry.ID, logEntry.DeviceID, logEntry.LogType, logEntry.Note, logEntry.CreatedAt, logEntry.CreatedBy)
	if err != nil {
		log.Errorf("Error creating log entry: %v", err)
		return err
	}
	return nil
}

func (d *Dao) DeleteLog(ctx context.Context, tx pgx.Tx, id string) error {
	log.Printf("Deleting log with ID: %s", id)

	_, err := tx.Exec(ctx, `
		DELETE FROM device_logs WHERE id = $1
	`, id)
	if err != nil {
		log.Errorf("Error deleting log entry: %v", err)
		return err
	}
	return nil
}
