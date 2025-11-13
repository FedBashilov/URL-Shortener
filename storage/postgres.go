package storage

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connString string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")
	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) Save(longURL string) (string, error) {
	shortCode := generateShortCode()

	query := `INSERT INTO urls (short_code, original_url) VALUES ($1, $2) 
	          RETURNING short_code`

	var code string
	err := p.db.QueryRow(query, shortCode, longURL).Scan(&code)

	if err != nil {
		return "", fmt.Errorf("failed to save URL: %w", err)
	}

	return code, nil
}

func (p *PostgresStorage) Get(shortCode string) (string, bool, error) {
	var originalURL string

	err := p.db.QueryRow(
		"SELECT original_url FROM urls WHERE short_code = $1",
		shortCode,
	).Scan(&originalURL)

	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	_, err = p.db.Exec(
		"UPDATE urls SET access_count = access_count + 1 WHERE short_code = $1",
		shortCode,
	)
	if err != nil {
		fmt.Printf("Warning: failed to update counter: %v\n", err)
	}

	return originalURL, true, nil
}

func (p *PostgresStorage) GetStats(shortCode string) (map[string]interface{}, error) {
	var originalURL string
	var accessCount int
	var createdAt time.Time

	err := p.db.QueryRow(
		"SELECT original_url, access_count, created_at FROM urls WHERE short_code = $1",
		shortCode,
	).Scan(&originalURL, &accessCount, &createdAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"original_url": originalURL,
		"access_count": accessCount,
		"created_at":   createdAt,
		"short_code":   shortCode,
	}

	return stats, nil
}

func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, codeLength)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}
