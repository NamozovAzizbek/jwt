package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

var SecretKey = []byte("keyKeyKey")

func main() {

	http.HandleFunc("/home", verifyJWT(handlePage))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}
}

func handlePage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var message Message
	err := json.NewDecoder(request.Body).Decode(&message)
	if err != nil {
		return
	}
	err = json.NewEncoder(writer).Encode(message)
	if err != nil {
		return
	}
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["user"] = "username"
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	fmt.Println(tokenString)
	return tokenString, nil
}
func verifyJWT(endpointHandler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["Token"] != nil {
			token, err := jwt.Parse(request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodECDSA)
				if !ok {
					writer.WriteHeader(http.StatusUnauthorized)
					_, err := writer.Write([]byte("You're Unauthorized"))
					if err != nil {
						return nil, err
					}
				}
				return "", nil
			})
			if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err2 := writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
				if err2 != nil {
					return
				}
			}
			if token.Valid {
				endpointHandler(writer, request)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}

	})
}

func authPage(w http.ResponseWriter) {
	token, err := generateJWT()
	if err != nil {
		return
	}
	client := &http.Client{}
	request, _ := http.NewRequest("POST", "<http://localhost:8080/>", nil)
	request.Header.Set("Token", token)
	_, _ = client.Do(request)
}

func extractClims(_ http.ResponseWriter, r *http.Request) (string, error) {
	if r.Header["Token"] != nil {
		//tokenString := r.Header["Token"][0]
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {
				return nil, fmt.Errorf("there's an error with the signing method")
			}
			return SecretKey, nil
		})
		if err != nil {
			return "Error Parsing token", err
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			username := claims["username"].(string)
			return username, nil
		}

	}
	return "unable to extract claims", nil
}
