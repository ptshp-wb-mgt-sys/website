// Package middleware/db.go contains database middleware
package middleware

import (
	"context"
	"net/http"
	"pet-mgt/backend/internal/store"
)

// DBContextKey for storing database in request context
const DBContextKey = "db"

// InjectDB injects the database into the request context
func InjectDB(db store.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), DBContextKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetDBFromContext extracts database from request context
func GetDBFromContext(ctx context.Context) (store.Database, bool) {
	db, ok := ctx.Value(DBContextKey).(store.Database)
	return db, ok
}
