package types

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

func (d *Dao) GetType(ctx context.Context, tx pgx.Tx, id string) (*utils.Type, error) {
	log.Printf("Fetching type with ID: %s", id)
	query := `SELECT id, name, description FROM types WHERE id = $1`
	var t utils.Type
	err := tx.QueryRow(ctx, query, id).Scan(&t.ID, &t.Name, &t.Description)
	if err != nil {
		log.Errorf("Error fetching type with ID %s: %v", id, err)
		return nil, err
	}
	return &t, nil
}

func (d *Dao) GetTypes(ctx context.Context, tx pgx.Tx) ([]*utils.Type, error) {
	log.Println("Fetching all types")
	query := `SELECT id, name, description FROM types`
	rows, err := tx.Query(ctx, query)
	if err != nil {
		log.Errorf("Error fetching types: %v", err)
		return nil, err
	}
	defer rows.Close()

	var types []*utils.Type
	for rows.Next() {
		var t utils.Type
		if err := rows.Scan(&t.ID, &t.Name, &t.Description); err != nil {
			log.Errorf("Error scanning type: %v", err)
			return nil, err
		}
		types = append(types, &t)
	}
	if err := rows.Err(); err != nil {
		log.Errorf("Error with rows: %v", err)
		return nil, err
	}
	return types, nil
}

func (d *Dao) CreateType(ctx context.Context, tx pgx.Tx, t *utils.Type) error {
	log.Printf("Creating type: %+v", t)
	if t.ID == "" {
		var err error
		t.ID, err = gonanoid.New()
		if err != nil {
			log.Errorf("Error generating ID: %v", err)
			return err
		}
	}
	query := `INSERT INTO types (id, name, description) VALUES ($1, $2, $3)`
	_, err := tx.Exec(ctx, query, t.ID, t.Name, t.Description)
	if err != nil {
		log.Errorf("Error creating type: %v", err)
		return err
	}
	return nil
}

func (d *Dao) UpdateType(ctx context.Context, tx pgx.Tx, t *utils.Type) error {
	log.Printf("Updating type: %+v", t)
	query := `UPDATE types SET name = $1, description = $2 WHERE id = $3`
	_, err := tx.Exec(ctx, query, t.Name, t.Description, t.ID)
	if err != nil {
		log.Errorf("Error updating type: %v", err)
		return err
	}
	return nil
}

func (d *Dao) DeleteType(ctx context.Context, tx pgx.Tx, id string) error {
	log.Printf("Deleting type with ID: %s", id)
	query := `DELETE FROM types WHERE id = $1`
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		log.Errorf("Error deleting type with ID %s: %v", id, err)
		return err
	}
	return nil
}
