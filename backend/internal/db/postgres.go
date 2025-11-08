package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitPostgres() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/synapse?sslmode=disable"
	}

	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func CreateSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		summary TEXT,
		source_url TEXT,
		type TEXT NOT NULL,
		category TEXT,
		tags TEXT[] DEFAULT '{}',
		embedding_id TEXT,
		image_url TEXT,
		embed_html TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS item_relations (
		item_id UUID REFERENCES items(id) ON DELETE CASCADE,
		related_item_id UUID REFERENCES items(id) ON DELETE CASCADE,
		similarity_score FLOAT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		PRIMARY KEY (item_id, related_item_id)
	);

	CREATE INDEX IF NOT EXISTS idx_items_created_at ON items(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_items_tags ON items USING GIN(tags);
	CREATE INDEX IF NOT EXISTS idx_relations_item ON item_relations(item_id);
	`

	_, err := Pool.Exec(context.Background(), schema)
	return err
}

