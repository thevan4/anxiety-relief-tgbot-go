package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/db/models"
)

type DBWork interface {
	GetAllUsers() ([]*models.User, error)
	InsertOrUpdateUsers(users []*models.User) error
	LogTechniqueStart(userID int64, techniqueID, category string) error
	LogTechniqueComplete(userID int64, techniqueID string) error
	Close() error
}

type PostgresDB struct {
	db *sql.DB
}

func MustGetPostgresDatabase(dataSourceName string, migrationFiles []string) (DBWork, func() error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("failed to open postgres: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}

	// Применяем миграции
	for _, migrationFile := range migrationFiles {
		migrationSQL, err := os.ReadFile(migrationFile)
		if err != nil {
			log.Fatalf("failed to read migration file %s: %v", migrationFile, err)
		}

		if _, err = db.Exec(string(migrationSQL)); err != nil {
			log.Fatalf("failed to execute migration %s: %v", migrationFile, err)
		}
	}

	log.Println("Database initialized successfully")

	return &PostgresDB{db: db}, db.Close
}

func (p *PostgresDB) GetAllUsers() ([]*models.User, error) {
	rows, err := p.db.Query(`
		SELECT user_id, username, first_seen, last_seen 
		FROM users
	`)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.FirstSeen, &u.LastSeen); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &u)
	}

	return users, nil
}

func (p *PostgresDB) InsertOrUpdateUsers(users []*models.User) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO users (user_id, username, first_seen, last_seen)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET
			username = EXCLUDED.username,
			last_seen = EXCLUDED.last_seen
	`)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, u := range users {
		if _, err := stmt.Exec(u.UserID, u.Username, u.FirstSeen, u.LastSeen); err != nil {
			return fmt.Errorf("insert/update user %d: %w", u.UserID, err)
		}
	}

	return tx.Commit()
}

func (p *PostgresDB) LogTechniqueStart(userID int64, techniqueID, category string) error {
	_, err := p.db.Exec(`
		INSERT INTO technique_usage (user_id, technique_id, category, started_at)
		VALUES ($1, $2, $3, NOW())
	`, userID, techniqueID, category)

	if err != nil {
		return fmt.Errorf("log technique start: %w", err)
	}

	return nil
}

func (p *PostgresDB) LogTechniqueComplete(userID int64, techniqueID string) error {
	_, err := p.db.Exec(`
		UPDATE technique_usage
		SET completed = TRUE, completed_at = NOW()
		WHERE user_id = $1 AND technique_id = $2
		AND id = (
			SELECT id FROM technique_usage
			WHERE user_id = $1 AND technique_id = $2
			ORDER BY started_at DESC
			LIMIT 1
		)
	`, userID, techniqueID)

	if err != nil {
		return fmt.Errorf("log technique complete: %w", err)
	}

	return nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}
