package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type AccountStmt struct {
	FindUserLoginByUsername *sql.Stmt
}

func InitStmt(db *sql.DB) *AccountStmt {
	accountStmt := new(AccountStmt)

	stmt, err := db.Prepare(
		`select id, username, created_at, updated_at 
        from user_login where username = $1`)
	if err != nil {
		log.Fatalln(err)
	}
	accountStmt.FindUserLoginByUsername = stmt

	return accountStmt
}

type Root struct {
	db          *sql.DB
	accountStmt *AccountStmt
}

type UserLogin struct {
	Id        uuid.UUID
	Username  string
	CreatedAt string
	UpdatedAt string
}

func findUserLoginByUsername(root *Root,
	ctx context.Context, username string) (*UserLogin, error) {

	result := new(UserLogin)

	row := root.accountStmt.FindUserLoginByUsername.QueryRowContext(ctx, username)
	err := row.Scan(&result.Id, &result.Username,
		&result.CreatedAt, &result.UpdatedAt)

	return result, err
}

func main() {
	config := "user=postgres password=1 dbname=baseweb sslmode=disable"
	db, err := sql.Open("postgres", config)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	root := new(Root)
	root.db = db
	root.accountStmt = InitStmt(db)

	router := mux.NewRouter()
	router.HandleFunc("/", root.homeHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func (root *Root) homeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := findUserLoginByUsername(root, r.Context(), "admin")
	if err != nil {
		log.Println("[ERROR]", err)
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

