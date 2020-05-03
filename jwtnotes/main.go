package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

var userMap map[string]string
var SecretKey string

type MyCustomClaims struct {
	// This will hold a users username after authenticating.
	// Ignore `json:"username"` it's required by JSON
	Username string `json:"username"`

	// This will hold claims that are recommended having (Expiration, issuer)
	jwt.StandardClaims
}

func Validate(protectedPage http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("Auth")
		if err != nil {
			fmt.Println("req.Cookie:", err)
			http.Redirect(res, req, "/", http.StatusForbidden)
			return
		}

		// Cookies concatenate the key/value. Remove the Auth= part
		splitCookie := strings.Split(cookie.String(), "Auth=")

		token, err := jwt.ParseWithClaims(splitCookie[1], &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Prevents a known exploit
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
			}
			return []byte(SecretKey), nil
		})

		if nil != err {
			fmt.Println("ParseWithClaims err:", err)
			http.Redirect(res, req, "/", http.StatusForbidden)
			return
		}

		// Validate the token and save the token's claims to a context
		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			context.Set(req, "Claims", claims)
		} else {
			http.Redirect(res, req, "/", http.StatusForbidden)
			return
		}

		res.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8000")
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		protectedPage(res, req)
	})
}

func setToken(res http.ResponseWriter, req *http.Request, userName string) {

	// Expires the token and cookie in 1 hours
	expireCookie := time.Now().Add(time.Hour * 1)
	expireToken := expireCookie.Unix()

	// We'll manually assign the claims but in production you'd insert values from a database
	claims := MyCustomClaims{
		userName,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "example.com",
		},
	}

	// Create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with a secret.
	signedToken, _ := token.SignedString([]byte(SecretKey))

	// This cookie will store the token on the client side
	cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}

	http.SetCookie(res, &cookie)
	fmt.Println("set cookie:", cookie.Name, ":", cookie.Value)

}

func login(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer func() {
		r.Body.Close()
	}()

	body := map[string]string{}
	err = json.Unmarshal(buf, &body)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}
	username := body["username"]
	password := body["userpwd"]

	if pwd, found := userMap[username]; found && pwd == password {
		setToken(w, r, username)
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Write([]byte("login"))
		fmt.Println("login")
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

}

func logout(w http.ResponseWriter, r *http.Request) {
	deleteCookie := http.Cookie{Name: "Auth", Value: "", Expires: time.Now()}
	http.SetCookie(w, &deleteCookie)
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Write([]byte("logout"))
	fmt.Println("logout")
}

func api(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "Claims").(*MyCustomClaims)
	w.Write([]byte(claims.Username + " at api"))
	fmt.Println(claims.Username + " at api")
	context.Clear(r)
}

func main() {

	userMap = map[string]string{}
	userMap["user1"] = "111111"
	userMap["user2"] = "222222"

	SecretKey = "secret"

	router := mux.NewRouter()

	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("POST")
	router.HandleFunc("/api", Validate(api)).Methods("POST")

	log.Fatal(http.ListenAndServe(":9876", router))
}
