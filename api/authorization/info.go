package authorization

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt"
)

type BasicToken struct {
	data string
}

func (b *BasicToken) parse(rawToken string) error {
	b.data = rawToken
	return nil
}

func (b *BasicToken) get() []byte {
	return []byte(b.data)
}

func (b *BasicToken) UserId() string {
	return ""
}

type JwtToken struct {
	data   string
	parsed *jwt.Token
	parts  []string
}

func (j *JwtToken) parse(rawToken string) error {

	j.data = rawToken

	parsed, parts, err := new(jwt.Parser).ParseUnverified(j.data, jwt.MapClaims{})

	if err != nil {
		return err
	}

	j.parts = parts
	j.parsed = parsed

	return nil
}

func (j *JwtToken) getClaim(k string) string {
	if claims, ok := j.parsed.Claims.(jwt.MapClaims); ok {
		return fmt.Sprint(claims[k])
	}
	return ""
}

func (j *JwtToken) get() []byte {
	return []byte(j.data)
}

func (j *JwtToken) UserId() string {
	return j.getClaim("user_id")
}

type Token interface {
	parse(rawToken string) error
	get() []byte
	UserId() string
}

func NewToken(data string) *Token {
	var token Token

	token = new(JwtToken)
	err := (token).parse(data)

	if err != nil {
		token = new(BasicToken)
		// Can't fail
		err = (token).parse(data)
	}

	return &token
}

type Info struct {
	Token    *Token
	CertData []byte
}

type key int

var infoKey key

func NewContext(ctx context.Context, info *Info) context.Context {
	return context.WithValue(ctx, infoKey, info)
}

func InfoFromContext(ctx context.Context) (Info, bool) {
	info, ok := ctx.Value(infoKey).(*Info)
	if info == nil {
		return Info{}, ok
	}

	return *info, ok
}

func (i Info) Scheme() string {
	if i.Token != nil {
		return BearerScheme
	}

	if len(i.CertData) > 0 {
		return CertScheme
	}

	return UnknownScheme
}

func (i Info) Hash() string {
	var key []byte
	if i.Token != nil {
		key = (*i.Token).get()
	}
	key = append(key, i.CertData...)
	hasher := sha256.New()
	return hex.EncodeToString(hasher.Sum(key))
}
