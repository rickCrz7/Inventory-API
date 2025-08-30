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

func (d *Dao) GetProperties(ctx context.Context, tx pgx.Tx, device_id string) (*utils.DeviceProperty, error) {
	log.Printf("Fetching property with Device ID: %s", device_id)
	var property utils.DeviceProperty
	query := `SELECT id, device_id, type_property_id, value FROM device_properties WHERE device_id = $1`
	err := tx.QueryRow(ctx, query, device_id).Scan(&property.ID, &property.DeviceID, &property.TypePropertyID, &property.Value)
	if err != nil {
		log.Errorf("Error fetching property with Device ID %s: %v", device_id, err)
		return nil, err
	}
	return &property, nil
}

func (d *Dao) CreateProperty(ctx context.Context, tx pgx.Tx, property *utils.DeviceProperty) error {
	log.Printf("Creating property for Device ID: %s", property.DeviceID)
	if property.ID == "" {
		var err error
		property.ID, err = gonanoid.New()
		if err != nil {
			log.Errorf("Error generating ID for property: %v", err)
			return err
		}
	}

	query := `INSERT INTO device_properties (id, device_id, type_property_id, value) VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(ctx, query, property.ID, property.DeviceID, property.TypePropertyID, property.Value)
	if err != nil {
		log.Errorf("Error creating property for Device ID %s: %v", property.DeviceID, err)
		return err
	}
	return nil
}

func (d *Dao) UpdateProperty(ctx context.Context, tx pgx.Tx, property *utils.DeviceProperty) error {
	log.Printf("Updating property for Device ID: %s", property.DeviceID)

	query := `UPDATE device_properties SET type_property_id = $1, value = $2 WHERE id = $3`
	_, err := tx.Exec(ctx, query, property.TypePropertyID, property.Value, property.ID)
	if err != nil {
		log.Errorf("Error updating property for Device ID %s: %v", property.DeviceID, err)
		return err
	}
	return nil
}

func (d *Dao) DeleteProperty(ctx context.Context, tx pgx.Tx, id string) error {
	log.Printf("Deleting property with ID: %s", id)

	query := `DELETE FROM device_properties WHERE id = $1`
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		log.Errorf("Error deleting property with ID %s: %v", id, err)
		return err
	}
	return nil
}