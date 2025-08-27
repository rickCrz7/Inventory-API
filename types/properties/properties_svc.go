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

func (s *Service) GetProperties(ctx context.Context, type_id string) ([]*utils.TypeProperty, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	properties, err := s.dao.GetProperties(ctx, tx, type_id)
	if err != nil {
		log.Errorf("Failed to get properties: %v", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return nil, err
	}
	return properties, nil
}

func (s *Service) CreateProperty(ctx context.Context, property *utils.TypeProperty) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.CreateProperty(ctx, tx, property); err != nil {
		log.Errorf("Failed to create property: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}

func (s *Service) UpdateProperty(ctx context.Context, property *utils.TypeProperty) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.UpdateProperty(ctx, tx, property); err != nil {
		log.Errorf("Failed to update property: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}

func (s *Service) DeleteProperty(ctx context.Context, id string) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.DeleteProperty(ctx, tx, id); err != nil {
		log.Errorf("Failed to delete property: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}