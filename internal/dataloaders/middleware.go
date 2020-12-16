package dataloaders

import (
	"context"
	"encoding/json"
	"net/http"

	"bitbucket.org/antuitinc/esp-df-api/internal/datamodels"
)

type userInfo struct {
	Id         string
	Name       string
	GivenName  string
	FamilyName string
	Email      string
	Locale     string
}

// Middleware stores Loaders as a request-scoped context value.
func Middleware(repo datamodels.DBRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			loaders := newLoaders(ctx, repo)
			augmentedCtx := context.WithValue(ctx, key, loaders)

			// Get UserInfo from request header
			var user_info userInfo
			h_user_info := r.Header.Get("user-info")
			json.Unmarshal([]byte(h_user_info), &user_info)

			augmentedCtxA := context.WithValue(augmentedCtx, "userId", user_info.Id)
			augmentedCtxB := context.WithValue(augmentedCtxA, "userLocale", user_info.Locale)

			r = r.WithContext(augmentedCtxB)
			next.ServeHTTP(w, r)
		})
	}
}
