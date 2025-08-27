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
		INSERT INTO types (id, name, description) VALUES ($1, $2, $3)
	`, id, "Test Type", "This is a test type")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	for i := 0; i < 3; i++ {
		pid, _ := gonanoid.New()
		_, err = tx.Exec(ctx, `
			INSERT INTO type_properties (id, type_id, name, data_type, required) VALUES ($1, $2, $3, $4, $5)
		`, pid, id, fmt.Sprintf("Test Property %d", i), fmt.Sprintf("This is a test property %d", i), true)
		if err != nil {
			t.Fatalf("Error inserting mock data: %v", err)
		}
	}
	t.Run("GetProperties", func(t *testing.T) {
		properties, err := dao.GetProperties(ctx, tx, id)
		if err != nil {
			t.Fatalf("Error getting properties: %v", err)
		}
		if len(properties) < 3 {
			t.Errorf("Expected at least 3 properties, but got %d", len(properties))
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
		INSERT INTO types (id, name, description) VALUES ($1, $2, $3)
	`, id, "Test Type", "This is a test type")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}

	t.Run("CreateProperty", func(t *testing.T) {
		pid, _ := gonanoid.New()
		property := &utils.TypeProperty{
			ID:       pid,
			TypeID:   id,
			Name:     "New Property",
			DataType: "string",
			Required: true,
		}
		err := dao.CreateProperty(ctx, tx, property)
		if err != nil {
			t.Fatalf("Error creating property: %v", err)
		}
		propertyFetched, err := dao.GetProperties(ctx, tx, id)
		if err != nil {
			t.Fatalf("Error fetching properties: %v", err)
		}
		if len(propertyFetched) == 0 {
			t.Errorf("Expected to fetch the created property, but got none")
		}
		if propertyFetched[0].ID != property.ID {
			t.Errorf("Expected to fetch property with ID %s, but got %s", property.ID, propertyFetched[0].ID)
		}
		if propertyFetched[0].Name != property.Name {
			t.Errorf("Expected to fetch property with Name %s, but got %s", property.Name, propertyFetched[0].Name)
		}
		if propertyFetched[0].DataType != property.DataType {
			t.Errorf("Expected to fetch property with DataType %s, but got %s", property.DataType, propertyFetched[0].DataType)
		}
		if propertyFetched[0].Required != property.Required {
			t.Errorf("Expected to fetch property with Required %v, but got %v", property.Required, propertyFetched[0].Required)
		}
	})
	t.Run("CreatePropertyWithoutID", func(t *testing.T) {
		property := &utils.TypeProperty{
			TypeID:   id,
			Name:     "New Property",
			DataType: "string",
			Required: true,
		}
		err := dao.CreateProperty(ctx, tx, property)
		if err != nil {
			t.Fatalf("Error creating property: %v", err)
		}
		propertyFetched, err := dao.GetProperties(ctx, tx, id)
		if err != nil {
			t.Fatalf("Error fetching properties: %v", err)
		}
		if len(propertyFetched) == 0 {
			t.Errorf("Expected to fetch the created property, but got none")
		}
		if propertyFetched[0].ID == "" {
			t.Errorf("Expected to fetch property with generated ID, but got empty ID")
		}
		if propertyFetched[0].Name != property.Name {
			t.Errorf("Expected to fetch property with Name %s, but got %s", property.Name, propertyFetched[0].Name)
		}
		if propertyFetched[0].DataType != property.DataType {
			t.Errorf("Expected to fetch property with DataType %s, but got %s", property.DataType, propertyFetched[0].DataType)
		}
		if propertyFetched[0].Required != property.Required {
			t.Errorf("Expected to fetch property with Required %v, but got %v", property.Required, propertyFetched[0].Required)
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
		INSERT INTO types (id, name, description) VALUES ($1, $2, $3)
	`, id, "Test Type", "This is a test type")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}

	pid, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO type_properties (id, type_id, name, data_type, required) VALUES ($1, $2, $3, $4, $5)
	`, pid, id, "Test Property", "This is a test property", true)
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}

	t.Run("UpdateProperty", func(t *testing.T) {
		property := &utils.TypeProperty{
			ID:       pid,
			TypeID:   id,
			Name:     "Updated Property",
			DataType: "string",
			Required: true,
		}
		err := dao.UpdateProperty(ctx, tx, property)
		if err != nil {
			t.Fatalf("Error updating property: %v", err)
		}
		propertyFetched, err := dao.GetProperties(ctx, tx, id)
		if err != nil {
			t.Fatalf("Error fetching properties: %v", err)
		}
		if len(propertyFetched) == 0 {
			t.Errorf("Expected to fetch the updated property, but got none")
		}
		if propertyFetched[0].ID != property.ID {
			t.Errorf("Expected to fetch property with ID %s, but got %s", property.ID, propertyFetched[0].ID)
		}
		if propertyFetched[0].Name != property.Name {
			t.Errorf("Expected to fetch property with Name %s, but got %s", property.Name, propertyFetched[0].Name)
		}
		if propertyFetched[0].DataType != property.DataType {
			t.Errorf("Expected to fetch property with DataType %s, but got %s", property.DataType, propertyFetched[0].DataType)
		}
		if propertyFetched[0].Required != property.Required {
			t.Errorf("Expected to fetch property with Required %v, but got %v", property.Required, propertyFetched[0].Required)
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
		INSERT INTO types (id, name, description) VALUES ($1, $2, $3)
	`, id, "Test Type", "This is a test type")
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}

	pid, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO type_properties (id, type_id, name, data_type, required) VALUES ($1, $2, $3, $4, $5)
	`, pid, id, "Test Property", "This is a test property", true)
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}

	t.Run("DeleteProperty", func(t *testing.T) {
		err := dao.DeleteProperty(ctx, tx, pid)
		if err != nil {
			t.Fatalf("Error deleting property: %v", err)
		}
		propertyFetched, err := dao.GetProperties(ctx, tx, id)
		if err != nil {
			t.Fatalf("Error fetching properties: %v", err)
		}
		if len(propertyFetched) != 0 {
			t.Errorf("Expected no properties to be fetched, but got %d", len(propertyFetched))
		}
	})
}