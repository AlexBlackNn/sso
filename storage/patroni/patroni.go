package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/XSAM/otelsql"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"sso/internal/config"
	"sso/internal/domain/models"
	"sso/storage"
)

const ErrCodeUserAlreadyExists = "23505"

type Storage struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
}

var tracer = otel.Tracer("sso service")

func New(cfg *config.Config) (*Storage, error) {
	dbWrite, err := otelsql.Open("pgx", cfg.StoragePatroni.Master)
	if err != nil {
		return nil, fmt.Errorf(
			"DATA LAYER: storage.postgres.New: couldn't open a database for Write: %w",
			err,
		)
	}
	dbRead, err := otelsql.Open("pgx", cfg.StoragePatroni.Slave)
	if err != nil {
		return nil, fmt.Errorf(
			"DATA LAYER: storage.postgres.New: couldn't open a database for Read: %w",
			err,
		)
	}
	return &Storage{dbRead: dbRead, dbWrite: dbWrite}, nil
}

func (s *Storage) Stop() error {
	err1 := s.dbWrite.Close()
	err2 := s.dbRead.Close()
	return fmt.Errorf("%w, %w", err1, err2)
}

// SaveUser saves user to db.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (context.Context, int64, error) {
	ctx, span := tracer.Start(ctx, "data layer Patroni: SaveUser",
		trace.WithAttributes(attribute.String("handler", "SaveUser")))
	defer span.End()

	var id int
	query := "INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id"
	err := s.dbWrite.QueryRowContext(ctx, query, email, passHash).Scan(&id)
	// https://stackoverflow.com/questions/34963064/go-pq-and-postgres-appropriate-error-handling-for-constraints
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == ErrCodeUserAlreadyExists {
			return ctx, int64(0), storage.ErrUserExists
		}
	}
	if err != nil {
		return ctx, 0, fmt.Errorf(
			"DATA LAYER: storage.postgres.SaveUser: couldn't save user  %w",
			err,
		)
	}

	return ctx, int64(id), nil
}

func (s *Storage) GetUser(ctx context.Context, value any) (context.Context, models.User, error) {
	ctx, span := tracer.Start(ctx, "data layer Patroni: GetUser",
		trace.WithAttributes(attribute.String("handler", "GetUser")))
	defer span.End()

	var row *sql.Row
	switch sqlParam := value.(type) {
	case int:
		query := "SELECT id, email, pass_hash, is_admin FROM users WHERE (id = $1);"
		row = s.dbRead.QueryRowContext(ctx, query, sqlParam)
	case string:
		query := "SELECT id, email, pass_hash, is_admin FROM users WHERE (email = $1);"
		row = s.dbRead.QueryRowContext(ctx, query, sqlParam)
	default:
		return ctx, models.User{}, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetUser: %w",
			storage.ErrWrongParamType,
		)
	}

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx, models.User{}, fmt.Errorf(
				"DATA LAYER: storage.postgres.GetUser: %w",
				storage.ErrUserNotFound,
			)
		}
		return ctx, models.User{}, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetUser: %w",
			err,
		)
	}
	return ctx, user, nil
}
