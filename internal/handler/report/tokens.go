package report

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken() (string, error) {

	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("AUTH_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func IsValidToken(token string) (bool, error) {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("AUTH_SECRET")), nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}
