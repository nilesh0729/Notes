package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	minSecretKeySize = 32
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretkey string) (Maker, error) {
	if len(secretkey) < minSecretKeySize {
		return nil, fmt.Errorf("the size of the secret key must be  atleast %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretkey}, nil
}

func (maker *JWTMaker) CreateToken(username string, Duration time.Duration) (string, error){

	payload, err := NewPayload(username,Duration)
	if err != nil{
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string)(*Payload, error){
	keyFunc := func(token *jwt.Token)(interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwttoken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil{
		if errors.Is(err, jwt.ErrTokenExpired){
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	Payload, ok := jwttoken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return Payload, nil
}
