package devices

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

func (d *Dao) GetDevice(ctx context.Context, tx pgx.Tx, id string) (*utils.Device, error) {
	log.Printf("Fetching device with ID: %s", id)
	query := `SELECT id, serial_number, name, type_id, owner_id, purchase_date, status
	FROM devices
	WHERE id = $1`
	var device utils.Device
	err := tx.QueryRow(ctx, query, id).Scan(&device.ID, &device.SerialNumber, &device.Name, &device.TypeID, &device.OwnerID, &device.PurchaseDate, &device.Status)
	if err != nil {
		log.Errorf("Error fetching device with ID %s: %v", id, err)
		return nil, err
	}
	return &device, nil
}

func (d *Dao) GetDevices(ctx context.Context, tx pgx.Tx) ([]*utils.Device, error) {
	log.Println("Fetching all devices")
	query := `SELECT id, serial_number, name, type_id, owner_id, purchase_date, status
	FROM devices`
	rows, err := tx.Query(ctx, query)
	if err != nil {
		log.Errorf("Error fetching devices: %v", err)
		return nil, err
	}
	defer rows.Close()

	var devices []*utils.Device
	for rows.Next() {
		var device utils.Device
		if err := rows.Scan(&device.ID, &device.SerialNumber, &device.Name, &device.TypeID, &device.OwnerID, &device.PurchaseDate, &device.Status); err != nil {
			log.Errorf("Error scanning device row: %v", err)
			return nil, err
		}
		devices = append(devices, &device)
	}
	if err := rows.Err(); err != nil {
		log.Errorf("Error iterating over device rows: %v", err)
		return nil, err
	}
	return devices, nil
}

func (d *Dao) CreateDevice(ctx context.Context, tx pgx.Tx, device *utils.Device) error {
	log.Printf("Creating device: %+v", device)
	if device.ID == "" {
		id, err := gonanoid.New()
		if err != nil {
			log.Errorf("Error generating ID for new device: %v", err)
			return err
		}
		device.ID = id
	}
	query := `INSERT INTO devices (id, serial_number, name, type_id, owner_id, purchase_date, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := tx.Exec(ctx, query, device.ID, device.SerialNumber, device.Name, device.TypeID, device.OwnerID, device.PurchaseDate, device.Status)
	if err != nil {
		log.Errorf("Error creating device: %v", err)
		return err
	}
	return nil
}

func (d *Dao) UpdateDevice(ctx context.Context, tx pgx.Tx, device *utils.Device) error {
	log.Printf("Updating device: %+v", device)
	query := `UPDATE devices
	SET serial_number = $1, name = $2, type_id = $3, owner_id = $4, purchase_date = $5, status = $6
	WHERE id = $7`
	_, err := tx.Exec(ctx, query, device.SerialNumber, device.Name, device.TypeID, device.OwnerID, device.PurchaseDate, device.Status, device.ID)
	if err != nil {
		log.Errorf("Error updating device with ID %s: %v", device.ID, err)
		return err
	}
	return nil
}

func (d *Dao) DeleteDevice(ctx context.Context, tx pgx.Tx, id string) error {
	log.Printf("Deleting device with ID: %s", id)
	query := `DELETE FROM devices WHERE id = $1`
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		log.Errorf("Error deleting device with ID %s: %v", id, err)
		return err
	}
	return nil
}
