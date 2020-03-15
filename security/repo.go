package security

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type Repo struct {
	findUserLoginByUsername      *sql.Stmt
	findPermissionsByUserLoginId *sql.Stmt
	getUserLogin                 *sql.Stmt
}

func InitRepo(db *sql.DB) *Repo {
	findUserLoginByUsername, err := db.Prepare(
		`select id, username, password from user_login where username = $1`)
	if err != nil {
		log.Fatalln(err)
	}

	findPermissionsByUserLoginId, err := db.Prepare(
		`select p.name from security_permission p
         inner join security_group_permission gp on gp.security_permission_id = p.id
         inner join security_group g on gp.security_group_id = g.id
         inner join user_login_security_group u on u.security_group_id = g.id
         where u.user_login_id = $1`)
	if err != nil {
		log.Fatalln(err)
	}

	getUserLogin, err := db.Prepare(
		`select id, username, password from user_login where id = $1`)
	if err != nil {
		log.Fatalln(err)
	}

	return &Repo{
		findUserLoginByUsername:      findUserLoginByUsername,
		findPermissionsByUserLoginId: findPermissionsByUserLoginId,
		getUserLogin:                 getUserLogin,
	}
}

func (repo *Repo) FindUserLoginByUsername(
	ctx context.Context, username string) *UserLogin {

	defer log.Printf("FindUserLoginByUsername %s\n", username)

	row := repo.findUserLoginByUsername.QueryRowContext(ctx, username)

	user := &UserLogin{}
	err := row.Scan(&user.Id, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Fatalln(err)
	}

	return user
}

func (repo *Repo) FindPermissionsByUserLoginId(
	ctx context.Context, id uuid.UUID) []string {

	defer log.Printf("FindPermissionsByUserLoginId %s\n", id)

	result := make([]string, 0)
	rows, err := repo.findPermissionsByUserLoginId.QueryContext(ctx, id)
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, name)
	}

	return result
}

func (repo *Repo) GetUserLogin(ctx context.Context, id uuid.UUID) *UserLogin {
	defer log.Printf("GetUserLogin %s\n", id)

	row := repo.getUserLogin.QueryRowContext(ctx, id)
	user := &UserLogin{}
	err := row.Scan(&user.Id, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Fatalln(err)
	}
	return user
}
