package properties

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

func (s *Service) GetProperties(ctx context.Context, id string) (*utils.DeviceProperty, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	props, err := s.dao.GetProperties(ctx, tx, id)
	if err != nil {
		log.Errorf("Error getting properties: %v", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return nil, err
	}
	return props, nil
}

func (s *Service) CreateProperty(ctx context.Context, prop *utils.DeviceProperty) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.CreateProperty(ctx, tx, prop); err != nil {
		log.Errorf("Error creating property: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}
	return nil
}

func (s *Service) UpdateProperty(ctx context.Context, prop *utils.DeviceProperty) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.UpdateProperty(ctx, tx, prop); err != nil {
		log.Errorf("Error updating property: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}
	return nil
}

func (s *Service) DeleteProperty(ctx context.Context, id string) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.DeleteProperty(ctx, tx, id); err != nil {
		log.Errorf("Error deleting property: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}
	return nil
}
