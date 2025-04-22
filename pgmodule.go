package pgmodule

import (
	"context"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/jackc/pgx/v4"
)

func init() {
	caddy.RegisterModule(PgModule{})
}

type PgModule struct {
	ConnectionString string `json:"connection_string"`
}

func (PgModule) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.pgmodule",
		New: func() caddy.Module { return new(PgModule) },
	}
}

func (m PgModule) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Connect to PostgreSQL
	conn, err := pgx.Connect(context.Background(), m.ConnectionString)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	// Insert request data into PostgreSQL
	_, err = conn.Exec(context.Background(), "INSERT INTO requests (method, path) VALUES ($1, $2)", r.Method, r.URL.Path)
	if err != nil {
		return err
	}

	// Continue to the next handler
	return next.ServeHTTP(w, r)
}
