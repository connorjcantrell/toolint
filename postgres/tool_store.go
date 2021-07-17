package postgres

import (
	"fmt"

	"github.com/connorjcantrell/toolint"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ToolStore struct {
	*sqlx.DB
}

func (s *ToolStore) Tool(id uuid.UUID) (toolint.Tool, error) {
	var p toolint.Tool
	if err := s.Get(&p, `SELECT * FROM tools WHERE id = $1`, id); err != nil {
		return toolint.Tool{}, fmt.Errorf("error getting tool: %w", err)
	}
	return p, nil
}

func (s *ToolStore) ToolsByCategory(category string) ([]toolint.Tool, error) {
	var tt []toolint.Tool
	var query = `
		SELECT
		FROM tools
		WHERE category_id = $1
		GROUP BY tools.id`
	if err := s.Select(&tt, query, category); err != nil {
		return []toolint.Tool{}, fmt.Errorf("error getting tools: %w", err)
	}
	return tt, nil
}

func (s *ToolStore) Tools() ([]toolint.Tool, error) {
	var tt []toolint.Tool
	var query = `
		SELECT
		FROM tools
		GROUP BY tools.id,
		ORDER BY name DESC`
	if err := s.Select(&tt, query); err != nil {
		return []toolint.Tool{}, fmt.Errorf("error getting tools: %w", err)
	}
	return tt, nil
}

func (s *ToolStore) CreateTool(t *toolint.Tool) error {
	if err := s.Get(t, `INSERT INTO tools VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		t.ID,
		t.Name,
		t.Model,
		t.Make,
		t.Category); err != nil {
		return fmt.Errorf("error creating tool: %w", err)
	}
	return nil
}

func (s *ToolStore) UpdateTool(t *toolint.Tool) error {
	if err := s.Get(t, `UPDATE tools SET name = $1, model = $2, make = $3, category_id = $4 WHERE id = $5 RETURNING *`,
		t.Name,
		t.Model,
		t.Make,
		t.Category,
		t.ID); err != nil {
		return fmt.Errorf("error updating tool: %w", err)
	}
	return nil
}

func (s *ToolStore) DeleteTool(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM tools WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting tool: %w", err)
	}
	return nil
}
