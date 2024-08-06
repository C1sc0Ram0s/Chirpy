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
	Users  map[int]User  `json:"users"`
}
type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
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

// Create a new Chirp and saves to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

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

// Create a new user and save to disk
func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Find max ID
	id := 1
	for k := range dbStruct.Users {
		if k >= id {
			id = k + 1
		}
	}
	user := User{
		ID:    id,
		Email: email,
	}
	dbStruct.Users[user.ID] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetChirps returns all chirps in the DB
func (db *DB) GetChirps(chirpId ...int) ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	chirps := []Chirp{}

	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	if chirpId != nil {
		if _, exists := dbStruct.Chirps[chirpId[0]]; !exists {
			return []Chirp{}, errors.New("chirp does not exist")
		} else {
			chirps = append(chirps, dbStruct.Chirps[chirpId[0]])
		}

		return chirps, nil
	}

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
	emptyDB := DBStructure{Chirps: make(map[int]Chirp), Users: make(map[int]User)}
	return db.writeDB(emptyDB)
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile(db.path)
	var dbStruct DBStructure
	if err != nil {
		// If the file doesn't exist, return an empty strut with initialized map
		if errors.Is(err, os.ErrNotExist) {
			dbStruct = DBStructure{
				Chirps: make(map[int]Chirp),
				Users:  make(map[int]User),
			}
			return dbStruct, nil
		} else {
			return dbStruct, err
		}
	}

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
