package models

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email          string `json:"Email"`
	HashedPassword []byte `json:"HashedPassword"`
	Name           string `json:"Name"`
}

func (u User) Validate() error {
	if len(u.Email) <= 0 {
		return errors.New("email can't be blank")
	}

	if len(u.HashedPassword) <= 0 {
		return errors.New("password can't be blank")
	}

	return nil
}

func (u *User) HashPassword(password string) ([]byte, error) {
	if len(password) <= 0 {
		return nil, errors.New("password can't be blank")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func (u *User) VerifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u User) GenerateToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}
