package devices

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
	viper.AddConfigPath("../config")
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

func TestGetDevice(t *testing.T) {
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
	t.Run("GetDevice", func(t *testing.T) {
		device, err := dao.GetDevice(ctx, tx, d_id)
		if err != nil {
			t.Fatalf("Error getting device: %v", err)
		}
		if device == nil {
			t.Fatal("Expected device to be found")
		}
		if device.ID != d_id {
			t.Fatalf("Expected device ID to be %s, got %s", d_id, device.ID)
		}
	})
	t.Run("GetDeviceNotFound", func(t *testing.T) {
		device, err := dao.GetDevice(ctx, tx, "nonexistent-id")
		if err == nil {
			t.Fatal("Expected error getting device, got nil")
		}
		if device != nil {
			t.Fatal("Expected device to be nil")
		}
	})
}

func TestGetDevices(t *testing.T) {
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
	for i := 0; i < 5; i++ {
		d_id, _ := gonanoid.New()
		_, err = tx.Exec(ctx, `
			INSERT INTO devices (id, serial_number, name, purchase_date, status, owner_id, type_id) VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, d_id, fmt.Sprintf("SN12345%d", i), fmt.Sprintf("Test Device %d", i), "2023-01-01", "active", id, t_id)
		if err != nil {
			t.Fatalf("Error inserting mock data: %v", err)
		}
	}
	// Test Cases
	t.Run("GetDevices", func(t *testing.T) {
		devices, err := dao.GetDevices(ctx, tx)
		if err != nil {
			t.Fatalf("Error getting devices: %v", err)
		}
		if len(devices) < 5 {
			t.Fatalf("Expected at least 5 devices, got %d", len(devices))
		}
	})
}

func TestCreateDevice(t *testing.T) {
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
	t.Run("CreateDevice", func(t *testing.T) {
		purchaseDate, err := time.Parse("2006-01-02", "2023-01-01")
		if err != nil {
			t.Fatalf("Error parsing purchase date: %v", err)
		}
		device := &utils.Device{
			SerialNumber: "SN123456",
			Name:         "Test Device",
			PurchaseDate: purchaseDate,
			Status:       "active",
			OwnerID:      id,
			TypeID:       t_id,
		}
		err = dao.CreateDevice(ctx, tx, device)
		if err != nil {
			t.Fatalf("Error creating device: %v", err)
		}
		row := tx.QueryRow(ctx, `
			SELECT id, serial_number, name, purchase_date, status, owner_id, type_id FROM devices WHERE serial_number = $1
		`, "SN123456")
		var created_device utils.Device
		err = row.Scan(&created_device.ID, &created_device.SerialNumber, &created_device.Name, &created_device.PurchaseDate, &created_device.Status, &created_device.OwnerID, &created_device.TypeID)
		if err != nil {
			t.Fatalf("Error getting created device: %v", err)
		}
		if created_device.Name != "Test Device" {
			t.Fatalf("Expected device name to be 'Test Device', got %q", created_device.Name)
		}
	})
	t.Run("CreateDeviceWithExistingID", func(t *testing.T) {
		purchaseDate, err := time.Parse("2006-01-02", "2023-01-01")
		if err != nil {
			t.Fatalf("Error parsing purchase date: %v", err)
		}
		device := &utils.Device{
			ID:           "existing-id",
			SerialNumber: "SN123456",
			Name:         "Test Device",
			PurchaseDate: purchaseDate,
			Status:       "active",
			OwnerID:      id,
			TypeID:       t_id,
		}
		err = dao.CreateDevice(ctx, tx, device)
		if err != nil {
			t.Fatalf("Error creating device: %v", err)
		}
		created_device, err := dao.GetDevice(ctx, tx, "existing-id")
		if err != nil {
			t.Fatalf("Error getting created device: %v", err)
		}
		if created_device.SerialNumber != "SN123456" {
			t.Fatalf("Expected device serial number to be 'SN123456', got %q", created_device.SerialNumber)
		}
	})
}

func TestUpdateDevice(t *testing.T) {
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
	t.Run("UpdateDevice", func(t *testing.T) {
		purchaseDate, err := time.Parse("2006-01-02", "2023-01-01")
		if err != nil {
			t.Fatalf("Error parsing purchase date: %v", err)
		}
		device := &utils.Device{
			ID:           d_id,
			SerialNumber: "SN123456",
			Name:         "Test Device",
			PurchaseDate: purchaseDate,
			Status:       "active",
			OwnerID:      id,
			TypeID:       t_id,
		}
		err = dao.UpdateDevice(ctx, tx, device)
		if err != nil {
			t.Fatalf("Error updating device: %v", err)
		}
		updatedDevice, err := dao.GetDevice(ctx, tx, d_id)
		if err != nil {
			t.Fatalf("Error getting updated device: %v", err)
		}
		if updatedDevice.Name != "Test Device" {
			t.Fatalf("Expected device name to be 'Test Device', got %q", updatedDevice.Name)
		}
	})
}

func TestDeleteDevice(t *testing.T) {
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
	t.Run("DeleteDevice", func(t *testing.T) {
		err = dao.DeleteDevice(ctx, tx, d_id)
		if err != nil {
			t.Fatalf("Error deleting device: %v", err)
		}
		_, err = dao.GetDevice(ctx, tx, d_id)
		if err == nil {
			t.Fatalf("Expected error getting deleted device, got nil")
		}
	})
}