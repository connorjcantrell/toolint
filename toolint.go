package toolint

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
}

type Tool struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Model    string    `db:"model"`
	Make     string    `db:"make"`
	Category string    `db:"category"`
}

type ToolEntry struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	ToolID    uuid.UUID `db:"post_id"`
	Condition string    `db:"condition"`
}

type Store interface {
	UserStore
	ToolStore
	ToolEntryStore
}

type UserStore interface {
	User(id uuid.UUID) (User, error)
	UserByUsername(username string) (User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id uuid.UUID) error
}
type ToolStore interface {
	Tool(id uuid.UUID) (Tool, error)
	Tools() ([]Tool, error)
	ToolsByCategory(category string) ([]Tool, error)
	CreateTool(t Tool) (Tool, error)
	UpdateTool(t Tool) (Tool, error)
	DeleteTool(id uuid.UUID) error
}

type ToolEntryStore interface {
	ToolEntry(id uuid.UUID) (ToolEntry, error)
	ToolEntries() ([]ToolEntry, error)
	ToolEntriesByUser(userID uuid.UUID) ([]ToolEntry, error)
	CreateToolEntry(t *ToolEntry) error
	UpdateToolEntry(t *ToolEntry) error
	DeleteToolEntry(id uuid.UUID) error
}
