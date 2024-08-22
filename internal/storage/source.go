package storage

import (
	"context"
	"time"

	"github.com/disbeliefff/JobHunter/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type SourceStorage struct {
	db *sqlx.DB
}

func NewSourceStorage(db *sqlx.DB) *SourceStorage {
	return &SourceStorage{
		db: db,
	}
}

func (s SourceStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var sources []dbSource
	if err := conn.SelectContext(ctx, &sources, `SELECT * FROM sources`); err != nil {
		return nil, err
	}

	return lo.Map(sources, func(source dbSource, _ int) model.Source {
		return model.Source(source)
	}), nil
}

func (s SourceStorage) SourceByID(ctx context.Context, id int) (*model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var source dbSource
	if err := conn.GetContext(ctx, &source, `SELECT * FROM sources WHERE id = ?`, id); err != nil {
		return nil, err
	}
	return (*model.Source)(&source), nil
}

func (s SourceStorage) Add(ctx context.Context, source model.Source) (int, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var id int

	row := conn.QueryRowxContext(
		ctx,
		`INSERT INTO sources (name, feed_url, priority)
					VALUES ($1, $2, $3) RETURNING id;`,
		source.Name, source.FeedURL, source.CreatedAt,
	)

	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s SourceStorage) Delete(ctx context.Context, id int) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `DELETE FROM sources WHERE id = $1`, id); err != nil {
		return err
	}

	return nil
}

type dbSource struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
