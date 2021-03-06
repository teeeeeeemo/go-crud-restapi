package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

/* 토큰 생성 */
func CreateToken(user_id uint32) (string, error) {

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

/* 토큰 유효성 검증 */
func TokenValid(r *http.Request) error {

	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return nil

}

/* 토큰 추출 */
func ExtractToken(r *http.Request) string {

	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	bearerTokenSplited := strings.Split(bearerToken, " ")
	if len(bearerTokenSplited) == 2 {
		bearerUpper := strings.ToUpper("bearer") // TODO 상수화 필요
		if strings.Compare(bearerUpper, strings.ToUpper(bearerTokenSplited[0])) == 0 {
			return bearerTokenSplited[1]
		}
	}
	return ""

}

/* 토큰 아이디 추출 */
func ExtractTokenID(r *http.Request) (uint32, error) {

	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil

}

/* 터메널에서 클레임을 예쁘게 표시 */
func Pretty(data interface{}) {

	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))

}
