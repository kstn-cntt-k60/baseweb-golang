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
		log.Fatalln(err)
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
		log.Fatalln(err)
	}
	if ok != "ok" {
		log.Fatalln("value not ok")
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
	repo *Repo, id uuid.UUID, valid bool) (*UserLogin, []string) {

	permissions := make([]string, 0)
	if !valid {
		return nil, permissions
	}

	user := repo.GetUserLogin(ctx, id)
	if user == nil {
		return nil, permissions
	}

	permissions = repo.FindPermissionsByUserLoginId(ctx, id)
	return user, permissions
}

type Handler func(http.ResponseWriter, *http.Request)

func Authenticated(redisClient *redis.Client, repo *Repo, handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		redisClient = redisClient.WithContext(ctx)

		token := r.Header.Get("X-Auth-Token")

		id, valid := sessionValid(redisClient, token)
		user, permissions := getUserLoginInfo(ctx, repo, id, valid)

		next := func(user *UserLogin, permissions []string) {
			ctx = context.WithValue(ctx, "userLogin", user)
			ctx = context.WithValue(ctx, "permissions", permissions)
			r = r.WithContext(ctx)
			handler(w, r)
		}

		if user != nil {
			next(user, permissions)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(401)
			return
		}

		user = repo.FindUserLoginByUsername(ctx, username)
		if user == nil {
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
			log.Fatalln(err)
		}
		w.Header().Add("X-Auth-Token", token)
		next(user, permissions)
	}
}
