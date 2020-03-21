package security

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"baseweb/basic"

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
		return "", fmt.Errorf("newSession: %w", err)
	}
	return token, err
}

var InvalidSession error = errors.New("Invalid Session")

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
	if err != nil {
		return id, fmt.Errorf("sessionValid: %w", err)
	}
	if ok != "ok" {
		return id, errors.New("sessionValid: value not ok")
	}

	err = redisClient.Expire(key, SESSION_EXPIRATION_TIME).Err()
	if err != nil {
		return id, err
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

func Authenticated(redisClient *redis.Client,
	repo *Repo, handler basic.Handler) basic.Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		redisClient = redisClient.WithContext(ctx)

		next := func(user UserLogin, permissions []string) error {
			ctx = context.WithValue(ctx, "userLogin", user)
			ctx = context.WithValue(ctx, "permissions", permissions)
			r = r.WithContext(ctx)
			return handler(w, r)
		}

		token := r.Header.Get("X-Auth-Token")

		id, err := sessionValid(redisClient, token)
		user, permissions, err := getUserLoginInfo(ctx, repo, id, err)

		if err == nil {
			return next(user, permissions)
		}
		if err != InvalidSession && err != sql.ErrNoRows {
			return fmt.Errorf("Authenticated: %w", err)
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}

		user, err = repo.FindUserLoginByUsername(ctx, username)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password), []byte(password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}

		permissions, err = repo.FindPermissionsByUserLoginId(ctx, user.Id)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}

		token, err = newSession(redisClient, user.Id)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}

		w.Header().Add("X-Auth-Token", token)
		return next(user, permissions)
	}
}

func Authorized(perm string, handler basic.Handler) basic.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
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
			return nil
		}

		return handler(w, r)
	}
}
