package owners

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

func (s *Service) GetOwner(ctx context.Context, id string) (*utils.Owner, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	owner, err := s.dao.GetOwner(ctx, tx, id)
	if err != nil {
		log.Errorf("Failed to get owner: %v", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return nil, err
	}
	return owner, nil
}

func (s *Service) GetOwnerByCampusID(ctx context.Context, campusID string) (*utils.Owner, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	owner, err := s.dao.GetOwnerByCampusID(ctx, tx, campusID)
	if err != nil {
		log.Errorf("Failed to get owner: %v", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return nil, err
	}
	return owner, nil
}

func (s *Service) GetOwnerByEmail(ctx context.Context, email string) (*utils.Owner, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	owner, err := s.dao.GetOwnerByEmail(ctx, tx, email)
	if err != nil {
		log.Errorf("Failed to get owner: %v", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return nil, err
	}
	return owner, nil
}

func (s *Service) GetOwners(ctx context.Context) ([]*utils.Owner, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	owners, err := s.dao.GetOwners(ctx, tx)
	if err != nil {
		log.Errorf("Failed to get owners: %v", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return nil, err
	}
	return owners, nil
}

func (s *Service) CreateOwner(ctx context.Context, owner *utils.Owner) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.CreateOwner(ctx, tx, owner); err != nil {
		log.Errorf("Failed to create owner: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}

func (s *Service) UpdateOwner(ctx context.Context, owner *utils.Owner) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.UpdateOwner(ctx, tx, owner); err != nil {
		log.Errorf("Failed to update owner: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}

func (s *Service) DeleteOwner(ctx context.Context, id string) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.DeleteOwner(ctx, tx, id); err != nil {
		log.Errorf("Failed to delete owner: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}
