package utils

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenMetadata struct {
	Expires int64
	UserId  string
}

func ExtractTokenMetadata(r *http.Request) (*TokenMetadata, error) {
	token, err := verifyToken(r)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	var isValid bool = ok && token.Valid

	if isValid {
		expires := int64(claims["exp"].(float64))

		userId := claims["userId"].(string)

		return &TokenMetadata{
			Expires: expires,
			UserId:  userId,
		}, nil
	}

	return nil, err
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	var tokenString string = extractToken(r)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)

	if err != nil {
		return nil, err
	}

	return token, nil
}

// extract token from http header
func extractToken(r *http.Request) string {
	var header string = r.Header.Get("Authorization")
	token := strings.Split(header, " ")

	var isEmpty bool = header == "" || len(token) < 2

	if isEmpty {
		return ""
	}

	return token[1]
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(GetValue("JWT_SECRET_KEY")), nil
}

func CheckToken(r *http.Request) (*TokenMetadata, error) {
	var now int64 = time.Now().Unix()

	claims, err := ExtractTokenMetadata(r)

	if err != nil {
		return nil, err
	}

	var expires int64 = claims.Expires

	if now > expires {
		return nil, err
	}

	return claims, err
}

func GenerateAccessToken(userId string) (string, error) {
	secret := GetValue("JWT_SECRET_KEY")

	minutesCount, _ := strconv.Atoi(GetValue("JWT_SECRET_KEY_EXPIRE_MINUTES"))

	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()
	claims["userId"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	decodedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return decodedToken, nil
}
