package security

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Root struct {
	Repo *Repo
}

func (root *Root) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		UserLogin   ClientUserLogin `json:"userLogin"`
		Permissions []string        `json:"securityPermissions"`
	}

	ctx := r.Context()

	userLogin := ctx.Value("userLogin").(UserLogin)
	permissions := ctx.Value("permissions").([]string)

	user, err := root.Repo.GetClientUserLogin(ctx, userLogin.Id)
	if err == sql.ErrNoRows {
		log.Panicln("User not found")
	}
	if err == context.Canceled || err == context.DeadlineExceeded {
		w.WriteHeader(500)
		return
	}
	if err != nil {
		log.Panicln(err)
	}

	response := Response{
		UserLogin:   user,
		Permissions: permissions,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panicln(err)
	}
}

func (root *Root) SecurityPermissionHandler(
	w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Groups           []Group           `json:"securityGroups"`
		Permissions      []Permission      `json:"securityPermissions"`
		GroupPermissions []GroupPermission `json:"securityGroupPermissions"`
	}

	ctx := r.Context()

	groups, err := root.Repo.GetAllGroup(ctx)
	if err == context.Canceled || err == context.DeadlineExceeded {
		w.WriteHeader(500)
		return
	}
	if err != nil {
		log.Panicln(err)
	}

	permissions, err := root.Repo.GetAllPermission(ctx)
	if err == context.Canceled || err == context.DeadlineExceeded {
		w.WriteHeader(500)
		return
	}
	if err != nil {
		log.Panicln(err)
	}

	groupPerms, err := root.Repo.GetAllGroupPermission(ctx)
	if err == context.Canceled || err == context.DeadlineExceeded {
		w.WriteHeader(500)
		return
	}
	if err != nil {
		log.Panicln(err)
	}

	response := Response{
		Groups:           groups,
		Permissions:      permissions,
		GroupPermissions: groupPerms,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panicln(err)
	}
}

func (root *Root) SaveGroupPermissonsHandler(
	w http.ResponseWriter, r *http.Request) {

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
		w.WriteHeader(500)
		return
	}

	for _, permId := range request.ToBeInserted {
		err := root.Repo.InsertGroupPermission(ctx, request.GroupId, permId)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}

	for _, permId := range request.ToBeDeleted {
		err := root.Repo.DeleteGroupPermission(ctx, request.GroupId, permId)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}

	groupPerms, err := root.Repo.GetAllGroupPermission(ctx)
	if err == context.Canceled || err == context.DeadlineExceeded {
		w.WriteHeader(500)
		return
	}
	if err != nil {
		log.Panicln(err)
	}

	response := Response{
		GroupPermissions: groupPerms,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panicln(err)
	}
}

func (root *Root) AddSecurityGroupHandler(
	w http.ResponseWriter, r *http.Request) {

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
		w.WriteHeader(500)
		return
	}

	id, err := root.Repo.InsertGroup(ctx, request.Name)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	group, err := root.Repo.GetGroup(ctx, id)
	if err == sql.ErrNoRows {
		log.Panicln("Group not found")
	}
	if err == context.Canceled || err == context.DeadlineExceeded {
		w.WriteHeader(500)
		return
	}
	if err != nil {
		log.Panicln(err)
	}

	response := Response{
		Group: group,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panicln(err)
	}
}
