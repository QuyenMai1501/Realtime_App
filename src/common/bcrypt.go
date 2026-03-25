package common

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {
	//const salt = 2^14
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	return string(bytes), err
}

func CheckPassword(hashedPw, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(pw))
	return err == nil
}