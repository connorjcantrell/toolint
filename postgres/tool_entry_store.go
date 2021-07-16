package postgres

import (
	"fmt"

	toolint "github.com/connorjcantrell/toolint"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ToolEntryStore struct {
	*sqlx.DB
}

func (s *ToolEntryStore) ToolEntry(id uuid.UUID) (toolint.ToolEntry, error) {
	var p toolint.ToolEntry
	if err := s.Get(&p, `SELECT * FROM tool_entries WHERE id = $1`, id); err != nil {
		return toolint.ToolEntry{}, fmt.Errorf("error getting tool: %w", err)
	}
	return p, nil
}

func (s *ToolEntryStore) ToolEntriesByUser(userID uuid.UUID) ([]toolint.ToolEntry, error) {
	var tt []toolint.ToolEntry
	var query = `
		SELECT
		FROM tool_entries
		WHERE user_id = $1
		GROUP BY tool_entries.id`
	if err := s.Select(&tt, query, userID); err != nil {
		return []toolint.ToolEntry{}, fmt.Errorf("error getting tool_entries: %w", err)
	}
	return tt, nil
}

func (s *ToolEntryStore) ToolEntries() ([]toolint.ToolEntry, error) {
	var tt []toolint.ToolEntry
	var query = `
		SELECT
		FROM tool_entries
		GROUP BY tool_entries.id`
	if err := s.Select(&tt, query); err != nil {
		return []toolint.ToolEntry{}, fmt.Errorf("error getting tool_entries: %w", err)
	}
	return tt, nil
}

func (s *ToolEntryStore) CreateToolEntry(t *toolint.ToolEntry) error {
	if err := s.Get(t, `INSERT INTO tool_entries VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		t.ID,
		t.UserID,
		t.ToolID,
		t.Condition); err != nil {
		return fmt.Errorf("error creating tool: %w", err)
	}
	return nil
}

func (s *ToolEntryStore) UpdateToolEntry(t *toolint.ToolEntry) error {
	if err := s.Get(t, `UPDATE tool_entries SET name = $1, model = $2, price = $3, category_id = $4 WHERE id = $5 RETURNING *`,
		t.ID,
		t.UserID,
		t.ToolID,
		t.Condition); err != nil {
		return fmt.Errorf("error updating tool: %w", err)
	}
	return nil
}

func (s *ToolEntryStore) DeleteToolEntry(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM tool_entries WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting tool: %w", err)
	}
	return nil
}
