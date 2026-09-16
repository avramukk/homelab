package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
	pool *pgxpool.Pool
}

// connector holds the (possibly not-yet-established) database handle. It keeps
// startup order-independent: the app becomes ready once Postgres is reachable,
// instead of failing permanently if it boots first.
type connector struct {
	dsn string

	mu    sync.RWMutex
	store *store
}

func newConnector(dsn string) *connector { return &connector{dsn: dsn} }

// get returns a live store, attempting a connection if none exists yet.
func (c *connector) get(ctx context.Context) *store {
	c.mu.RLock()
	s := c.store
	c.mu.RUnlock()
	if s != nil {
		return s
	}

	ns, err := newStore(ctx, c.dsn)
	if err != nil {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.store != nil { // lost the race
		ns.Close()
		return c.store
	}
	c.store = ns
	return c.store
}

// run retries until the database is reachable, then returns. The pool handles
// later reconnects on its own.
func (c *connector) run(ctx context.Context, logger *slog.Logger) {
	for {
		if c.get(ctx) != nil {
			logger.Info("database connected", "event", "db_connected")
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
		}
	}
}

func (c *connector) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.store != nil {
		c.store.Close()
	}
}

type item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// newStore connects to Postgres and applies the (idempotent) schema. Returns an
// error instead of exiting so readiness can report the problem.
func newStore(ctx context.Context, dsn string) (*store, error) {
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is empty")
	}
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 8
	cfg.MinConns = 1
	cfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}

	s := &store{pool: pool}
	if err := s.migrate(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return s, nil
}

func (s *store) Close() { s.pool.Close() }

func (s *store) Ping(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return s.pool.Ping(pingCtx)
}

func (s *store) migrate(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS items (
			id         text PRIMARY KEY,
			name       text NOT NULL,
			created_at timestamptz NOT NULL DEFAULT now()
		)`)
	return err
}

func (s *store) ListItems(ctx context.Context) ([]item, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, created_at FROM items ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]item, 0)
	for rows.Next() {
		var it item
		if err := rows.Scan(&it.ID, &it.Name, &it.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (s *store) CreateItem(ctx context.Context, name string) (item, error) {
	it := item{ID: newID(), Name: name}
	err := s.pool.QueryRow(ctx,
		`INSERT INTO items (id, name) VALUES ($1, $2) RETURNING created_at`,
		it.ID, it.Name,
	).Scan(&it.CreatedAt)
	return it, err
}

func newID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b)
}
