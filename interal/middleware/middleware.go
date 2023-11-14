package middleware

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type claims struct {
	jwt.StandardClaims
	telegramUserID int64
}

func JwtHashing(password string, telegramUserID int64) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(time.Hour * 24 * 365)},
		telegramUserID: telegramUserID,
	})
	token, err := t.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetUserID(token string) (int64, error) {
	claims := &claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return 0, fmt.Errorf("jwt.ParseWithClaims failed err: %v", err)
	}
	return claims.telegramUserID, nil
}

func tokenIsValid(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return fmt.Errorf("some problem token, jwt.Parse(tokenStr...) falied %v", err)
	}
	if !token.Valid {
		return errors.New("token is invalid")
	}
	return nil
}
