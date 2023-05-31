package router

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/erlnerlngga/backend-socius/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type ClaimsType struct {
	User_ID string `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// create jwt
func CreateJWT(user_id string) (string, error) {
	// declare expiration time with 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	// declare jwt claims
	claims := &ClaimsType{
		User_ID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// declare token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// create token jwt string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

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
		claims := new(ClaimsType)

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
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
