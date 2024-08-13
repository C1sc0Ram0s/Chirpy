package database

import (
	"errors"
)

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (db *DB) CreateChirp(userID int, body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := 1
	for _, chirp := range dbStructure.Chirps {
		if chirp.ID >= id {
			id = chirp.ID + 1
		}
	}

	chirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorID: userID,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id, userID int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, exists := dbStructure.Chirps[id]
	if !exists {
		return err
	}

	if chirp.AuthorID == userID {
		delete(dbStructure.Chirps, id)
		err = db.writeDB(dbStructure)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("unauthorized delete")

}
