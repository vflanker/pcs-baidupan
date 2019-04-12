package webjwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/iikira/BaiduPCS-Go/pcsutil/jwted25519"
	"strconv"
	"time"
)

func (wj *WebJWT) KeyFunc(token *jwt.Token) (interface{}, error) {
	wj.lazyInit()
	return wj.ed25519PubKey, nil
}

// StandardClaims jwt 标准声明
func (wj *WebJWT) StandardClaims() *jwt.StandardClaims {
	var (
		nowTime = time.Now()
		nowStr  = strconv.FormatInt(nowTime.UnixNano(), 10)
	)
	return &jwt.StandardClaims{
		ExpiresAt: nowTime.Add(wj.Expires).Unix(),
		Id:        nowStr,
	}
}

// Sign 签名
func (wj *WebJWT) Sign(claims jwt.Claims) string {
	wj.lazyInit()
	token := jwt.NewWithClaims(jwted25519.SigningMethodED25519, claims)
	tokenString, err := token.SignedString(wj.ed25519PrivKey)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func (wj *WebJWT) SignStandardClaims() string {
	return wj.Sign(wj.StandardClaims())
}

// Verify 验证签名
func (wj *WebJWT) Verify(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, wj.KeyFunc)
	return err == nil && token.Valid
}
