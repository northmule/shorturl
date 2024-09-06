package user

import (
	"crypto/sha512"
	"fmt"
)

// PasswordHash хэш пароля согласно общим правилом хранения в бд
func PasswordHash(password string) string {
	hashAlg512 := sha512.New()
	hashAlg512.Write([]byte(password))
	return fmt.Sprintf("%x", hashAlg512.Sum(nil))
}
