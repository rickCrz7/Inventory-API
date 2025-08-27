package properties

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

func (d *Dao) GetProperties(ctx context.Context, tx pgx.Tx, type_id string) ([]*utils.TypeProperty, error) {
	log.Printf("Fetching properties with Type ID: %s", type_id)
	query := `SELECT id, type_id, name, data_type, required
	FROM type_properties
	WHERE type_id = $1`
	rows, err := tx.Query(ctx, query, type_id)
	if err != nil {
		log.Errorf("Could not get properties for type %s: %v", type_id, err)
		return nil, err
	}
	defer rows.Close()

	var properties []*utils.TypeProperty
	for rows.Next() {
		var property utils.TypeProperty
		if err := rows.Scan(&property.ID, &property.TypeID, &property.Name, &property.DataType, &property.Required); err != nil {
			log.Errorf("Could not scan property: %v", err)
			return nil, err
		}
		properties = append(properties, &property)
	}
	if err := rows.Err(); err != nil {
		log.Errorf("Error occurred while fetching properties: %v", err)
		return nil, err
	}
	return properties, nil
}

func (d *Dao) CreateProperty(ctx context.Context, tx pgx.Tx, property *utils.TypeProperty) error {
	log.Printf("Creating property: %v", property)
	if property.ID == "" {
		var err error
		property.ID, err = gonanoid.New()
		if err != nil {
			log.Errorf("Could not generate property ID: %v", err)
			return err
		}
	}
	query := `INSERT INTO type_properties (id, type_id, name, data_type, required)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, query, property.ID, property.TypeID, property.Name, property.DataType, property.Required)
	if err != nil {
		log.Errorf("Could not create property: %v", err)
		return err
	}
	return nil
}

func (d *Dao) UpdateProperty(ctx context.Context, tx pgx.Tx, property *utils.TypeProperty) error {
	log.Printf("Updating property: %v", property)
	query := `UPDATE type_properties
	SET type_id = $2, name = $3, data_type = $4, required = $5
	WHERE id = $1`
	_, err := tx.Exec(ctx, query, property.ID, property.TypeID, property.Name, property.DataType, property.Required)
	if err != nil {
		log.Errorf("Could not update property %s: %v", property.ID, err)
		return err
	}
	return nil
}

func (d *Dao) DeleteProperty(ctx context.Context, tx pgx.Tx, id string) error {
	log.Printf("Deleting property with ID: %s", id)
	query := `DELETE FROM type_properties WHERE id = $1`
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		log.Errorf("Could not delete property %s: %v", id, err)
		return err
	}
	return nil
}
