package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radifan9/social-media-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewUserRepository(db *pgxpool.Pool, rdb *redis.Client) *UserRepository {
	return &UserRepository{
		db:  db,
		rdb: rdb,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, email, hashedPassword string) (models.User, error) {
	// Begin transaction
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return models.User{}, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Println("failed to rollback transaction: ", rollbackErr)
			}
		}
	}()

	// Step 1: Create user
	query := `
		INSERT INTO
			users (email, password)
		VALUES
			($1, $2) RETURNING id,
			email`
	var user models.User

	if err := u.db.QueryRow(ctx, query, email, hashedPassword).Scan(&user.Id, &user.Email); err != nil {
		return models.User{}, fmt.Errorf("failed to register user: %w", err)
	}

	// Step 2: Create Profile
	// var profileID string
	_, err = u.createProfile(ctx, tx, user.Id)
	if err != nil {
		return models.User{}, err
	}

	// Commit transaction if everything succeeds
	if err = tx.Commit(ctx); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *UserRepository) createProfile(ctx context.Context, tx pgx.Tx, userID string) (string, error) {
	query := `
		insert into
			user_profiles (user_id)
		values
			($1) returning user_id`

	var id string
	err := tx.QueryRow(ctx, query, userID).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
