package webjwt

import (
	cryptorand "crypto/rand"
	"golang.org/x/crypto/ed25519"
	"time"
)

const (
	// ExpriesToken token超时时间
	ExpriesToken = 24 * time.Hour
)

type (
	// WebJWT jwt 对象
	WebJWT struct {
		Expires        time.Duration
		ed25519PubKey  ed25519.PublicKey
		ed25519PrivKey ed25519.PrivateKey
	}
)

// NewWebJWT 初始化jwt
func NewWebJWT() *WebJWT {
	return &WebJWT{
		Expires: ExpriesToken,
	}
}

func (wj *WebJWT) lazyInit() {
	if wj.ed25519PubKey == nil || wj.ed25519PrivKey == nil {
		wj.genKey()
	}
}

func (wj *WebJWT) genKey() {
	var err error
	wj.ed25519PubKey, wj.ed25519PrivKey, err = ed25519.GenerateKey(cryptorand.Reader)
	if err != nil {
		panic(err)
	}
}
