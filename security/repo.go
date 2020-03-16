package security

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type Repo struct {
	db                           *sql.DB
	findUserLoginByUsername      *sql.Stmt
	findPermissionsByUserLoginId *sql.Stmt
	getUserLogin                 *sql.Stmt
	getClientUserLogin           *sql.Stmt
	getAllGroup                  *sql.Stmt
	getAllPermission             *sql.Stmt
	getAllGroupPermission        *sql.Stmt
}

func InitRepo(db *sql.DB) *Repo {
	findUserLoginByUsername, err := db.Prepare(
		`select id, username, password from user_login where username = $1`)
	if err != nil {
		log.Panicln(err)
	}

	findPermissionsByUserLoginId, err := db.Prepare(
		`select p.name from security_permission p
         inner join security_group_permission gp on gp.security_permission_id = p.id
         inner join security_group g on gp.security_group_id = g.id
         inner join user_login_security_group u on u.security_group_id = g.id
         where u.user_login_id = $1`)
	if err != nil {
		log.Panicln(err)
	}

	getUserLogin, err := db.Prepare(
		`select id, username, password from user_login where id = $1`)
	if err != nil {
		log.Panicln(err)
	}

	getClientUserLogin, err := db.Prepare(
		`select id, username, created_at, updated_at
        from user_login where id = $1`)
	if err != nil {
		log.Panicln(err)
	}

	getAllGroup, err := db.Prepare(
		`select id, name, created_at from security_group`)
	if err != nil {
		log.Panicln(err)
	}

	getAllPermission, err := db.Prepare(
		`select id, name, created_at from security_permission`)
	if err != nil {
		log.Panicln(err)
	}

	getAllGroupPermission, err := db.Prepare(
		`select security_group_id, security_permission_id, created_at
        from security_group_permission`)
	if err != nil {
		log.Panicln(err)
	}

	return &Repo{
		db:                           db,
		findUserLoginByUsername:      findUserLoginByUsername,
		findPermissionsByUserLoginId: findPermissionsByUserLoginId,
		getUserLogin:                 getUserLogin,
		getClientUserLogin:           getClientUserLogin,
		getAllGroup:                  getAllGroup,
		getAllPermission:             getAllPermission,
		getAllGroupPermission:        getAllGroupPermission,
	}
}

func (repo *Repo) FindUserLoginByUsername(
	ctx context.Context, username string) (UserLogin, bool) {

	defer log.Printf("FindUserLoginByUsername %s\n", username)

	row := repo.findUserLoginByUsername.QueryRowContext(ctx, username)

	user := UserLogin{}
	err := row.Scan(&user.Id, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return user, false
	}
	if err != nil {
		log.Panicln(err)
	}

	return user, true
}

func (repo *Repo) FindPermissionsByUserLoginId(
	ctx context.Context, id uuid.UUID) []string {

	defer log.Printf("FindPermissionsByUserLoginId %s\n", id)

	result := make([]string, 0)
	rows, err := repo.findPermissionsByUserLoginId.QueryContext(ctx, id)
	if err != nil {
		log.Panicln(err)
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Panicln(err)
		}
		result = append(result, name)
	}

	return result
}

func (repo *Repo) GetUserLogin(ctx context.Context, id uuid.UUID) (UserLogin, bool) {
	defer log.Printf("GetUserLogin %s\n", id)

	row := repo.getUserLogin.QueryRowContext(ctx, id)
	user := UserLogin{}
	err := row.Scan(&user.Id, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return user, false
	}
	if err != nil {
		log.Panicln(err)
	}
	return user, true
}

func (repo *Repo) GetClientUserLogin(
	ctx context.Context, id uuid.UUID) (ClientUserLogin, bool) {

	defer log.Printf("GetClientUserLogin %s\n", id)

	row := repo.getClientUserLogin.QueryRowContext(ctx, id)
	user := ClientUserLogin{}
	err := row.Scan(&user.Id, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return user, false
	}
	if err != nil {
		log.Panicln(err)
	}

	return user, true
}

func (repo *Repo) GetAllGroup(ctx context.Context) (result []Group) {
	defer log.Println("GetAllGroup")

	rows, err := repo.getAllGroup.QueryContext(ctx)
	if err != nil {
		log.Panicln(err)
	}

	for rows.Next() {
		group := Group{}
		err := rows.Scan(&group.Id, &group.Name, &group.CreatedAt)
		if err != nil {
			log.Panicln(err)
		}
		result = append(result, group)
	}
	return
}

func (repo *Repo) GetAllPermission(ctx context.Context) (result []Permission) {
	defer log.Println("GetAllPermission")

	rows, err := repo.getAllPermission.QueryContext(ctx)
	if err != nil {
		log.Panicln(err)
	}

	for rows.Next() {
		perm := Permission{}
		err := rows.Scan(&perm.Id, &perm.Name, &perm.CreatedAt)
		if err != nil {
			log.Panicln(err)
		}
		result = append(result, perm)
	}
	return
}

func (repo *Repo) GetAllGroupPermission(ctx context.Context) (result []GroupPermission) {
	defer log.Println("GetAllGroupPermission")

	rows, err := repo.getAllGroupPermission.QueryContext(ctx)
	if err != nil {
		log.Panicln(err)
	}

	for rows.Next() {
		gp := GroupPermission{}
		err := rows.Scan(&gp.Id.GroupId, &gp.Id.PermissionId, &gp.CreatedAt)
		if err != nil {
			log.Panicln(err)
		}
		result = append(result, gp)
	}
	return
}

func (repo *Repo) InsertGroupPermission(
	ctx context.Context, groupId, permissionId int16) error {

	defer log.Println("InsertGroupPermission", groupId, permissionId)

	_, err := repo.db.ExecContext(ctx,
		`insert into security_group_permission
        (security_group_id, security_permission_id)
        values ($1, $2)`, groupId, permissionId)

	return err
}

func (repo *Repo) DeleteGroupPermission(
	ctx context.Context, groupId, permissionId int16) error {

	defer log.Println("DeleteGroupPermission", groupId, permissionId)

	_, err := repo.db.ExecContext(ctx,
		`delete from security_group_permission
        where security_group_id = $1
            and security_permission_id = $2`, groupId, permissionId)
	return err
}

func (repo *Repo) InsertGroup(ctx context.Context, name string) (int16, error) {
	defer log.Println("InsertGroup", name)

	row := repo.db.QueryRowContext(ctx,
		`select max(id) from security_group`)

	var maxId int16
	err := row.Scan(&maxId)
	if err != nil {
		return 0, err
	}

	_, err = repo.db.ExecContext(ctx,
		`insert into security_group(id, name)
        values ($1, $2)`, maxId+1, name)

	return maxId + 1, err
}

func (repo *Repo) GetGroup(ctx context.Context, id int16) (Group, bool) {
	defer log.Println("GetGroup", id)

	row := repo.db.QueryRowContext(ctx,
		`select id, name, created_at
        from security_group where id = $1`, id)

	group := Group{}
	err := row.Scan(&group.Id, &group.Name, &group.CreatedAt)
	if err == sql.ErrNoRows {
		return group, false
	}
	if err != nil {
		log.Panicln(err)
	}

	return group, true
}
