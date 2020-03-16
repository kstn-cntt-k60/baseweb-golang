package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const SESSION_EXPIRATION_TIME = 30 * time.Minute

func newSession(redisClient *redis.Client, id uuid.UUID) (string, error) {

	b := make([]byte, 20)
	rand.Read(b)

	token := fmt.Sprintf("%s:%s", id.String(),
		base64.StdEncoding.EncodeToString(b))

	key := fmt.Sprintf("session:%s", token)

	err := redisClient.Set(key, "ok", SESSION_EXPIRATION_TIME).Err()
	if err != nil {
		log.Panicln(err)
	}
	return token, err
}

func sessionValid(redisClient *redis.Client,
	token string) (uuid.UUID, bool) {

	id := uuid.New()
	if token == "" {
		return id, false
	}

	key := fmt.Sprintf("session:%s", token)
	ok, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		return id, false
	}
	if err != nil {
		log.Panicln(err)
	}
	if ok != "ok" {
		log.Panicln("value not ok")
	}

	tokens := strings.Split(token, ":")
	if len(tokens) == 0 {
		return id, false
	}

	id, err = uuid.Parse(tokens[0])
	if err != nil {
		return id, false
	}

	return id, true
}

func getUserLoginInfo(ctx context.Context,
	repo *Repo, id uuid.UUID, valid bool) (UserLogin, []string, bool) {

	user := UserLogin{}
	permissions := make([]string, 0)
	if !valid {
		return user, permissions, false
	}

	user, ok := repo.GetUserLogin(ctx, id)
	if !ok {
		return user, permissions, false
	}

	permissions = repo.FindPermissionsByUserLoginId(ctx, id)
	return user, permissions, true
}

type Handler func(http.ResponseWriter, *http.Request)

func Authenticated(redisClient *redis.Client, repo *Repo, handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		redisClient = redisClient.WithContext(ctx)

		next := func(user UserLogin, permissions []string) {
			ctx = context.WithValue(ctx, "userLogin", user)
			ctx = context.WithValue(ctx, "permissions", permissions)
			r = r.WithContext(ctx)
			handler(w, r)
		}

		token := r.Header.Get("X-Auth-Token")

		id, valid := sessionValid(redisClient, token)
		user, permissions, ok := getUserLoginInfo(ctx, repo, id, valid)

		if ok {
			next(user, permissions)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(401)
			return
		}

		user, ok = repo.FindUserLoginByUsername(ctx, username)
		if !ok {
			w.WriteHeader(401)
			return
		}

		err := bcrypt.CompareHashAndPassword(
			[]byte(user.Password), []byte(password))
		if err != nil {
			w.WriteHeader(401)
			return
		}

		permissions = repo.FindPermissionsByUserLoginId(ctx, user.Id)

		token, err = newSession(redisClient, user.Id)
		if err != nil {
			log.Panicln(err)
		}
		w.Header().Add("X-Auth-Token", token)
		next(user, permissions)
	}
}

func Authorized(perm string, handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		permissions := ctx.Value("permissions").([]string)

		authorized := false
		for _, e := range permissions {
			if perm == e {
				authorized = true
				break
			}
		}

		if !authorized {
			w.WriteHeader(403)
			return
		}

		handler(w, r)
	}
}
