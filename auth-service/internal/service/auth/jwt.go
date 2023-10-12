package auth

import (
	"auth-service/internal/model"
	"auth-service/internal/repository/user"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type userForJwt struct {
	id        string
	username  string
	email     string
	roles     []string
	createdAt time.Time
	updatedAt time.Time
}

type claims struct {
	user *userForJwt
	jwt.RegisteredClaims
}

func ToUserFromRepository(user *user.User) *userForJwt {
	return &userForJwt{
		id:        user.UUID,
		username:  user.Username,
		email:     user.Email,
		roles:     user.Roles,
		createdAt: user.CreatedAt,
		updatedAt: user.UpdatedAt,
	}
}

func ToUserFromModel(user *model.User) *userForJwt {
	return &userForJwt{
		id:        user.UUID,
		username:  user.Username,
		email:     user.Email,
		roles:     user.Roles,
		createdAt: user.CreatedAt,
		updatedAt: user.UpdatedAt,
	}
}

var jwtKey = []byte("my_secret_key")

func GenerateJWTToken(user *userForJwt) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claim := &claims{
		user: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   user.id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("generating token: %s", err)
	}

	return tokenString, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
