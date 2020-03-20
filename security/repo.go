package security

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db                           *sqlx.DB
	findUserLoginByUsername      *sqlx.Stmt
	findPermissionsByUserLoginId *sqlx.Stmt
	getUserLogin                 *sqlx.Stmt
	getClientUserLogin           *sqlx.Stmt
	getAllGroup                  *sqlx.Stmt
	getAllPermission             *sqlx.Stmt
	getAllGroupPermission        *sqlx.Stmt
}

func InitRepo(db *sqlx.DB) *Repo {
	query := `select id, username, password from user_login where username = ?`
	findUserLoginByUsername, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panicln(err)
	}

	query = `select p.name from security_permission p
         inner join security_group_permission gp on gp.security_permission_id = p.id
         inner join security_group g on gp.security_group_id = g.id
         inner join user_login_security_group u on u.security_group_id = g.id
         where u.user_login_id = ?`
	findPermissionsByUserLoginId, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panicln(err)
	}

	query = `select id, username, password from user_login where id = ?`
	getUserLogin, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panicln(err)
	}

	query = `select id, username, created_at, updated_at
        from user_login where id = ?`
	getClientUserLogin, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panicln(err)
	}

	query = `select id, name, created_at from security_group`
	getAllGroup, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panicln(err)
	}

	query = `select * from security_permission`
	getAllPermission, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panicln(err)
	}

	query = `select security_group_id,
            security_permission_id, created_at
            from security_group_permission`
	getAllGroupPermission, err := db.Preparex(db.Rebind(query))
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
	ctx context.Context, username string) (UserLogin, error) {

	log.Println("FindUserLoginByUsername", username)

	user := UserLogin{}
	return user, repo.findUserLoginByUsername.GetContext(ctx, &user, username)
}

func (repo *Repo) FindPermissionsByUserLoginId(
	ctx context.Context, id uuid.UUID) ([]string, error) {

	log.Println("FindPermissionsByUserLoginId", id)

	result := make([]string, 0)
	err := repo.findPermissionsByUserLoginId.SelectContext(ctx, &result, id)
	return result, err
}

func (repo *Repo) GetUserLogin(ctx context.Context, id uuid.UUID) (UserLogin, error) {
	log.Println("GetUserLogin", id)

	user := UserLogin{}
	return user, repo.getUserLogin.GetContext(ctx, &user, id)
}

func (repo *Repo) GetClientUserLogin(
	ctx context.Context, id uuid.UUID) (ClientUserLogin, error) {

	log.Println("GetClientUserLogin", id)

	user := ClientUserLogin{}
	err := repo.getClientUserLogin.GetContext(ctx, &user, id)
	return user, err
}

func (repo *Repo) GetAllGroup(ctx context.Context) ([]Group, error) {
	log.Println("GetAllGroup")

	result := make([]Group, 0)
	return result, repo.getAllGroup.SelectContext(ctx, &result)
}

func (repo *Repo) GetAllPermission(ctx context.Context) ([]Permission, error) {
	log.Println("GetAllPermission")

	result := make([]Permission, 0)
	return result, repo.getAllPermission.SelectContext(ctx, &result)
}

func (repo *Repo) GetAllGroupPermission(ctx context.Context) ([]GroupPermission, error) {
	log.Println("GetAllGroupPermission")

	result := make([]GroupPermission, 0)

	rows, err := repo.getAllGroupPermission.QueryContext(ctx)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		gp := GroupPermission{}
		err = rows.Scan(&gp.Id.GroupId, &gp.Id.PermissionId, &gp.CreatedAt)
		if err != nil {
			return result, err
		}
		result = append(result, gp)
	}

	return result, rows.Err()
}

func (repo *Repo) InsertGroupPermission(
	ctx context.Context, groupId, permissionId int16) error {

	log.Println("InsertGroupPermission", groupId, permissionId)

	query := `insert into security_group_permission
        (security_group_id, security_permission_id)
        values (?, ?)`
	_, err := repo.db.ExecContext(ctx,
		repo.db.Rebind(query), groupId, permissionId)

	return err
}

func (repo *Repo) DeleteGroupPermission(
	ctx context.Context, groupId, permissionId int16) error {

	log.Println("DeleteGroupPermission", groupId, permissionId)

	query := `delete from security_group_permission
        where security_group_id = ?
        and security_permission_id = ?`
	_, err := repo.db.ExecContext(ctx,
		repo.db.Rebind(query), groupId, permissionId)

	return err
}

func (repo *Repo) InsertGroup(ctx context.Context, name string) (int16, error) {
	log.Println("InsertGroup", name)

	var maxId int16
	err := repo.db.GetContext(ctx, &maxId,
		`select max(id) from security_group`)
	if err != nil {
		return 0, err
	}

	query := `insert into security_group(id, name) values (?, ?)`
	_, err = repo.db.ExecContext(ctx,
		repo.db.Rebind(query), maxId+1, name)

	return maxId + 1, err
}

func (repo *Repo) GetGroup(ctx context.Context, id int16) (Group, error) {
	log.Println("GetGroup", id)

	query := `select id, name, created_at
        from security_group where id = ?`

	group := Group{}
	return group, repo.db.GetContext(ctx, &group, repo.db.Rebind(query), id)
}
