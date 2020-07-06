package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"baseweb/account"
	"baseweb/basic"
	"baseweb/export"
	"baseweb/facility"
	importProduct "baseweb/import"
	"baseweb/order"
	"baseweb/product"
	"baseweb/salesman"
	"baseweb/salesroute"
	"baseweb/schedule"
	"baseweb/security"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Root struct {
	router         *mux.Router
	db             *sqlx.DB
	redisClient    *redis.Client
	securityRepo   *security.Repo
	security       *security.Root
	accountRepo    *account.Repo
	account        *account.Root
	productRepo    *product.Repo
	product        *product.Root
	facilityRepo   *facility.Repo
	facility       *facility.Root
	importRepo     *importProduct.Repo
	importProduct  *importProduct.Root
	orderRepo      *order.Repo
	order          *order.Root
	exportRepo     *export.Repo
	export         *export.Root
	salesrouteRepo *salesroute.Repo
	salesroute     *salesroute.Root
	salesman       *salesman.Root
	salesmanRepo   *salesman.Repo
	schedule       *schedule.Root
	scheduleRepo   *schedule.Repo
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
	psqlHost := os.Getenv("PSQL_HOST")
	if psqlHost == "" {
		psqlHost = "localhost"
	}

	psqlUser := os.Getenv("PSQL_USER")
	if psqlUser == "" {
		psqlUser = "postgres"
	}

	psqlPasswd := os.Getenv("PSQL_PASSWD")
	if psqlPasswd == "" {
		psqlPasswd = "1"
	}

	psqlSsl := os.Getenv("PSQL_SSL")
	if psqlSsl == "" {
		psqlSsl = "disable"
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	config := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=baseweb sslmode=%s",
		psqlHost, psqlUser, psqlPasswd, psqlSsl)
	db := sqlx.MustConnect("postgres", config)
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", redisHost),
		DB:   0,
	})

	securityRepo := security.InitRepo(db)
	accountRepo := account.InitRepo(db)
	productRepo := product.InitRepo(db)
	facilityRepo := facility.InitRepo(db)
	importRepo := importProduct.InitRepo(db)
	orderRepo := order.InitRepo(db)
	exportRepo := export.InitRepo(db)
	salesrouteRepo := salesroute.InitRepo(db)
	salesmanRepo := salesman.InitRepo(db)
	scheduleRepo := schedule.InitRepo(db)

	router := mux.NewRouter()

	root := &Root{
		router:         router,
		db:             db,
		redisClient:    redisClient,
		securityRepo:   securityRepo,
		security:       security.InitRoot(securityRepo),
		accountRepo:    accountRepo,
		account:        account.InitRoot(accountRepo),
		productRepo:    productRepo,
		product:        product.InitRoot(productRepo),
		facilityRepo:   facilityRepo,
		facility:       facility.InitRoot(facilityRepo),
		importRepo:     importRepo,
		importProduct:  importProduct.InitRoot(importRepo),
		orderRepo:      orderRepo,
		order:          order.InitRoot(orderRepo),
		exportRepo:     exportRepo,
		export:         export.InitRoot(exportRepo),
		salesrouteRepo: salesrouteRepo,
		salesroute:     salesroute.InitRoot(salesrouteRepo),
		salesmanRepo:   salesmanRepo,
		salesman:       salesman.InitRoot(salesmanRepo),
		scheduleRepo:   scheduleRepo,
		schedule:       schedule.InitRoot(scheduleRepo),
	}

	root.GetAuthorized("/", "VIEW_EDIT_USER_LOGIN", root.homeHandler)

	SecurityRoutes(root)
	AccountRoutes(root)
	ProductRoutes(root)
	FacilityRoutes(root)
	ImportRoutes(root)
	OrderRoutes(root)
	ExportRoutes(root)
	SalesrouteRoutes(root)
	SalesmanRoutes(root)
	ScheduleRoutes(root)

	http.Handle("/", router)

	log.Println("Server is running")
	err := http.ListenAndServe(":8080",
		http.HandlerFunc(applyJson(http.DefaultServeMux)))
	log.Fatal(err)
}

func applyJson(handler http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	}
}
