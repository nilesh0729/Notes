package tokens

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto        *paseto.V2
	symmetrickKey []byte
}

func NewPasetoMaker(symmetric_key string) (Maker, error) {
	if len(symmetric_key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Invalid Key Size: keySize must be %d", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:        paseto.NewV2(),
		symmetrickKey: []byte(symmetric_key),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration)(string, error){
	payload, err := NewPayload(username, duration)
	if err != nil{
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetrickKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string)(*Payload, error){
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetrickKey, payload, nil)
	if err != nil{
		return nil, ErrInvalidToken
	}
	err = payload.valid()
	if err != nil{
		return nil, ErrExpiredToken
	}

	return payload, nil
	
}
