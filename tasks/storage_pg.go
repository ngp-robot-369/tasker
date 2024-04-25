package tasks

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type storagePG struct {
	db *sql.DB
}

func MakeStoragePG() *storagePG {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASS")
	data := os.Getenv("PG_DATA")
	conf := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, data)

	db, err := sql.Open("postgres", conf)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil
	}

	// db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(WORKERS_POOL_SIZE)
	db.SetMaxOpenConns(WORKERS_POOL_SIZE)

	var (
		versionMin   = 9.5 // UPSERT syntax is used
		versionCheck = true
		dropTable    = true
	)

	if versionCheck {
		var version string
		err = db.QueryRow(`SHOW server_version`).Scan(&version)
		if err != nil {
			log.Panicf("Error getting postgres version: %v", err)
		}
		if ver, err := strconv.ParseFloat(version, 32); err == nil {
			if ver < versionMin {
				log.Panicf("Too old version of postgres: %v, minimal: %v", version, versionMin)
			}
		}
	}
	if dropTable {
		_, err = db.Exec(`DROP TABLE IF EXISTS TaskStatus`)
		if err != nil {
			log.Panicf("Error dropping table TaskStatus: %v", err)
		}
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS TaskStatus (
		id TEXT PRIMARY KEY,
		status TEXT,
		http_status_code INTEGER,
		response_headers JSONB,
		content_length BIGINT
		)`)
	if err != nil {
		log.Panicf("Error creating table TaskStatus: ", err)
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_task_status_id ON TaskStatus(ID)`)
	if err != nil {
		log.Panicf("Error creating index on TaskStatus: ", err)
	}

	return &storagePG{db: db}
}

func (s *storagePG) PutStatus(ctx context.Context, status TaskStatus) error {
	headers, _ := json.Marshal(status.ResponseHeaders)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO TaskStatus (id, status, http_status_code, response_headers, content_length)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE 
			SET status = $2, 
				http_status_code = $3,
				response_headers = $4,
				content_length = $5
		`, status.ID, status.Status, status.HTTPStatusCode, headers, status.ContentLength)

	return err
}

func (s *storagePG) GetStatus(ctx context.Context, id string) (*TaskStatus, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT status, http_status_code, response_headers, content_length FROM TaskStatus
		WHERE id = $1`, id)

	out := TaskStatus{
		ID: id,
	}
	err := row.Scan(&out.Status, &out.HTTPStatusCode, &out.ResponseHeaders, &out.ContentLength)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// -----------------------------------------------------
// Postgres convert interfaces implementation. See:
// - https://pkg.go.dev/database/sql/driver#Valuer
// - https://pkg.go.dev/database/sql#Scanner
// -----------------------------------------------------
func (h httpHeaders) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (a *httpHeaders) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
