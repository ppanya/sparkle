package sparklecrypto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
	"strings"
	"time"
)

type TokenSigner interface {
	Sign(value interface{}) (string, error)
	Parse(token string, value interface{}) error
}

type Base64TokenSigner struct{}

func NewBase64TokenSigner() *Base64TokenSigner {
	return &Base64TokenSigner{}
}

func (p *Base64TokenSigner) Sign(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(b)
	trailer := xid.New()
	return fmt.Sprintf("%s::%s", string(encoded), trailer.String()), err
}

func (p *Base64TokenSigner) Parse(token string, value interface{}) error {
	seq := strings.Split(token, "::")
	if len(seq) != 2 {
		return fmt.Errorf("invalid token format")
	}
	decoded, err := base64.StdEncoding.DecodeString(seq[0])
	if err != nil {
		return err
	}
	return json.Unmarshal(decoded, value)
}

type JWTTokenSigner struct {
	Secret []byte
}

func (j *JWTTokenSigner) Parse(token string, value interface{}) error {
	t, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		return j.Secret, nil
	})

	if err != nil {
		return err
	}

	if claim, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return json.Unmarshal([]byte(claim["data"].(string)), value)
	} else {
		return fmt.Errorf("invalid credentials")
	}

}

func NewJWT(secret string) *JWTTokenSigner {
	return &JWTTokenSigner{
		Secret: []byte(secret),
	}
}

func (j *JWTTokenSigner) Sign(value interface{}) (string, error) {
	vv, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	c := jwt.MapClaims{
		"data": string(vv),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString(j.Secret)
	return tokenString, err
}
