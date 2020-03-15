package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"baseweb/security"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Root struct {
	db           *sql.DB
	redisClient  *redis.Client
	securityRepo *security.Repo
}

func (root *Root) homeHandler(w http.ResponseWriter, r *http.Request) {
	user := root.securityRepo.FindUserLoginByUsername(r.Context(), "admin")
	if user == nil {
		log.Println("[ERROR]", "user not found")
		w.WriteHeader(404)
		return
	}

	log.Println(user)

	bytes, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintln(w, string(bytes))
}

func main() {
	config := "user=postgres password=1 dbname=baseweb sslmode=disable"
	db, err := sql.Open("postgres", config)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	root := &Root{
		db:           db,
		redisClient:  redisClient,
		securityRepo: security.InitRepo(db),
	}

	router := mux.NewRouter()
	router.HandleFunc("/", security.Authenticated(root.redisClient, root.securityRepo, root.homeHandler))

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
