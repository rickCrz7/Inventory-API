package properties

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"testing"

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

func TestGetProperties(t *testing.T) {
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
	// Insert a type_property for "color"
	typePropertyID, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO type_properties (id, type_id, name, data_type, required)
		VALUES ($1, $2, $3, $4, $5)
	`, typePropertyID, t_id, "color", "string", true)
	if err != nil {
		t.Fatalf("Error inserting type_property: %v", err)
	}
	p_id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO device_properties (id, device_id, type_property_id, value) VALUES ($1, $2, $3, $4)
		`, p_id, d_id, typePropertyID, "red")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}

	// Add another type_property for "weight"
	typePropertyID2, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO type_properties (id, type_id, name, data_type, required)
		VALUES ($1, $2, $3, $4, $5)
	`, typePropertyID2, t_id, "weight", "float", true)
	if err != nil {
		t.Fatalf("Error inserting type_property: %v", err)
	}
	p_id2, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO device_properties (id, device_id, type_property_id, value) VALUES ($1, $2, $3, $4)
		`, p_id2, d_id, typePropertyID2, "2.5")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}

	// Add another type_property for "is_wireless"
	typePropertyID3, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO type_properties (id, type_id, name, data_type, required)
		VALUES ($1, $2, $3, $4, $5)
	`, typePropertyID3, t_id, "is_wireless", "bool", true)
	if err != nil {
		t.Fatalf("Error inserting type_property: %v", err)
	}
	p_id3, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO device_properties (id, device_id, type_property_id, value) VALUES ($1, $2, $3, $4)
		`, p_id3, d_id, typePropertyID3, "true")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}

	t.Run("GetProperties", func(t *testing.T) {
		properties, err := dao.GetProperties(ctx, tx, d_id)
		if err != nil {
			t.Fatalf("Error getting properties: %v", err)
		}
		if properties == nil {
			t.Fatal("Expected properties to be found")
		}
		if len(properties) != 3 {
			t.Fatalf("Expected 3 properties, got %d", len(properties))
		}
	})
	t.Run("GetDeviceNotFound", func(t *testing.T) {
		properties, err := dao.GetProperties(ctx, tx, "nonexistent-id")
		if err != nil {
			t.Fatal("Expected no error, got:", err)
		}
		if properties != nil {
			t.Fatal("Expected properties to be nil")
		}
	})
}

func TestCreateProperty(t *testing.T) {
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
	// Insert a type_property for "color"
	typePropertyID, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO type_properties (id, type_id, name, data_type, required)
		VALUES ($1, $2, $3, $4, $5)
	`, typePropertyID, t_id, "color", "string", true)
	if err != nil {
		t.Fatalf("Error inserting type_property: %v", err)
	}

	t.Run("CreateProperty", func(t *testing.T) {
		property := &utils.DeviceProperty{
			DeviceID:       d_id,
			TypePropertyID: typePropertyID,
			Value:          "blue",
		}
		err := dao.CreateProperty(ctx, tx, property)
		if err != nil {
			t.Fatalf("Error creating property: %v", err)
		}
		row, err := tx.Query(ctx, `SELECT id, device_id, type_property_id, value FROM device_properties WHERE id = $1`, property.ID)
		if err != nil {
			t.Fatalf("Error querying created property: %v", err)
		}
		defer row.Close()
		if !row.Next() {
			t.Fatal("Expected to find the created property")
		}
		var fetchedProperty utils.DeviceProperty
		if err := row.Scan(&fetchedProperty.ID, &fetchedProperty.DeviceID, &fetchedProperty.TypePropertyID, &fetchedProperty.Value); err != nil {
			t.Fatalf("Error scanning fetched property: %v", err)
		}
		if fetchedProperty.DeviceID != property.DeviceID || fetchedProperty.TypePropertyID != property.TypePropertyID || fetchedProperty.Value != property.Value {
			t.Fatal("Fetched property does not match created property")
		}
	})
	t.Run("CreatePropertyWithExistingID", func(t *testing.T) {
		existingID, _ := gonanoid.New()
		property := &utils.DeviceProperty{
			ID:             existingID,
			DeviceID:       d_id,
			TypePropertyID: typePropertyID,
			Value:          "green",
		}
		err := dao.CreateProperty(ctx, tx, property)
		if err != nil {
			t.Fatalf("Error creating property with existing ID: %v", err)
		}
		row, err := tx.Query(ctx, `SELECT id, device_id, type_property_id, value FROM device_properties WHERE id = $1`, existingID)
		if err != nil {
			t.Fatalf("Error querying created property: %v", err)
		}
		defer row.Close()
		if !row.Next() {
			t.Fatal("Expected to find the created property")
		}
		var fetchedProperty utils.DeviceProperty
		if err := row.Scan(&fetchedProperty.ID, &fetchedProperty.DeviceID, &fetchedProperty.TypePropertyID, &fetchedProperty.Value); err != nil {
			t.Fatalf("Error scanning fetched property: %v", err)
		}
		if fetchedProperty.ID != existingID || fetchedProperty.DeviceID != property.DeviceID || fetchedProperty.TypePropertyID != property.TypePropertyID || fetchedProperty.Value != property.Value {
			t.Fatal("Fetched property does not match created property")
		}
	})
}

func TestUpdateProperty(t *testing.T) {
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
	// Insert a type_property for "color"
	typePropertyID, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO type_properties (id, type_id, name, data_type, required)
		VALUES ($1, $2, $3, $4, $5)
	`, typePropertyID, t_id, "color", "string", true)
	if err != nil {
		t.Fatalf("Error inserting type_property: %v", err)
	}
	p_id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO device_properties (id, device_id, type_property_id, value) VALUES ($1, $2, $3, $4)
		`, p_id, d_id, typePropertyID, "red")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	t.Run("UpdateProperty", func(t *testing.T) {
		property := &utils.DeviceProperty{
			ID:             p_id,
			DeviceID:       d_id,
			TypePropertyID: typePropertyID,
			Value:          "blue",
		}
		err := dao.UpdateProperty(ctx, tx, property)
		if err != nil {
			t.Fatalf("Error updating property: %v", err)
		}
		row, err := tx.Query(ctx, `SELECT id, device_id, type_property_id, value FROM device_properties WHERE id = $1`, property.ID)
		if err != nil {
			t.Fatalf("Error querying updated property: %v", err)
		}
		defer row.Close()
		if !row.Next() {
			t.Fatal("Expected to find the updated property")
		}
		var fetchedProperty utils.DeviceProperty
		if err := row.Scan(&fetchedProperty.ID, &fetchedProperty.DeviceID, &fetchedProperty.TypePropertyID, &fetchedProperty.Value); err != nil {
			t.Fatalf("Error scanning fetched property: %v", err)
		}
		if fetchedProperty.DeviceID != property.DeviceID || fetchedProperty.TypePropertyID != property.TypePropertyID || fetchedProperty.Value != property.Value {
			t.Fatal("Fetched property does not match updated property")
		}
	})
}

func TestDeleteProperty(t *testing.T) {
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
	// Insert a type_property for "color"
	typePropertyID, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO type_properties (id, type_id, name, data_type, required)
		VALUES ($1, $2, $3, $4, $5)
	`, typePropertyID, t_id, "color", "string", true)
	if err != nil {
		t.Fatalf("Error inserting type_property: %v", err)
	}
	p_id, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO device_properties (id, device_id, type_property_id, value) VALUES ($1, $2, $3, $4)
		`, p_id, d_id, typePropertyID, "red")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	t.Run("DeleteProperty", func(t *testing.T) {
		err := dao.DeleteProperty(ctx, tx, p_id)
		if err != nil {
			t.Fatalf("Error deleting property: %v", err)
		}
		row, err := tx.Query(ctx, `SELECT id FROM device_properties WHERE id = $1`, p_id)
		if err != nil {
			t.Fatalf("Error querying deleted property: %v", err)
		}
		defer row.Close()
		if row.Next() {
			t.Fatal("Expected not to find the deleted property")
		}
	})
}
