package owners

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

func (d *Dao) GetOwner(ctx context.Context, tx pgx.Tx, id string) (*utils.Owner, error) {
	log.Printf("Fetching owner with ID: %s", id)
	query := `SELECT id, first_name, last_name, campus_id, email
	FROM owners 
	WHERE id = $1`
	row := tx.QueryRow(ctx, query, id)
	var owner utils.Owner
	err := row.Scan(&owner.ID, &owner.FirstName, &owner.LastName,
		&owner.CampusID, &owner.Email)
	if err != nil {
		log.Errorf("Could not get owner %s: %v", id, err)
		return nil, err
	}
	return &owner, nil
}

func (d *Dao) GetOwnerByCampusID(ctx context.Context, tx pgx.Tx, campus_id string) (*utils.Owner, error) {
	log.Printf("Fetching owner with Campus ID: %s", campus_id)
	query := `SELECT id, first_name, last_name, campus_id, email
	FROM owners
	WHERE campus_id = $1`
	row := tx.QueryRow(ctx, query, campus_id)
	var owner utils.Owner
	err := row.Scan(&owner.ID, &owner.FirstName, &owner.LastName,
		&owner.CampusID, &owner.Email)
	if err != nil {
		log.Errorf("Could not get owner %s: %v", campus_id, err)
		return nil, err
	}
	return &owner, nil
}

func (d *Dao) GetOwnerByEmail(ctx context.Context, tx pgx.Tx, email string) (*utils.Owner, error) {
	log.Printf("Fetching owner with Email: %s", email)
	query := `SELECT id, first_name, last_name, campus_id, email
	FROM owners
	WHERE email = $1`
	row := tx.QueryRow(ctx, query, email)
	var owner utils.Owner
	err := row.Scan(&owner.ID, &owner.FirstName, &owner.LastName,
		&owner.CampusID, &owner.Email)
	if err != nil {
		log.Errorf("Could not get owner %s: %v", email, err)
		return nil, err
	}
	return &owner, nil
}

func (d *Dao) GetOwners(ctx context.Context, tx pgx.Tx) ([]*utils.Owner, error) {
	log.Printf("Fetching all owners")
	query := `SELECT id, first_name, last_name, campus_id, email
	FROM owners`
	rows, err := tx.Query(ctx, query)
	if err != nil {
		log.Errorf("Could not get owners: %v", err)
		return nil, err
	}
	defer rows.Close()

	var owners []*utils.Owner
	for rows.Next() {
		var owner utils.Owner
		if err := rows.Scan(&owner.ID, &owner.FirstName, &owner.LastName,
			&owner.CampusID, &owner.Email); err != nil {
			log.Errorf("Could not scan owner: %v", err)
			return nil, err
		}
		owners = append(owners, &owner)
	}
	if err := rows.Err(); err != nil {
		log.Errorf("Error occurred while fetching owners: %v", err)
		return nil, err
	}
	return owners, nil
}

func (d *Dao) CreateOwner(ctx context.Context, tx pgx.Tx, owner *utils.Owner) error {
	log.Printf("Creating owner: %v", owner)
	if owner.ID == "" {
		var err error
		owner.ID, err = gonanoid.New()
		if err != nil {
			log.Errorf("Could not generate owner ID: %v", err)
			return err
		}
	}
	query := `INSERT INTO owners (id, first_name, last_name, email, campus_id)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, query, owner.ID, owner.FirstName, owner.LastName, owner.Email, owner.CampusID)
	if err != nil {
		log.Errorf("Could not create owner: %v", err)
		return err
	}
	return nil
}

func (d *Dao) UpdateOwner(ctx context.Context, tx pgx.Tx, owner *utils.Owner) error {
	log.Printf("Updating owner: %v", owner)
	query := `UPDATE owners
	SET first_name = $1, last_name = $2, email = $3, campus_id = $4
	WHERE id = $5`
	_, err := tx.Exec(ctx, query, owner.FirstName, owner.LastName, owner.Email, owner.CampusID, owner.ID)
	if err != nil {
		log.Errorf("Could not update owner: %v", err)
		return err
	}
	return nil
}

func (d *Dao) DeleteOwner(ctx context.Context, tx pgx.Tx, id string) error {
	log.Printf("Deleting owner with ID: %s", id)
	query := `DELETE FROM owners WHERE id = $1`
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		log.Errorf("Could not delete owner %s: %v", id, err)
		return err
	}
	return nil
}
