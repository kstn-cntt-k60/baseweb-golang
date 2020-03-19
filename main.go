package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"baseweb/account"
	"baseweb/basic"
	"baseweb/security"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Root struct {
	router       *mux.Router
	db           *sql.DB
	redisClient  *redis.Client
	securityRepo *security.Repo
	security     *security.Root
	accountRepo  *account.Repo
	account      *account.Root
}

func UnwrapHandler(h basic.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.String())
		err := h(w, r)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("[CANCELED]", r.Method, r.URL.String(), err)
			} else if errors.Is(err, context.DeadlineExceeded) {
				log.Println("[TIMEOUT]", r.Method, r.URL.String(), err)
			} else {
				log.Println("[ERROR]", r.Method, r.URL.String(), err)
			}
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (root *Root) Authenticated(handler basic.Handler) basic.Handler {
	return security.Authenticated(root.redisClient, root.securityRepo, handler)
}

func (root *Root) GetAuthenticated(url string, handler basic.Handler) {
	root.router.HandleFunc(url,
		UnwrapHandler(
			root.Authenticated(handler))).Methods("GET")
}

func (root *Root) PostAuthenticated(url string, handler basic.Handler) {
	root.router.HandleFunc(url,
		UnwrapHandler(
			root.Authenticated(handler))).Methods("POST")
}

func (root *Root) GetAuthorized(url string,
	perm string, handler basic.Handler) {
	root.router.HandleFunc(url,
		UnwrapHandler(
			root.Authenticated(
				security.Authorized(perm, handler)))).Methods("GET")
}

func (root *Root) PostAuthorized(url string, perm string,
	handler basic.Handler) {
	root.router.HandleFunc(url,
		UnwrapHandler(
			root.Authenticated(
				security.Authorized(perm, handler)))).Methods("POST")
}

func (root *Root) homeHandler(w http.ResponseWriter, r *http.Request) error {
	time.Sleep(5 * time.Second)

	user, err := root.securityRepo.FindUserLoginByUsername(r.Context(), "admin")
	if err != nil {
		return err
	}

	log.Println(user)

	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, string(bytes))
	return nil
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

	router := mux.NewRouter()

	root := &Root{
		router:       router,
		db:           db,
		redisClient:  redisClient,
		securityRepo: securityRepo,
		security:     security.InitRoot(securityRepo),
		accountRepo:  accountRepo,
		account:      account.InitRoot(accountRepo),
	}

	root.GetAuthorized("/", "VIEW_EDIT_USER_LOGIN", root.homeHandler)
	root.PostAuthenticated("/api/login", root.security.LoginHandler)

	root.GetAuthorized(
		"/api/security/permission",
		"VIEW_EDIT_SECURITY_PERMISSION",
		root.security.SecurityPermissionHandler)

	root.PostAuthorized(
		"/api/security/save-group-permissions",
		"VIEW_EDIT_SECURITY_PERMISSION",
		root.security.SaveGroupPermissonsHandler)

	root.PostAuthorized(
		"/api/security/add-security-group",
		"VIEW_EDIT_SECURITY_GROUP",
		root.security.AddSecurityGroupHandler)

	root.PostAuthorized(
		"/api/account/add-party",
		"VIEW_EDIT_PARTY",
		root.account.AddPartyHandler)

	root.GetAuthorized(
		"/api/account/view-person",
		"VIEW_EDIT_PARTY",
		root.account.ViewPersonHandler)

	root.GetAuthorized(
		"/api/account/view-customer",
		"VIEW_EDIT_PARTY",
		root.account.ViewCustomerHandler)

	root.PostAuthorized(
		"/api/account/update-person",
		"VIEW_EDIT_PARTY",
		root.account.UpdatePersonHandler)

	root.PostAuthorized(
		"/api/account/delete-person",
		"VIEW_EDIT_PARTY",
		root.account.DeletePersonHandler)

	root.PostAuthorized(
		"/api/account/update-customer",
		"VIEW_EDIT_PARTY",
		root.account.UpdateCustomerHandler)

	root.PostAuthorized(
		"/api/account/delete-customer",
		"VIEW_EDIT_PARTY",
		root.account.DeleteCustomerHandler)

	http.Handle("/", router)

	err = http.ListenAndServe(":8080",
		http.HandlerFunc(applyJson(http.DefaultServeMux)))
	log.Fatalln(err)
}

func applyJson(handler http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	}
}
