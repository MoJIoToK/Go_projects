// Package sqlite provides an SQLite implementation of the storage.Storage interface.
// It uses foreign keys and transactions to maintain db integrity.
package sqlite

import (
	"database/sql"
	"fmt"
	"limits-app/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// Storage is an SQLite implementation of the storage.Storage interface.
// It manages connections to an SQLite database and executes queries within transactions.
type Storage struct {
	db *sql.DB
}

// NewStorage creates a new SQLite storage instance and initializes the schema.
// The provided path must be writable (e.g., "./db/limits.db").
// It enables foreign key constraints and creates tables if they don't exist.
func NewStorage(storagePath string) (*Storage, error) {
	const op = "database.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	_, err = db.Exec(`
	PRAGMA foreign_keys = ON;
	CREATE TABLE IF NOT EXISTS clients (
		client_code TEXT PRIMARY KEY
	);
	CREATE Table IF NOT EXISTS positions (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    client_code TEXT,
	    seccode TEXT,
	    open_balance REAL,
		FOREIGN KEY (client_code) REFERENCES clients(client_code)
	    ON DELETE CASCADE
	    ON UPDATE CASCADE
	);
	CREATE Table IF NOT EXISTS limits (
	    position_id INTEGER,
	    limit_kind INTEGER,
	    open_limit REAL,
	    FOREIGN KEY (position_id) REFERENCES positions(id)
	    ON DELETE CASCADE
	    ON UPDATE CASCADE
	);
`)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{db: db}, nil
}

// Save persists normalized client limits using a single transaction.
// It inserts clients (ignoring duplicates), then positions, then limits.
// On error, the transaction is rolled back automatically.
func (s Storage) Save(limits []models.ClientLimits) error {
	const op = "database.sqlite.SaveLimits"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	defer tx.Rollback()

	for _, client := range limits {
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO clients (client_code) VALUES (?)",
			client.ClientCode,
		)
		if err != nil {
			return fmt.Errorf("%s: insert client %s: %w", op, client.ClientCode, err)
		}

		for _, pos := range client.Positions {
			res, err := tx.Exec(
				"INSERT INTO positions (client_code, seccode, open_balance) VALUES (?, ?, ?)",
				client.ClientCode, pos.SecId, pos.OpenBalance,
			)
			if err != nil {
				return fmt.Errorf(
					"%s: insert position %s for client %s: %w",
					op, pos.SecId, client.ClientCode, err)
			}

			posID, err := res.LastInsertId()
			if err != nil {
				return fmt.Errorf("%s: get position ID: %w", op, err)
			}

			for _, limit := range pos.Limits {
				_, err := tx.Exec(
					"INSERT INTO limits (position_id, limit_kind, open_limit) VALUES (?, ?, ?)",
					posID, limit.LimitKind, limit.OpenLimit,
				)
				if err != nil {
					return fmt.Errorf(
						"%s: insert limit %d for position %d: %w",
						op, limit.LimitKind, posID, err)
				}
			}
		}

	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	return nil
}

// LoadAll retrieves all normalized db by joining clients → positions → limits.
// It reconstructs the hierarchical model (Client → Positions → Limits) from flat tables.
// Returns an empty slice if no db exists.
func (s Storage) LoadAll() ([]models.ClientLimits, error) {
	const op = "storage.sqlite.LoadAll"

	clientsRows, err := s.db.Query("SELECT * FROM clients")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer clientsRows.Close()

	var result []models.ClientLimits
	for clientsRows.Next() {
		var clientCode string
		if err := clientsRows.Scan(&clientCode); err != nil {
			return nil, fmt.Errorf("%s: scan client: %w", op, err)
		}

		posRows, err := s.db.Query("SELECT * FROM positions WHERE client_code = ?",
			clientCode)
		if err != nil {
			return nil, fmt.Errorf("%s: select positions: %w", op, err)
		}

		var positions []models.InstrumentPosition
		for posRows.Next() {
			var id int
			var clientCodeFromDB string
			var SecId string
			var openBalance float64
			if err := posRows.Scan(&id, &clientCodeFromDB, &SecId, &openBalance); err != nil {
				return nil, fmt.Errorf("%s: scan positions: %w", op, err)
			}

			limitsRows, err := s.db.Query(
				"SELECT limit_kind, open_limit FROM limits WHERE position_id = ?",
				id)
			if err != nil {
				return nil, fmt.Errorf("%s: select limits: %w", op, err)
			}

			limits := make(map[int]models.DepoLimit)
			for limitsRows.Next() {
				var kind int
				var openLimit float64
				if err := limitsRows.Scan(&kind, &openLimit); err != nil {
					return nil, fmt.Errorf("%s: scan limit: %w", op, err)
				}
				limits[kind] = models.DepoLimit{
					ClientCode:  clientCode,
					SecCode:     SecId,
					LimitKind:   kind,
					OpenLimit:   openLimit,
					OpenBalance: openBalance,
				}
			}
			limitsRows.Close()

			positions = append(positions, models.InstrumentPosition{
				SecId:       SecId,
				OpenBalance: openBalance,
				Limits:      limits,
			})
		}
		posRows.Close()

		result = append(result, models.ClientLimits{
			ClientCode: clientCode,
			Positions:  positions,
		})
	}

	return result, nil
}
