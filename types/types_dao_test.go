package types

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

func TestGetType(t *testing.T) {
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
	t.Run("GetType", func(t *testing.T) {
		typ, err := dao.GetType(ctx, tx, id)
		if err != nil {
			t.Fatalf("Error getting type: %v", err)
		}
		if typ == nil {
			t.Fatal("Expected type to be found, but got nil")
		}
		if typ.ID != id {
			t.Errorf("Expected type ID %q, but got %q", id, typ.ID)
		}
	})
	t.Run("GetTypeNotFound", func(t *testing.T) {
		typ, err := dao.GetType(ctx, tx, "non-existent-id")
		if err == nil {
			t.Fatal("Expected error getting type, but got none")
		}
		if typ != nil {
			t.Fatal("Expected type to be nil, but got a value")
		}
	})
}

func TestGetTypes(t *testing.T) {
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
	for i := 0; i < 3; i++ {
		id, _ := gonanoid.New()
		_, err = tx.Exec(ctx, `
			INSERT INTO types (id, name, description) VALUES ($1, $2, $3)
		`, id, fmt.Sprintf("Test Type %d", i), fmt.Sprintf("This is a test type %d", i))
		if err != nil {
			t.Fatalf("Error inserting mock data: %v", err)
		}
	}
	t.Run("GetTypes", func(t *testing.T) {
		types, err := dao.GetTypes(ctx, tx)
		if err != nil {
			t.Fatalf("Error getting types: %v", err)
		}
		if len(types) < 3 {
			t.Errorf("Expected at least 3 types, but got %d", len(types))
		}
	})
}

func TestCreateType(t *testing.T) {
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

	t.Run("CreateType", func(t *testing.T) {
		description := "This is a new type"
		typ := &utils.Type{
			Name:        "New Type",
			Description: &description,
		}
		err := dao.CreateType(ctx, tx, typ)
		if err != nil {
			t.Fatalf("Error creating type: %v", err)
		}
		if typ.ID == "" {
			t.Error("Expected type ID to be set, but got empty string")
		}
	})

	t.Run("CreateTypeWithExistingID", func(t *testing.T) {
		description := "This is a new type"
		typ := &utils.Type{
			ID:          "existing-id",
			Name:        "New Type",
			Description: &description,
		}
		err := dao.CreateType(ctx, tx, typ)
		if err != nil {
			t.Fatalf("Error creating type with existing ID: %v", err)
		}
		createdType, err := dao.GetType(ctx, tx, typ.ID)
		if err != nil {
			t.Fatalf("Error getting created type: %v", err)
		}
		if createdType == nil {
			t.Fatal("Expected created type to be found, but got nil")
		}
	})
}

func TestUpdateType(t *testing.T) {
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
	t.Run("UpdateType", func(t *testing.T) {
		description := "This is an updated test type"
		typ := &utils.Type{
			ID:          id,
			Name:        "Updated Type",
			Description: &description,
		}
		err := dao.UpdateType(ctx, tx, typ)
		if err != nil {
			t.Fatalf("Error updating type: %v", err)
		}
		updatedType, err := dao.GetType(ctx, tx, id)
		if err != nil {
			t.Fatalf("Error getting updated type: %v", err)
		}
		if updatedType == nil {
			t.Fatal("Expected updated type to be found, but got nil")
		}
		if updatedType.Name != typ.Name {
			t.Errorf("Expected type name %q, but got %q", typ.Name, updatedType.Name)
		}
		if *updatedType.Description != *typ.Description {
			t.Errorf("Expected type description %q, but got %q", *typ.Description, *updatedType.Description)
		}
	})
}

func TestDeleteType(t *testing.T) {
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
	t.Run("DeleteType", func(t *testing.T) {
		err := dao.DeleteType(ctx, tx, id)
		if err != nil {
			t.Fatalf("Error deleting type: %v", err)
		}
		deletedType, err := dao.GetType(ctx, tx, id)
		if err == nil {
			t.Fatal("Expected error getting deleted type, but got none")
		}
		if deletedType != nil {
			t.Fatal("Expected deleted type to be nil, but got a value")
		}
	})
}
