package logs

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

func (s *Service) GetLogs(ctx context.Context, deviceID string) ([]*utils.DeviceLog, error) {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		log.Errorf("Error starting transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	logs, err := s.dao.GetLogs(ctx, tx, deviceID)
	if err != nil {
		log.Errorf("Error fetching logs: %v", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return nil, err
	}

	return logs, nil
}

func (s *Service) CreateLog(ctx context.Context, logEntry *utils.DeviceLog) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.CreateLog(ctx, tx, logEntry); err != nil {
		log.Errorf("Error creating log: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func (s *Service) DeleteLog(ctx context.Context, id string) error {
	tx, err := s.pdb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.dao.DeleteLog(ctx, tx, id); err != nil {
		log.Errorf("Error deleting log: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil
}
