package devices

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rickCrz7/Inventory-API/utils"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	dao *Dao
	pdb *pgxpool.Pool
}

func NewService(dao *Dao, pdb *pgxpool.Pool) *Service {
	return &Service{
		dao: dao,
		pdb: pdb,
	}
}

func (s *Service) GetDevice(ctx context.Context, id string) (*utils.Device, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Error starting transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	device, err := s.dao.GetDevice(ctx, tx, id)
	if err != nil {
		log.Errorf("Error fetching device with ID %s: %v", id, err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return nil, err
	}

	return device, nil
}

func (s *Service) GetDevices(ctx context.Context) ([]*utils.Device, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Error starting transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	devices, err := s.dao.GetDevices(ctx, tx)
	if err != nil {
		log.Errorf("Error fetching devices: %v", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return nil, err
	}

	return devices, nil
}

func (s *Service) CreateDevice(ctx context.Context, device *utils.Device) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.CreateDevice(ctx, tx, device); err != nil {
		log.Errorf("Error creating device: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func (s *Service) UpdateDevice(ctx context.Context, device *utils.Device) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.UpdateDevice(ctx, tx, device); err != nil {
		log.Errorf("Error updating device: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func (s *Service) DeleteDevice(ctx context.Context, id string) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.DeleteDevice(ctx, tx, id); err != nil {
		log.Errorf("Error deleting device: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil
}
