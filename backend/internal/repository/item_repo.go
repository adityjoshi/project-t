package repository

import (
	"context"
	"database/sql"
	"fmt"
	"synapse/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"
)

type ItemRepository struct {
	pool *pgxpool.Pool
}

func NewItemRepository(pool *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{pool: pool}
}

func (r *ItemRepository) Create(ctx context.Context, item *models.Item) error {
	query := `
		INSERT INTO items (id, title, content, summary, source_url, type, category, tags, embedding_id, image_url, embed_html, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	
	tagsArray := pgtype.Array[string]{
		Elements: item.Tags,
		Valid:    true,
	}
	
	_, err := r.pool.Exec(ctx, query,
		item.ID, item.Title, item.Content, item.Summary, item.SourceURL,
		item.Type, item.Category, tagsArray, item.EmbeddingID, item.ImageURL, item.EmbedHTML, item.CreatedAt,
	)
	return err
}

func (r *ItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	query := `
		SELECT id, title, content, summary, source_url, type, category, tags, embedding_id, image_url, embed_html, created_at
		FROM items
		WHERE id = $1
	`
	
	var item models.Item
	var tagsArray pgtype.Array[string]
	var imageURL, embedHTML, category sql.NullString
	
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.Title, &item.Content, &item.Summary, &item.SourceURL,
		&item.Type, &category, &tagsArray, &item.EmbeddingID, &imageURL, &embedHTML, &item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	item.Tags = tagsArray.Elements
	if category.Valid {
		item.Category = category.String
	}
	if imageURL.Valid {
		item.ImageURL = imageURL.String
	}
	if embedHTML.Valid {
		item.EmbedHTML = embedHTML.String
	}
	return &item, nil
}

func (r *ItemRepository) GetAll(ctx context.Context) ([]models.Item, error) {
	query := `
		SELECT id, title, content, summary, source_url, type, category, tags, embedding_id, image_url, embed_html, created_at
		FROM items
		ORDER BY created_at DESC
	`
	
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return []models.Item{}, err
	}
	defer rows.Close()
	
	items := []models.Item{}
	for rows.Next() {
		var item models.Item
		var tagsArray pgtype.Array[string]
		var imageURL, embedHTML, category sql.NullString
		
		err := rows.Scan(
			&item.ID, &item.Title, &item.Content, &item.Summary, &item.SourceURL,
			&item.Type, &category, &tagsArray, &item.EmbeddingID, &imageURL, &embedHTML, &item.CreatedAt,
		)
		if err != nil {
			return []models.Item{}, err
		}
		
		item.Tags = tagsArray.Elements
		if category.Valid {
			item.Category = category.String
		}
		if imageURL.Valid {
			item.ImageURL = imageURL.String
		}
		if embedHTML.Valid {
			item.EmbedHTML = embedHTML.String
		}
		items = append(items, item)
	}
	
	return items, nil
}

func (r *ItemRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Item, error) {
	if len(ids) == 0 {
		return []models.Item{}, nil
	}
	
	query := `
		SELECT id, title, content, summary, source_url, type, category, tags, embedding_id, image_url, embed_html, created_at
		FROM items
		WHERE id = ANY($1)
	`
	
	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var items []models.Item
	for rows.Next() {
		var item models.Item
		var tagsArray pgtype.Array[string]
		var imageURL, embedHTML, category sql.NullString
		
		err := rows.Scan(
			&item.ID, &item.Title, &item.Content, &item.Summary, &item.SourceURL,
			&item.Type, &category, &tagsArray, &item.EmbeddingID, &imageURL, &embedHTML, &item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		item.Tags = tagsArray.Elements
		if category.Valid {
			item.Category = category.String
		}
		if imageURL.Valid {
			item.ImageURL = imageURL.String
		}
		if embedHTML.Valid {
			item.EmbedHTML = embedHTML.String
		}
		items = append(items, item)
	}
	
	return items, nil
}

func (r *ItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM items WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// SearchItems performs text search with filters
func (r *ItemRepository) SearchItems(ctx context.Context, filters *models.QueryFilters, limit int) ([]models.Item, error) {
	query := `
		SELECT id, title, content, summary, source_url, type, category, tags, embedding_id, image_url, embed_html, created_at
		FROM items
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	// Text search
	if filters.SearchTerms != "" {
		query += fmt.Sprintf(` AND (
			title ILIKE $%d OR 
			content ILIKE $%d OR 
			summary ILIKE $%d
		)`, argIndex, argIndex, argIndex)
		searchPattern := "%" + filters.SearchTerms + "%"
		args = append(args, searchPattern)
		argIndex++
	}

	// Type filter (only apply if search terms exist, or if type was explicitly set)
	// This allows searching for "video" to find items containing "video" even if type doesn't match
	if filters.Type != "" && filters.SearchTerms != "" {
		// If we have search terms, type filter is optional - search in all types but prefer the specified type
		// We'll handle this in post-processing or make it optional
		// For now, if type is set and search terms exist, we'll search in that type OR in content
		// This is a bit complex, so let's make type filter optional when search terms exist
	}
	if filters.Type != "" && filters.SearchTerms == "" {
		// Only apply type filter if no search terms (pure type filter)
		query += fmt.Sprintf(` AND type = $%d`, argIndex)
		args = append(args, filters.Type)
		argIndex++
	}

	// Date range filter
	if filters.DateFrom != nil {
		query += fmt.Sprintf(` AND created_at >= $%d`, argIndex)
		args = append(args, *filters.DateFrom)
		argIndex++
	}
	if filters.DateTo != nil {
		query += fmt.Sprintf(` AND created_at <= $%d`, argIndex)
		args = append(args, *filters.DateTo)
		argIndex++
	}

	// Tags filter
	if len(filters.Tags) > 0 {
		query += fmt.Sprintf(` AND tags && $%d`, argIndex)
		args = append(args, filters.Tags)
		argIndex++
	}

	// Author filter (search in content)
	if filters.Author != "" {
		query += fmt.Sprintf(` AND (content ILIKE $%d OR title ILIKE $%d)`, argIndex, argIndex)
		authorPattern := "%" + filters.Author + "%"
		args = append(args, authorPattern)
		argIndex++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argIndex)
	args = append(args, limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return []models.Item{}, err
	}
	defer rows.Close()

	items := []models.Item{}
	for rows.Next() {
		var item models.Item
		var tagsArray pgtype.Array[string]
		var imageURL, embedHTML sql.NullString

		err := rows.Scan(
			&item.ID, &item.Title, &item.Content, &item.Summary, &item.SourceURL,
			&item.Type, &tagsArray, &item.EmbeddingID, &imageURL, &embedHTML, &item.CreatedAt,
		)
		if err != nil {
			return []models.Item{}, err
		}

		item.Tags = tagsArray.Elements
		if imageURL.Valid {
			item.ImageURL = imageURL.String
		}
		if embedHTML.Valid {
			item.EmbedHTML = embedHTML.String
		}
		items = append(items, item)
	}

	return items, nil
}

