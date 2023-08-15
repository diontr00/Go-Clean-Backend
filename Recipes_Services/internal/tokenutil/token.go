package tokenutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"khanhanhtr/sample/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(user string, secret string, expiry time.Time) (model.JWTOutput, error) {
	claims := &model.JWTCustomClaims{
		Username: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return model.JWTOutput{}, err
	}
	return model.JWTOutput{
		Token:   t,
		Expires: expiry,
	}, nil
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected Signing method : %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractClaim(requestToken string, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(
		requestToken,
		&model.JWTCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected Signing method : %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

	if err != nil {
		return "", fmt.Errorf("Error parsing claims : %v", err)
	}

	claims, ok := token.Claims.(*model.JWTCustomClaims)
	if !ok && !token.Valid {
		return "", fmt.Errorf("Invalid Token : %v", token)
	}
	return claims.Username, nil
}

func GenerateSignature(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
