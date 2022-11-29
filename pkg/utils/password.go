package utils

import "golang.org/x/crypto/bcrypt"

// 哈希密码
func HashAndSalt(pwd []byte) string {
	if hash, err := bcrypt.GenerateFromPassword(pwd,
		bcrypt.MinCost); err != nil {
		return ""
	} else {
		return string(hash)
	}
}

// 验证密码
func ComparePassword(hashedPwd, plainPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd),
		[]byte(plainPwd)); err != nil {
		return false
	}
	return true
}
