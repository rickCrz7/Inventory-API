package types

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

func (s *Service) GetType(ctx context.Context, id string) (*utils.Type, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	typ, err := s.dao.GetType(ctx, tx, id)
	if err != nil {
		log.Errorf("Error getting type: %v", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return nil, err
	}
	return typ, nil
}

func (s *Service) GetTypes(ctx context.Context) ([]*utils.Type, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	types, err := s.dao.GetTypes(ctx, tx)
	if err != nil {
		log.Errorf("Error getting types: %v", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return nil, err
	}
	return types, nil
}

func (s *Service) CreateType(ctx context.Context, t *utils.Type) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.CreateType(ctx, tx, t); err != nil {
		log.Errorf("Error creating type: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}
	return nil
}

func (s *Service) UpdateType(ctx context.Context, t *utils.Type) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.UpdateType(ctx, tx, t); err != nil {
		log.Errorf("Error updating type: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}
	return nil
}

func (s *Service) DeleteType(ctx context.Context, id string) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.DeleteType(ctx, tx, id); err != nil {
		log.Errorf("Error deleting type: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}
	return nil
}
