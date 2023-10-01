package apperror

import (
	"errors"
	"net/http"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				} else if errors.Is(err, ErrNoContent) {
					w.WriteHeader(http.StatusNoContent)
					w.Write(ErrNoContent.Marshal())
					return
				} else if errors.Is(err, ErrUnauthorized) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(ErrUnauthorized.Marshal())
					return
				}

				/* else if errors.Is(err, NoAuthErr) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(ErrNotFound.Marshal())
					return
				} */ // и тд прочие прописанные в error.go/var ошибки

				err = err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(appErr.Marshal())
				return
			}

			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	}
}
