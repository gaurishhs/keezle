package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("invalid hash format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
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

// HashPassword hashes a password using Argon2id using the following parameters:
// - memory: 46 MiB
// - threads: 2
// - time: 2 iterations
// - keyLength: 32 bytes
// - saltLength: 16 bytes
// The resulting hash is encoded in a format compatible with the Argon2 specification.
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

func ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.time, p.memory, p.threads, p.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.time, &p.threads)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}

func generatePasswordHash(password string, p *params) (encodedHash string, err error) {
	salt, err := GenerateRandomBytes(int(p.saltLength))

	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.time, p.memory, p.threads, p.keyLength)
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.time, p.threads, base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash))

	return encodedHash, nil
}
