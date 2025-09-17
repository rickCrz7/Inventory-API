package logs

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"testing"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rickCrz7/Inventory-API/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var postgresURI string

func init() {
	viper.SetConfigName("app")
	viper.AddConfigPath("../../config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	postgresURI = viper.GetString("postgres.dev") // Change to "postgres.dev" for development/local db
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("\t%s:%d", filename, f.Line)
		},
	})
	log.SetLevel(log.DebugLevel)
}

func TestGetLogs(t *testing.T) {
	pdb, err := utils.OpenDB(postgresURI, false)
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer pdb.Close()

	ctx := context.Background()
	tx, err := pdb.Begin(ctx)
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	dao := NewDao()
	// Mock Data
	id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO owners (id, first_name, last_name, email)
		VALUES ($1, $2, $3, $4)
	`, id, "John", "Doe", "john.doe@example.com")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	t_id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO types (id, name, description) VALUES ($1, $2, $3)
	`, t_id, "Test Type", "This is a test type")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	d_id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO devices (id, serial_number, name, purchase_date, status, owner_id, type_id) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, d_id, "SN123456", "Test Device", "2023-01-01", "active", id, t_id)
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	// Add logs
	for i := 0; i < 5; i++ {
		logID, _ := gonanoid.New()
		_, err = tx.Exec(ctx, `
			INSERT INTO device_logs (id, device_id, note, log_type, created_at, created_by)
			VALUES ($1, $2, $3, $4, NOW(), $5)
		`, logID, d_id, fmt.Sprintf("Log message %d", i+1), "INFO", "tester")
		if err != nil {
			t.Fatalf("Error inserting log: %v", err)
		}
	}

	t.Run("GetLogs", func(t *testing.T) {
		logs, err := dao.GetLogs(ctx, tx, d_id)
		if err != nil {
			t.Fatalf("Error getting logs: %v", err)
		}
		if logs == nil {
			t.Fatal("Expected logs to be found")
		}
		if len(logs) != 5 {
			t.Fatalf("Expected 5 logs, got %d", len(logs))
		}
	})

	t.Run("GetLogsNotFound", func(t *testing.T) {
		logs, err := dao.GetLogs(ctx, tx, "nonexistent-id")
		if err != nil {
			t.Fatal("Expected no error, got:", err)
		}
		if logs != nil {
			t.Fatal("Expected logs to be nil")
		}
	})
}

func TestCreateLogs(t *testing.T) {
	pdb, err := utils.OpenDB(postgresURI, false)
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer pdb.Close()

	ctx := context.Background()
	tx, err := pdb.Begin(ctx)
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	dao := NewDao()
	// Mock Data
	id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO owners (id, first_name, last_name, email)
		VALUES ($1, $2, $3, $4)
	`, id, "John", "Doe", "john.doe@example.com")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	t_id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO types (id, name, description) VALUES ($1, $2, $3)
	`, t_id, "Test Type", "This is a test type")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	d_id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO devices (id, serial_number, name, purchase_date, status, owner_id, type_id) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, d_id, "SN123456", "Test Device", "2023-01-01", "active", id, t_id)
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}

	t.Run("CreateLog", func(t *testing.T) {
		logEntry := &utils.DeviceLog{
			DeviceID:  d_id,
			LogType:   "INFO",
			Note:      "Device created",
			CreatedAt: time.Now(),
			CreatedBy: "tester",
		}
		err := dao.CreateLog(ctx, tx, logEntry)
		if err != nil {
			t.Fatalf("Error creating log: %v", err)
		}
		// If CreateLog generates the ID, retrieve it for the test
		if logEntry.ID == "" {
			row := tx.QueryRow(ctx, `SELECT id FROM device_logs WHERE device_id = $1 AND log_type = $2 AND note = $3 AND created_by = $4 ORDER BY created_at DESC LIMIT 1`, d_id, "INFO", "Device created", "tester")
			err = row.Scan(&logEntry.ID)
			if err != nil {
				t.Fatalf("Could not retrieve inserted log ID: %v", err)
			}
		}
		// Check if log was inserted
		row := tx.QueryRow(ctx, `SELECT id, device_id, log_type, note, created_at, created_by FROM device_logs WHERE id = $1`, logEntry.ID)
		var got utils.DeviceLog
		err = row.Scan(&got.ID, &got.DeviceID, &got.LogType, &got.Note, &got.CreatedAt, &got.CreatedBy)
		if err != nil {
			t.Fatalf("Inserted log not found: %v", err)
		}
		if got.DeviceID != d_id || got.LogType != "INFO" || got.Note != "Device created" || got.CreatedBy != "tester" {
			t.Fatalf("Inserted log does not match: %+v", got)
		}
	})

	t.Run("CreateLogWithID", func(t *testing.T) {
		customID, _ := gonanoid.New()
		logEntry := &utils.DeviceLog{
			ID:        customID,
			DeviceID:  d_id,
			LogType:   "WARN",
			Note:      "Custom ID log",
			CreatedAt: time.Now(),
			CreatedBy: "tester2",
		}
		err := dao.CreateLog(ctx, tx, logEntry)
		if err != nil {
			t.Fatalf("Error creating log with custom ID: %v", err)
		}
		row := tx.QueryRow(ctx, `SELECT id FROM device_logs WHERE id = $1`, customID)
		var foundID string
		err = row.Scan(&foundID)
		if err != nil {
			t.Fatalf("Log with custom ID not found: %v", err)
		}
		if foundID != customID {
			t.Fatalf("Expected ID %s, got %s", customID, foundID)
		}
	})

	t.Run("CreateLogInvalidDeviceID", func(t *testing.T) {
		logEntry := &utils.DeviceLog{
			DeviceID:  "nonexistent-device",
			LogType:   "ERROR",
			Note:      "Invalid device",
			CreatedAt: time.Now(),
			CreatedBy: "tester3",
		}
		err := dao.CreateLog(ctx, tx, logEntry)
		if err == nil {
			t.Fatal("Expected error when creating log with invalid device_id, got nil")
		}
	})
}

func TestDeleteLogs(t *testing.T) {
	pdb, err := utils.OpenDB(postgresURI, false)
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer pdb.Close()

	ctx := context.Background()
	tx, err := pdb.Begin(ctx)
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	dao := NewDao()
	// Mock Data
	id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO owners (id, first_name, last_name, email)
		VALUES ($1, $2, $3, $4)
	`, id, "John", "Doe", "john.doe@example.com")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	t_id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO types (id, name, description) VALUES ($1, $2, $3)
	`, t_id, "Test Type", "This is a test type")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	d_id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO devices (id, serial_number, name, purchase_date, status, owner_id, type_id) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, d_id, "SN123456", "Test Device", "2023-01-01", "active", id, t_id)
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	logID, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
			INSERT INTO device_logs (id, device_id, note, log_type, created_at, created_by)
			VALUES ($1, $2, $3, $4, NOW(), $5)
		`, logID, d_id, fmt.Sprintf("Log message %d", 1), "INFO", "tester")
	if err != nil {
		t.Fatalf("Error inserting log: %v", err)
	}

	t.Run("DeleteLog", func(t *testing.T) {
		err := dao.DeleteLog(ctx, tx, logID)
		if err != nil {
			t.Fatalf("Error deleting log: %v", err)
		}
		// Verify deletion
		row := tx.QueryRow(ctx, `SELECT id FROM device_logs WHERE id = $1`, logID)
		var foundID string
		err = row.Scan(&foundID)
		if err == nil {
			t.Fatal("Expected no rows after deletion, but found one")
		}
	})
}
