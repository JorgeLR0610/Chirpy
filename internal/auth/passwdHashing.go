package auth

import (
	"github.com/alexedwards/argon2id"
)

const dummyHash = "$argon2id$v=19$m=65536,t=1,p=12$9+GvJhOgkFRakeSNNQ31fQ$lm/Es7TXa+Kqk+RhtlqzSn1mcAYO8jd04qfSzgk48eA"

func HashPassword(password string) (string, error) {
	hashedPasswd, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashedPasswd, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}



func SimulatePasswordCheck(password string) {
	_, _ = argon2id.ComparePasswordAndHash(password, dummyHash)
}
