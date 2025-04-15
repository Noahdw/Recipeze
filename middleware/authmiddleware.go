package rmiddleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"recipeze/model"
	"recipeze/service"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

type AuthMiddleware struct {
	store *sessions.CookieStore
	auth  service.AuthService
}

type CtxUserKey struct{}
type CtxGroupAuthorizedKey struct{}

func NewAuthMiddleware(auth service.AuthService) *AuthMiddleware {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.Options.HttpOnly = true
	return &AuthMiddleware{
		store: store,
		auth:  auth,
	}

}

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := a.store.Get(r, "session")
		sessionToken, ok := session.Values["session_token"]
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		sessionTokenStr, ok := sessionToken.(string)
		if !ok || sessionTokenStr == "" {
			next.ServeHTTP(w, r)
			return
		}
		user, err := a.auth.GetLoggedInUser(r.Context(), sessionTokenStr)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), CtxUserKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) AuthorizeGroup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			next.ServeHTTP(w, r)
			return
		}
		groupID, err := getGroupID(r)
		if err != nil {
			return
		}
		inGroup, err := a.auth.IsUserInGroup(r.Context(), groupID, user.ID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), CtxGroupAuthorizedKey{}, inGroup)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *model.User {
	userAny := ctx.Value(CtxUserKey{})
	user, ok := userAny.(*model.User)
	if !ok {
		return nil
	}
	return user
}

func getGroupID(r *http.Request) (int, error) {
	groupIDStr := chi.URLParam(r, "group_id")
	if groupIDStr == "" {
		return 0, fmt.Errorf("no group ID provided")
	}
	return strconv.Atoi(groupIDStr)
}
