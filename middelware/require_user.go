package middelware

import (
	"net/http"

	"../context"
	"../models"
)

type RequireUser struct {
	models.UserService
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("remember_token")
		if err != nil {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}

		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}

		// apply user to request context
		ctx := req.Context()
		ctx = context.WithUser(ctx, user)
		req = req.WithContext(ctx)

		next(res, req)
	})
}
