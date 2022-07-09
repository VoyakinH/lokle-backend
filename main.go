package main

import (
	"context"

	// "database/sql"
	"encoding/json"

	// "github.com/xuri/excelize/v2"
	"net/http"
	"os"

	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var ctx = context.Background()

const expCookieTime = 1382400

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

// var psqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 	"185.225.35.60", 5432, "postgres", "1991055q", "glasha")

type Auth struct {
	Email    string
	Password string
}

func createUserSession(w http.ResponseWriter, r *http.Request) {
	var a Auth
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionID, _ := uuid.NewRandom()

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  sessionID.String(),
		MaxAge: expCookieTime,
	}

	_, err = rdb.SetNX(ctx, sessionID.String(), 0, expCookieTime*time.Second).Result()
	if err != nil {
		w.WriteHeader(523)
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			return
		}
		return
	}

	http.SetCookie(w, cookie)
}

func deleteUserSession(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("session-id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rdb.Del(ctx, tokenCookie.Value).Val()

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

func checkUserSession(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("session-id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	val := rdb.Exists(ctx, tokenCookie.Value).Val()
	if val < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rdb.Expire(ctx, tokenCookie.Value, expCookieTime*time.Second)

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  tokenCookie.Value,
		MaxAge: expCookieTime,
	}

	http.SetCookie(w, cookie)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/login", createUserSession)
	mux.HandleFunc("/api/v1/logout", deleteUserSession)
	mux.HandleFunc("/api/v1/validate_user", checkUserSession)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		return
	}
}
