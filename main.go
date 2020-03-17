package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"baseweb/account"
	"baseweb/security"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Root struct {
	db           *sql.DB
	redisClient  *redis.Client
	securityRepo *security.Repo
	security     *security.Root
	accountRepo  *account.Repo
	account      *account.Root
}

func (root *Root) homeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := root.securityRepo.FindUserLoginByUsername(r.Context(), "admin")
	if err == sql.ErrNoRows {
		log.Println("[ERROR]", "user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err == context.Canceled || err == context.DeadlineExceeded {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Panicln(err)
	}

	log.Println(user)

	bytes, err := json.Marshal(user)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Fprintln(w, string(bytes))
}

func (root *Root) Authenticated(handler security.Handler) security.Handler {
	return security.Authenticated(root.redisClient, root.securityRepo, handler)
}

func main() {
	config := "user=postgres password=1 dbname=baseweb sslmode=disable"
	db, err := sql.Open("postgres", config)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	securityRepo := security.InitRepo(db)
	accountRepo := account.InitRepo(db)

	root := &Root{
		db:           db,
		redisClient:  redisClient,
		securityRepo: securityRepo,
		security:     security.InitRoot(securityRepo),
		accountRepo:  accountRepo,
		account:      account.InitRoot(accountRepo),
	}

	router := mux.NewRouter()

	router.HandleFunc("/", root.Authenticated(
		security.Authorized(
			"VIEW_EDIT_USER_LOGIN", root.homeHandler))).Methods("GET")

	router.HandleFunc("/api/login",
		root.Authenticated(root.security.LoginHandler)).Methods("POST")

	router.HandleFunc("/api/security/permission",
		root.Authenticated(
			security.Authorized("VIEW_EDIT_SECURITY_PERMISSION",
				root.security.SecurityPermissionHandler))).Methods("GET")

	router.HandleFunc("/api/security/save-group-permissions",
		root.Authenticated(
			security.Authorized("VIEW_EDIT_SECURITY_PERMISSION",
				root.security.SaveGroupPermissonsHandler))).Methods("POST")

	router.HandleFunc("/api/security/add-security-group",
		root.Authenticated(
			security.Authorized("VIEW_EDIT_SECURITY_GROUP",
				root.security.AddSecurityGroupHandler))).Methods("POST")

	router.HandleFunc("/api/account/add-party",
		root.Authenticated(
			security.Authorized("VIEW_EDIT_PARTY",
				root.account.AddPartyHandler))).Methods("POST")

	router.HandleFunc("/api/account/view-person",
		root.Authenticated(
			security.Authorized("VIEW_EDIT_PARTY",
				root.account.ViewPersonHandler))).Methods("GET")

	http.Handle("/", router)

	err = http.ListenAndServe(":8080", http.HandlerFunc(applyJson(http.DefaultServeMux)))
	log.Fatalln(err)
}

func applyJson(handler http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	}
}
