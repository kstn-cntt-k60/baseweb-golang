package security

import (
	"context"
	"crypto/rand"
	"database/sql"
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

type InvalidSessionError struct {
}

func (*InvalidSessionError) Error() string {
	return "invalid session"
}

var InvalidSession = &InvalidSessionError{}

func sessionValid(redisClient *redis.Client,
	token string) (uuid.UUID, error) {

	id := uuid.New()
	if token == "" {
		return id, InvalidSession
	}

	key := fmt.Sprintf("session:%s", token)

	ok, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		return id, InvalidSession
	}
	if err == context.Canceled || err == context.DeadlineExceeded {
		return id, err
	}
	if err != nil {
		log.Panicln(err)
	}
	if ok != "ok" {
		log.Panicln("value not ok")
	}

	tokens := strings.Split(token, ":")
	if len(tokens) == 0 {
		return id, InvalidSession
	}

	id, err = uuid.Parse(tokens[0])
	if err != nil {
		return id, InvalidSession
	}

	return id, nil
}

func getUserLoginInfo(ctx context.Context,
	repo *Repo, id uuid.UUID, err error) (UserLogin, []string, error) {

	user := UserLogin{}
	permissions := make([]string, 0)
	if err != nil {
		return user, permissions, err
	}

	user, err = repo.GetUserLogin(ctx, id)
	if err != nil {
		return user, permissions, err
	}

	permissions, err = repo.FindPermissionsByUserLoginId(ctx, id)
	if err != nil {
		return user, permissions, err
	}

	return user, permissions, nil
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

		id, err := sessionValid(redisClient, token)
		user, permissions, err := getUserLoginInfo(ctx, repo, id, err)

		if err == nil {
			next(user, permissions)
			return
		}
		if err == context.Canceled || err == context.DeadlineExceeded {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err != InvalidSession && err != sql.ErrNoRows {
			log.Panicln(err)
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err = repo.FindUserLoginByUsername(ctx, username)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password), []byte(password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		permissions, err = repo.FindPermissionsByUserLoginId(ctx, user.Id)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err = newSession(redisClient, user.Id)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
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
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(w, r)
	}
}
