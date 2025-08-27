package owners

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

func TestGetOwner(t *testing.T) {
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
	t.Run("GetOwner", func(t *testing.T) {
		owner, err := dao.GetOwner(ctx, tx, id)
		if err != nil {
			t.Fatalf("Error getting owner: %v", err)
		}
		if owner == nil {
			t.Fatal("Expected owner to be found, but got nil")
		}
		if owner.ID != id {
			t.Errorf("Expected owner ID %s, but got %s", id, owner.ID)
		}
	})
	t.Run("GetOwnerNotFound", func(t *testing.T) {
		owner, err := dao.GetOwner(ctx, tx, "nonexistent-id")
		if err == nil {
			t.Fatalf("Expected error getting owner, but got none")
		}
		if owner != nil {
			t.Errorf("Expected owner to be nil, but got %v", owner)
		}
	})
}

func TestGetOwnerByCampusID(t *testing.T) {
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
	campusID, _ := gonanoid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO owners (id, first_name, last_name, email, campus_id)
		VALUES ($1, $2, $3, $4, $5)
	`, id, "John", "Doe", "john.doe@example.com", campusID)
	if err != nil {
		t.Fatalf("Error inserting mock data: %v", err)
	}
	t.Run("GetOwnerByCampusID", func(t *testing.T) {
		owner, err := dao.GetOwnerByCampusID(ctx, tx, campusID)
		if err != nil {
			t.Fatalf("Error getting owner: %v", err)
		}
		if owner == nil {
			t.Fatal("Expected owner to be found, but got nil")
		}
		if owner.CampusID == nil || *owner.CampusID != campusID {
			t.Errorf("Expected owner Campus ID %s, but got %v", campusID, owner.CampusID)
		}
	})
	t.Run("GetOwnerByCampusIDNotFound", func(t *testing.T) {
		owner, err := dao.GetOwnerByCampusID(ctx, tx, "nonexistent-campus-id")
		if err == nil {
			t.Fatalf("Expected error getting owner, but got none")
		}
		if owner != nil {
			t.Errorf("Expected owner to be nil, but got %v", owner)
		}
	})
}

func TestGetOwnerByEmail(t *testing.T) {
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
	t.Run("GetOwnerByEmail", func(t *testing.T) {
		owner, err := dao.GetOwnerByEmail(ctx, tx, "john.doe@example.com")
		if err != nil {
			t.Fatalf("Error getting owner: %v", err)
		}
		if owner == nil {
			t.Fatal("Expected owner to be found, but got nil")
		}
		if owner.ID != id {
			t.Errorf("Expected owner ID %s, but got %s", id, owner.ID)
		}
	})
	t.Run("GetOwnerByEmailNotFound", func(t *testing.T) {
		owner, err := dao.GetOwnerByEmail(ctx, tx, "nonexistent-email@example.com")
		if err == nil {
			t.Fatalf("Expected error getting owner, but got none")
		}
		if owner != nil {
			t.Errorf("Expected owner to be nil, but got %v", owner)
		}
	})
}

func TestGetOwners(t *testing.T) {
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
			INSERT INTO owners (id, first_name, last_name, email)
			VALUES ($1, $2, $3, $4)
		`, id, fmt.Sprintf("John%d", i), fmt.Sprintf("Doe%d", i), fmt.Sprintf("john.doe%d@example.com", i))
		if err != nil {
			t.Fatalf("Error inserting mock data: %v", err)
		}
	}
	t.Run("GetOwners", func(t *testing.T) {
		owners, err := dao.GetOwners(ctx, tx)
		if err != nil {
			t.Fatalf("Error getting owners: %v", err)
		}
		if len(owners) < 3 {
			t.Errorf("Expected at least 3 owners, but got %d", len(owners))
		}
	})
}

func TestCreateOwner(t *testing.T) {
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

	t.Run("CreateOwner", func(t *testing.T) {
		owner := &utils.Owner{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
		}
		err := dao.CreateOwner(ctx, tx, owner)
		if err != nil {
			t.Fatalf("Error creating owner: %v", err)
		}
		createdOwner, err := dao.GetOwnerByEmail(ctx, tx, owner.Email)
		if err != nil {
			t.Fatalf("Error getting created owner: %v", err)
		}
		if createdOwner == nil {
			t.Fatal("Expected created owner to be found, but got nil")
		}
		if owner.ID == "" {
			t.Error("Expected owner ID to be set, but got empty string")
		}
	})

	t.Run("CreateOwnerWithExistingID", func(t *testing.T) {
		owner := &utils.Owner{
			ID:        "existing-id",
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
		}
		err := dao.CreateOwner(ctx, tx, owner)
		if err != nil {
			t.Fatalf("Error creating owner with existing ID: %v", err)
		}
		createdOwner, err := dao.GetOwner(ctx, tx, owner.ID)
		if err != nil {
			t.Fatalf("Error getting created owner: %v", err)
		}
		if createdOwner == nil {
			t.Fatal("Expected created owner to be found, but got nil")
		}
	})
}

func TestUpdateOwner(t *testing.T) {
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

	t.Run("UpdateOwner", func(t *testing.T) {
		owner := &utils.Owner{
			ID:        id,
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
		}
		err := dao.UpdateOwner(ctx, tx, owner)
		if err != nil {
			t.Fatalf("Error updating owner: %v", err)
		}
		updatedOwner, err := dao.GetOwner(ctx, tx, owner.ID)
		if err != nil {
			t.Fatalf("Error getting updated owner: %v", err)
		}
		if updatedOwner == nil {
			t.Fatal("Expected updated owner to be found, but got nil")
		}
		if updatedOwner.FirstName != owner.FirstName {
			t.Errorf("Expected first name %s, but got %s", owner.FirstName, updatedOwner.FirstName)
		}
	})
}

func TestDeleteOwner(t *testing.T) {
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

	t.Run("DeleteOwner", func(t *testing.T) {
		err := dao.DeleteOwner(ctx, tx, id)
		if err != nil {
			t.Fatalf("Error deleting owner: %v", err)
		}
		deletedOwner, err := dao.GetOwner(ctx, tx, id)
		if err == nil {
			t.Fatalf("Expected error getting deleted owner, but got none")
		}
		if deletedOwner != nil {
			t.Errorf("Expected deleted owner to be nil, but got %+v", deletedOwner)
		}
	})
}
