package tokens

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

const DefaultExpiration = time.Hour * 24

func New(id string) (string, error) {
	issuedAt := time.Now()
	claims := &jwt.StandardClaims{
		Audience:  "Authorization",
		ExpiresAt: issuedAt.Add(DefaultExpiration).Unix(),
		Id:        id,
		IssuedAt:  issuedAt.Unix(),
		Issuer:    "go-subscriptions-workflows",
		Subject:   id,
	}
	return newToken(claims)
}

func newToken(claims *jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetSecretKey())
}

func GetSecretKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

type TokenPayload struct {
	Issuer    string
	Audience  string
	UserID    string
	Subject   string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func Parse(tokenString string) (*TokenPayload, error) {
	claims := new(jwt.StandardClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, onParse)
	if err != nil {
		return nil, err
	}
	var ok bool
	claims, ok = token.Claims.(*jwt.StandardClaims)
	if !token.Valid || !ok {
		return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorMalformed)
	}
	return &TokenPayload{
		Issuer:    claims.Issuer,
		Audience:  claims.Audience,
		UserID:    claims.Id,
		Subject:   claims.Subject,
		IssuedAt:  time.Unix(claims.IssuedAt, 0),
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
	}, nil
}

func onParse(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return GetSecretKey(), nil
}
