package security

import (
	"encoding/json"
	"net/http"
)

type Root struct {
	repo *Repo
}

func InitRoot(repo *Repo) *Root {
	return &Root{
		repo: repo,
	}
}

func (root *Root) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	type Response struct {
		UserLogin   ClientUserLogin `json:"userLogin"`
		Permissions []string        `json:"securityPermissions"`
	}

	ctx := r.Context()

	userLogin := ctx.Value("userLogin").(UserLogin)
	permissions := ctx.Value("permissions").([]string)

	user, err := root.repo.GetClientUserLogin(ctx, userLogin.Id)
	if err != nil {
		return err
	}

	response := Response{
		UserLogin:   user,
		Permissions: permissions,
	}

	return json.NewEncoder(w).Encode(response)
}

func (root *Root) SecurityPermissionHandler(
	w http.ResponseWriter, r *http.Request) error {

	type Response struct {
		Groups           []Group           `json:"securityGroups"`
		Permissions      []Permission      `json:"securityPermissions"`
		GroupPermissions []GroupPermission `json:"securityGroupPermissions"`
	}

	ctx := r.Context()

	groups, err := root.repo.GetAllGroup(ctx)
	if err != nil {
		return err
	}

	permissions, err := root.repo.GetAllPermission(ctx)
	if err != nil {
		return err
	}

	groupPerms, err := root.repo.GetAllGroupPermission(ctx)
	if err != nil {
		return err
	}

	response := Response{
		Groups:           groups,
		Permissions:      permissions,
		GroupPermissions: groupPerms,
	}

	return json.NewEncoder(w).Encode(response)
}

func (root *Root) SaveGroupPermissonsHandler(
	w http.ResponseWriter, r *http.Request) error {

	type Request struct {
		GroupId      int16   `json:"securityGroupId"`
		ToBeInserted []int16 `json:"toBeInserted"`
		ToBeDeleted  []int16 `json:"toBeDeleted"`
	}

	type Response struct {
		GroupPermissions []GroupPermission `json:"securityGroupPermissions"`
	}

	ctx := r.Context()

	request := Request{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return err
	}

	for _, permId := range request.ToBeInserted {
		err := root.repo.InsertGroupPermission(ctx, request.GroupId, permId)
		if err != nil {
			return err
		}
	}

	for _, permId := range request.ToBeDeleted {
		err := root.repo.DeleteGroupPermission(ctx, request.GroupId, permId)
		if err != nil {
			return err
		}
	}

	groupPerms, err := root.repo.GetAllGroupPermission(ctx)
	if err != nil {
		return err
	}

	response := Response{
		GroupPermissions: groupPerms,
	}

	return json.NewEncoder(w).Encode(response)
}

func (root *Root) AddSecurityGroupHandler(
	w http.ResponseWriter, r *http.Request) error {

	type Request struct {
		Name string `json:"name"`
	}

	type Response struct {
		Group Group `json:"securityGroup"`
	}

	ctx := r.Context()

	request := Request{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return err
	}

	id, err := root.repo.InsertGroup(ctx, request.Name)
	if err != nil {
		return err
	}

	group, err := root.repo.GetGroup(ctx, id)
	if err != nil {
		return err
	}

	response := Response{
		Group: group,
	}

	return json.NewEncoder(w).Encode(response)
}
