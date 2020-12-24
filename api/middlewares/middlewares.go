package middlewares

import (
	"net/http"
)

/* 응답을 JSON 포맷으로 세팅 */
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

/* 인증 토큰 유효성 검사 */
// func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		err := auth.TokenValid(r)
// 		if err != nil {
// 			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
// 			return
// 		}
// 		next(w, r)
// 	}
// }
