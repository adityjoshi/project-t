package repository

import (
	"context"
	"synapse/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"
)

type RelationRepository struct {
	pool *pgxpool.Pool
}

func NewRelationRepository(pool *pgxpool.Pool) *RelationRepository {
	return &RelationRepository{pool: pool}
}

func (r *RelationRepository) Create(ctx context.Context, itemID, relatedItemID uuid.UUID, score float64) error {
	query := `
		INSERT INTO item_relations (item_id, related_item_id, similarity_score)
		VALUES ($1, $2, $3)
		ON CONFLICT (item_id, related_item_id) 
		DO UPDATE SET similarity_score = $3, created_at = NOW()
	`
	_, err := r.pool.Exec(ctx, query, itemID, relatedItemID, score)
	return err
}

func (r *RelationRepository) GetRelated(ctx context.Context, itemID uuid.UUID, limit int) ([]models.RelatedItem, error) {
	query := `
		SELECT i.id, i.title, i.content, i.summary, i.source_url, i.type, i.tags, i.embedding_id, i.created_at, ir.similarity_score
		FROM item_relations ir
		JOIN items i ON ir.related_item_id = i.id
		WHERE ir.item_id = $1
		ORDER BY ir.similarity_score DESC
		LIMIT $2
	`
	
	rows, err := r.pool.Query(ctx, query, itemID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var related []models.RelatedItem
	for rows.Next() {
		var rel models.RelatedItem
		var tagsArray pgtype.Array[string]
		
		err := rows.Scan(
			&rel.Item.ID, &rel.Item.Title, &rel.Item.Content, &rel.Item.Summary,
			&rel.Item.SourceURL, &rel.Item.Type, &tagsArray, &rel.Item.EmbeddingID,
			&rel.Item.CreatedAt, &rel.SimilarityScore,
		)
		if err != nil {
			return nil, err
		}
		
		rel.Item.Tags = tagsArray.Elements
		related = append(related, rel)
	}
	
	return related, nil
}

