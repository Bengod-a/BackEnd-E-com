package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type jwtServices struct {
	secretKey string
}

func JWTAuthService() *jwtServices {
	return &jwtServices{secretKey: GetSecretKey()}
}

func GetSecretKey() string {
	secret := os.Getenv("SECRET_KEY")
	if secret == "SECRET_KEY" {
		secret = ""
	}
	return secret
}

func (service *jwtServices) GenerateToken(id int, email string, role string, cooldown time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * cooldown).Unix(),
		"iat":   time.Now().Unix(),
	}

	fmt.Println("claims-->", claims)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(service.secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token: algorithm %v is not valid", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})
}
