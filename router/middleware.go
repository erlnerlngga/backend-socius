package router

import (
	"log"
	"net/http"

	"strings"

	"github.com/erlnerlngga/backend-socius/util"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

// middleware to handle jwt verification
func WithJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		// get auth header
		authHeader := r.Header.Get("Authorization")

		// sanity check
		if authHeader == "" {
			util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "no auth header"})
			return
		}

		// split the header space
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "Invalid auth header"})
		}

		tokenString := headerParts[1]

		// init claims
		claims := new(util.ClaimsType)

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return util.JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "Signature Invalid"})
				return
			}

			if strings.HasPrefix(err.Error(), "token is expired by") {
				util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "expired token"})
				return
			}

			util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: err.Error()})
			return
		}

		if !token.Valid {
			util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "token invalid"})
		}

		next.ServeHTTP(w, r)
	})
}

func WithJWTAuthWS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		thoken := chi.URLParam(r, "token")

		// tokenString := thoken
		log.Println("tokenString", thoken)

		// init claims
		claims := new(util.ClaimsType)

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match

		token, err := jwt.ParseWithClaims(thoken, claims, func(t *jwt.Token) (interface{}, error) {
			return util.JwtKey, nil
		})

		log.Println("TOKEN", token)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.Println("1. WithJWTAuthWS ", err)
				util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "Signature Invalid"})
				return
			}

			if strings.HasPrefix(err.Error(), "token is expired by") {
				log.Println("2. WithJWTAuthWS ", err)
				util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "expired token"})
				return
			}

			log.Println("3. WithJWTAuthWS ", err)
			util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: err.Error()})
			return
		}

		if !token.Valid {
			log.Println("4. WithJWTAuthWS ", err)
			util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "token invalid"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
