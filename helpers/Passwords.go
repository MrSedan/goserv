package helpers

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	PBKDF2_HASH_ALGORITHM string = "sha512"
	PBKDF2_ITERATIONS     int    = 15000
	SCRYPT_N              int    = 32768
	SCRYPT_R              int    = 8
	SCRYPT_P              int    = 1

	SALT_BYTES int = 64
	HASH_BYTES int = 64
)
const (
	HASH_SECTIONS        int = 4
	HASH_ALGORITHM_INDEX int = 0
	HASH_ITERATION_INDEX int = 1
	HASH_SALT_INDEX      int = 2
	HASH_PBKDF2_INDEX    int = 3
	HASH_SCRYPT_R_INDEX  int = 4
	HASH_SCRYPT_P_INDEX  int = 5
)

func CreateHash(password string) (string, error) {
	salt := make([]byte, SALT_BYTES)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	var hash []byte
	hash = pbkdf2.Key([]byte(password), salt, PBKDF2_ITERATIONS, HASH_BYTES, sha512.New)
	return fmt.Sprintf(
		"%s:%d:%s:%s", PBKDF2_HASH_ALGORITHM, PBKDF2_ITERATIONS,
		base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(hash),
	), err
}

func ValidatePassword(password, correctHash string) bool {
	params := strings.Split(correctHash, ":")
	if len(params) < HASH_SECTIONS {
		return false
	}
	it, err := strconv.Atoi(params[HASH_ITERATION_INDEX])
	if err != nil {
		return false
	}
	salt, err := base64.StdEncoding.DecodeString(params[HASH_SALT_INDEX])
	if err != nil {
		return false
	}
	hash, err := base64.StdEncoding.DecodeString(params[HASH_PBKDF2_INDEX])
	if err != nil {
		return false
	}
	var testHash []byte
	testHash = pbkdf2.Key([]byte(password), salt, it, len(hash), sha512.New)
	return subtle.ConstantTimeCompare(hash, testHash) == 1
}
