package tokens

import "time"

type Maker interface {
	CreateToken(username string, Duration time.Duration) (string, error)
	VerifyToken(token string)(*Payload, error)
}
