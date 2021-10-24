package authentication

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

type RefreshTokenData struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenDetails struct {
	AccessToken string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	td.AccessToken, err = at.SignedString([]byte("jdnfksdmfksd"))
	if err != nil {
		return nil, err
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte("jdnfksdmfksd"))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error)  {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("jdnfksdmfksd"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) (string, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return "", err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return "", err
	}
	mapClaims := token.Claims.(jwt.MapClaims)
	t := mapClaims["user_id"].(string)
	return t, nil
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func RefreshToken(refresh string) (*TokenDetails, error) {

	token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("jdnfksdmfksd"), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {

		userId := claims["user_id"].(string)

		ts, createErr := CreateToken(userId)
		if  createErr != nil {
			return nil, err
		}
		return ts, nil

	} else {
		return nil, errors.New( "refresh expired")
	}
}