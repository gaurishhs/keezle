package utils

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type params struct {
	// memory is the amount of memory to use for hashing in KiB
	memory uint32
	// threads is the number of threads to use for hashing
	threads uint8
	// time is the number of iterations to perform
	time       uint32
	keyLength  uint32
	saltLength uint32
}

func HashPassword(password string) (string, error) {
	p := &params{
		memory:     47104, // 46 MiB
		threads:    2,
		time:       2,
		keyLength:  32,
		saltLength: 16,
	}
	hash, err := generatePasswordHash(password, p)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func generatePasswordHash(password string, p *params) (encodedHash string, err error) {
	salt, err := generateRandomBytes(p.saltLength)

	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.time, p.memory, p.threads, p.keyLength)
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.time, p.threads, base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash))

	return encodedHash, nil
}
