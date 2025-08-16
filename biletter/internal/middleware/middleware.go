package middleware

import (
	"biletter/internal/adapters/db/postgresql"
	"biletter/internal/domain/entity"
	"biletter/pkg/logging"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jackc/pgconn"
	"net/http"
	"strings"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func (f AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)
	if err != nil {
		return
	}
}

type Middleware func(h AppHandler, s postgresql.PostgresStorage, l *logging.Logger, module, action string) AppHandler
type Chain []Middleware

func New(middlewares ...Middleware) Chain {
	var slice Chain
	return append(slice, middlewares...)
}

func (c Chain) Then(handler AppHandler, storage postgresql.PostgresStorage, logger *logging.Logger, module, action string) AppHandler {
	for i := range c {
		handler = c[len(c)-1-i](handler, storage, logger, module, action)
	}
	return handler
}

func ErrorMiddleware(h AppHandler, _ postgresql.PostgresStorage, _ *logging.Logger, _, _ string) AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var appErr *AppErr
		var validErr validation.Errors
		var pgConnErr *pgconn.PgError

		w.Header().Set("Content-Type", "application/json")

		err := h(w, r)

		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) || err.Error() == "not found" {
					w.WriteHeader(http.StatusNotFound)
					_, _ = w.Write(ErrNotFound.Marshal())

					return nil

				} else if errors.Is(err, BadRequest) {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write(BadRequest.Marshal())

					return nil

				} else if errors.Is(err, Unauthorized) {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = w.Write(Unauthorized.Marshal())

					return nil

				} else if errors.Is(err, RequestTimeout) {
					w.WriteHeader(http.StatusRequestTimeout)
					_, _ = w.Write(Unauthorized.Marshal())

					return nil

				} else if errors.Is(err, Forbidden) || err.Error() == "forbidden" {
					w.WriteHeader(http.StatusForbidden)
					_, _ = w.Write(Forbidden.Marshal())

					return nil

				} else if errors.Is(err, FailedDependency) {
					w.WriteHeader(http.StatusFailedDependency)
					_, _ = w.Write(FailedDependency.Marshal())

					return nil

				} else if appErr.Code == "DS-000002" {
					w.WriteHeader(http.StatusAccepted)
					_, _ = w.Write(appErr.Marshal())

					return nil

				} else if appErr.Code == "DS-500010" {
					w.WriteHeader(http.StatusNotExtended)
					_, _ = w.Write(appErr.Marshal())

					return nil
				}

				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(BadRequest.Marshal())

				return nil

			} else if errors.As(err, &validErr) {
				for _, e := range validErr {
					unwrapErr := errors.Unwrap(e)

					if unwrapErr != nil {
						Validation.Message += unwrapErr.Error()
						break

					} else {
						Validation.Message += e.Error() + "\n"
					}
				}

				Validation.Message = strings.TrimRight(Validation.Message, "\n")

				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(Validation.Marshal())

				Validation.Message = ""

				return nil

			} else if errors.As(err, &pgConnErr) {
				if pgConnErr.Code == "23503" {
					BadRequest.Message = "Выполнение невозможно"
				}

				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(BadRequest.Marshal())
				BadRequest.Message = "bad request"
				return nil

			} else if err.Error() == "forbidden" {
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write(Forbidden.Marshal())
				BadRequest.Message = "bad request"
				return nil
			}

			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write(systemError(err).Marshal())
		}
		return nil
	}
}

var authC = newAuthCache()

func BasicAuthMiddleware(h AppHandler, client postgresql.PostgresStorage, log *logging.Logger, _, _ string) AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		const realm = "biletter"

		// быстрая проверка и парс
		email, pass, ok := r.BasicAuth()
		if !ok || email == "" {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s", charset="UTF-8"`, realm))
			return Unauthorized
		}

		// ключ кэша: email \x00 sha256(pass) lower-hex
		sum := sha256.Sum256([]byte(pass))
		var hexPass [64]byte
		hex.Encode(hexPass[:], sum[:])
		cacheKey := email + "\x00" + string(hexPass[:])

		if e, ok := authC.get(cacheKey); ok {
			if e.neg {
				w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s", charset="UTF-8"`, realm))
				return Unauthorized
			}
			// кладём user в контекст и дальше
			r = r.WithContext(context.WithValue(r.Context(), entity.UserCtxKey, e.user))
			return h(w, r)
		}

		// блок быстрых атак: проверим счетчик фейлов (если используешь Redis)
		// if blocked(email, ip) { return Unauthorized }

		u, err := client.GetUserByEmail(r.Context(), email)
		if err != nil || !u.IsActive {
			authC.setBad(cacheKey)
			w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s", charset="UTF-8"`, realm))
			return Unauthorized
		}

		if !verifyPassword(u, pass) {
			authC.setBad(cacheKey)
			w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s", charset="UTF-8"`, realm))
			return Unauthorized
		}

		// success
		authC.setOK(cacheKey, u)
		r = r.WithContext(context.WithValue(r.Context(), entity.UserCtxKey, u))
		return h(w, r)
	}
}

func verifyPassword(u entity.AuthUser, given string) bool {
	if u.PasswordPlain != nil {
		return subtle.ConstantTimeCompare([]byte(*u.PasswordPlain), []byte(given)) == 1
	}
	// sha256 → hex без fmt и без lower
	sum := sha256.Sum256([]byte(given))
	var buf [64]byte
	hex.Encode(buf[:], sum[:]) // всегда lower-hex
	// ВАЖНО: убедись, что password_hash в БД приведён к lower-hex
	return subtle.ConstantTimeCompare([]byte(u.PasswordHash), buf[:]) == 1
}

func constantTimeEqual(a, b string) bool {
	ab, bb := []byte(a), []byte(b)
	if len(ab) != len(bb) {
		return false
	}
	var v byte
	for i := range ab {
		v |= ab[i] ^ bb[i]
	}
	return v == 0
}
