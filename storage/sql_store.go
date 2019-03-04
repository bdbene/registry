package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/satori/go.uuid"

	_ "github.com/mattn/go-sqlite3"
)

// SqlStore stores committed values in an append-only log,
// caches proposed changes in memory.
type SqlStore struct {
	location       string
	proposedCache  map[string]string
	committedStore *sql.DB
	mutex          *sync.Mutex
}

const createStatement string = "CREATE TABLE IF NOT EXISTS schemas (name TEXT, version TEXT, schema TEXT, PRIMARY KEY (name, version))"
const insertSchema string = "INSERT INTO schemas (name, version, schema) VALUES (?, ?, ?)"
const queryVersions string = "SELECT version FROM schemas WHERE name=?"
const querySchema string = "SELECT name, schema FROM schemas WHERE version=?"

// NewSqlStore creates a data store that persists committed data to disk
// in an append only log.
func NewSqlStore(config *StorageConfig) (*SqlStore, error) {
	location := "./data/database.db"

	db, err := sql.Open("sqlite3", location)
	if err != nil {
		return nil, &StorageError{err.Error()}
	}

	statement, err := db.Prepare(createStatement)
	if err != nil {
		return nil, &StorageError{err.Error()}
	}

	_, err = statement.Exec()
	if err != nil {
		return nil, &StorageError{err.Error()}
	}

	mutex := &sync.Mutex{}

	return &SqlStore{location, make(map[string]string, 20), db, mutex}, nil
}

// Propose a value by saving it in memory before proposing it to other nodes.
func (db *SqlStore) Propose(key, value string) error {
	db.mutex.Lock()
	db.proposedCache[key] = value
	db.mutex.Unlock()

	return nil
}

// Commit an accepted proposed change to disk.
func (db *SqlStore) Commit(key string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	val, ok := db.proposedCache[key]

	if !ok {
		return &StorageError{fmt.Sprintf("Cannot commit key '%s' as it has not been proposed.", key)}
	}

	version, err := uuid.NewV4()
	if err != nil {
		return &StorageError{err.Error()}
	}

	_, err = db.committedStore.Exec(insertSchema, key, version, val)
	if err != nil {
		return &StorageError{err.Error()}
	}

	delete(db.proposedCache, key)

	log.Printf("Schema written to disk; %s, %s: %s\n", key, version, val)
	return nil
}

// LookupVersions retrieves all versions of a schema.
func (db *SqlStore) LookupVersions(key string) ([]string, error) {
	rows, err := db.committedStore.Query(queryVersions, key)
	if err != nil {
		return nil, &StorageError{err.Error()}
	}

	defer rows.Close()
	versions := make([]string, 1)

	for rows.Next() {
		var version string

		err = rows.Scan(&version)
		if err != nil {
			return nil, &StorageError{err.Error()}
		}

		versions = append(versions, version)
	}

	return versions, nil
}

// Lookup a schema based on its name and version.
func (db *SqlStore) Lookup(key string, version string) (string, error) {
	row := db.committedStore.QueryRow(querySchema, key, version)

	var schema string

	err := row.Scan(&schema)
	if err != nil {
		return "", &StorageError{err.Error()}
	}

	return schema, nil
}

// Close the database.
func (db *SqlStore) Close() {
	db.committedStore.Close()
}
