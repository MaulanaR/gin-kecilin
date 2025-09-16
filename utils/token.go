package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

var jwtKey []byte

func SetJWTKey(key []byte) {
	jwtKey = key
}

func GetJWTKey() []byte {
	return jwtKey
}
func ValidateToken(token string) (*Claims, error) {
	claims := &Claims{}
	secret := GetJWTKey()

	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New("Invalid token")
	}
	return claims, nil
}

func HashPassword(password *string) (*string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPwd := string(bytes)
	return &hashedPwd, nil
}

func GenerateToken(email, userID string) (string, string, error) {
	expTime := time.Now().Add(time.Hour * 24).Unix()
	refreshExpTime := time.Now().Add(time.Hour * 24 * 7).Unix()

	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime,
		},
	}

	refClaims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpTime,
		},
	}

	// generate tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAT, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims)
	signedRT, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return signedAT, signedRT, nil
}

func VerifyPassword(inputPwd, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(inputPwd))
	if err != nil {
		return false, err
	}
	return true, nil
}
