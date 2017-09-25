package controllers

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// password for signing signature of the JWT
var tokenEncode = []byte(global.Config.JWTtokenPassword)

const expire = 60 * 60 * 24 * 7 // token expires in one week

// User ...
type User struct {
	ID       string
	Username string
}

// CreateToken ...
func CreateToken(user User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["userInfo"] = User{ID: user.ID, Username: user.Username}
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token.Claims = claims

	signedString, err := token.SignedString(tokenEncode)
	if err != nil {
		return "", err
	}

	return signedString, nil
}
