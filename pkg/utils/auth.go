package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

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

type mClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

// 创建token
func CreateToken(username string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, mClaims{
		jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 3600*viper.GetInt64("app_jwt_timeout"),
		},
		username,
	}).SignedString([]byte(viper.GetString("app_jwt_key")))
}

// 验证token
func ParseToken(tokestr string) (*mClaims, error) {
	token, err := jwt.ParseWithClaims(tokestr, &mClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("app_jwt_key")), nil
	})
	if claims, ok := token.Claims.(*mClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
