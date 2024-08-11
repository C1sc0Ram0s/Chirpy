package database

import (
	"time"
)

type RefreshToken struct {
	UserID     int       `json:"user_id"`
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

func (db *DB) StoreRefreshToken(userID int, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		UserID:     userID,
		Token:      token,
		Expiration: time.Now().Add(time.Hour),
	}
	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRefreshTokenUser(token string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	refreshToken, exists := dbStructure.RefreshTokens[token]
	if !exists {
		return User{}, ErrNotExist
	}

	if refreshToken.Expiration.Before(time.Now()) {
		return User{}, ErrNotExist
	}

	user, err := db.GetUser(refreshToken.UserID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) RevokeToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.RefreshTokens, token)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
