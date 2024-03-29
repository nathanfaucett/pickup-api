package jwt

import (
	"encoding/json"
	"fmt"

	"github.com/aicacia/pickup/app/config"
	"github.com/dgrijalva/jwt-go/v4"
)

type ToMapClaims interface {
	ToMapClaims() (jwt.MapClaims, error)
}

func anyToMapClaims(value any) (jwt.MapClaims, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var result jwt.MapClaims
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type Claims struct {
	Subject          int32 `json:"sub" validate:"required"`
	NotBeforeSeconds int64 `json:"nbf" validate:"required"`
	IssuedAtSeconds  int64 `json:"iat" validate:"required"`
	ExpiresAtSeconds int64 `json:"exp" validate:"required"`
}

func (claims *Claims) ToMapClaims() (jwt.MapClaims, error) {
	return anyToMapClaims(claims)
}

func ParseClaimsFromToken[C Claims](tokenString string) (*C, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(config.Get().JWT.Secret), nil
	}, jwt.WithoutAudienceValidation())
	if err != nil {
		return nil, err
	}
	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		bytes, err := json.Marshal(&mapClaims)
		if err != nil {
			return nil, err
		}
		var claims C
		if err := json.Unmarshal(bytes, &claims); err != nil {
			return nil, err
		}
		return &claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func CreateToken[C ToMapClaims](claims C) (string, error) {
	mapClaims, err := claims.ToMapClaims()
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenString, err := token.SignedString([]byte(config.Get().JWT.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
