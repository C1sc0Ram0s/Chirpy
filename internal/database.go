package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}
type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// NewDB creates a new DB connection
// and creates a new database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	mux := sync.RWMutex{}
	db := DB{
		path: path,
		mux:  &mux,
	}

	// Check if DB exists, else create one
	_, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		err = db.ensureDB()
		if err != nil {
			return nil, err
		}
	}

	return &db, nil
}

// Create a new Chirp and saves to disk/file
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	delete(dbStruct.Chirps, 0)

	// Find max ID
	id := 1
	for k := range dbStruct.Chirps {
		if k >= id {
			id = k + 1
		}
	}

	chirp := Chirp{
		Id:   id,
		Body: body,
	}
	dbStruct.Chirps[chirp.Id] = chirp

	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

// GetChirps returns all chirps in the DB
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := []Chirp{}
	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})
	return chirps, nil
}

// ensureDB creates a new DB file if it doesn't exist
func (db *DB) ensureDB() error {
	emptyDB := DBStructure{Chirps: make(map[int]Chirp)}
	return db.writeDB(emptyDB)
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile(db.path)
	if err != nil {
		// If the file doesn't exist, return an empty strut with initialized map
		if errors.Is(err, os.ErrNotExist) {
			return DBStructure{Chirps: make(map[int]Chirp)}, nil
		} else {
			return DBStructure{}, err
		}
	}

	var dbStruct DBStructure
	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		return DBStructure{}, nil
	}

	return dbStruct, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	dbStructData, err := json.Marshal(&dbStructure)
	if err != nil {
		return fmt.Errorf("failed to marshal db struct: %w", err)
	}

	err = os.WriteFile(db.path, dbStructData, 0666)
	if err != nil {
		return fmt.Errorf("failed to write db file: %w", err)
	}
	return nil
}
